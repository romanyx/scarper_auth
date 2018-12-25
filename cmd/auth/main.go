package main

import (
	"crypto/rsa"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "net/http/pprof"

	_ "github.com/lib/pq"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/heptiolabs/healthcheck"
	"github.com/pkg/errors"
	"github.com/romanyx/scraper_auth/internal/auth"
	grpcCli "github.com/romanyx/scraper_auth/internal/broker/grpc"
	"github.com/romanyx/scraper_auth/internal/change"
	"github.com/romanyx/scraper_auth/internal/notifier/smtp"
	"github.com/romanyx/scraper_auth/internal/reg"
	"github.com/romanyx/scraper_auth/internal/reset"
	"github.com/romanyx/scraper_auth/internal/storage/postgres"
	"github.com/romanyx/scraper_auth/internal/validation"
	"github.com/romanyx/scraper_auth/internal/verify"
	authEng "github.com/romanyx/scraper_auth/kit/auth"
	"github.com/romanyx/scraper_auth/proto"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"google.golang.org/grpc"
)

const (
	shutdownTimeout = 15 * time.Second
	readTimeout     = 15 * time.Second
	writeTimeout    = 15 * time.Second
	alg             = "RS256"
)

var (
	version = "unset"
)

func main() {
	var (
		addr           = flag.String("addr", ":8080", "address of gRPC server")
		dsn            = flag.String("dsn", "", "postgres database DSN")
		tokenExp       = flag.Duration("expire", time.Hour, "token live time")
		privateKeyFile = flag.String("key", "", "private key file path")
		keyID          = flag.String("id", "", "private key id")
		healthAddr     = flag.String("health", ":8081", "health check addr")
		debugAddr      = flag.String("debug", ":1234", "debug server addr")
	)

	// Setup db connection
	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		log.Fatalf("failed to connect db: %v\n", err)
	}
	defer db.Close()

	// Authentication setup.
	keyContents, err := ioutil.ReadFile(*privateKeyFile)
	if err != nil {
		log.Fatalf("reading auth private key: %v", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyContents)
	if err != nil {
		log.Fatalf("parsing auth private key: %v", err)
	}

	publicKeyLookup := authEng.NewSingleKeyFunc(*keyID, key.Public().(*rsa.PublicKey))
	authenticator, err := authEng.NewAuthenticator(key, *keyID, alg, publicKeyLookup)
	if err != nil {
		log.Fatalf("constructing authenticator: %v", err)
	}

	// Register stats and trace exporters to export
	// the collected data.
	view.RegisterExporter(&exporter.PrintExporter{})

	// Register the views to collect server request count.
	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		log.Fatal(err)
	}

	// Setup gRPC server.
	grpcServer := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}))

	// Setup handlers.
	srv := setupServer(authenticator, nil, db, *tokenExp)
	proto.RegisterAuthServer(grpcServer, srv)

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Make a channel for errors.
	errChan := make(chan error)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			errChan <- errors.Wrap(err, "failed to serve grpc")
		}
	}()

	// Debug server.
	debugServer := setupDebugServer(*debugAddr)
	go func() {
		if err := debugServer.ListenAndServe(); err != nil {
			errChan <- errors.Wrap(err, "debug server")
		}
	}()
	defer debugServer.Close()

	// Health checker handler.
	health := healthcheck.NewHandler()

	// Build and start healt server.
	healthServer := http.Server{
		Addr:         *healthAddr,
		Handler:      health,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	go func() {
		if err := healthServer.ListenAndServe(); err != nil {
			errChan <- errors.Wrap(err, "health server")
		}
	}()
	defer healthServer.Close()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		log.Fatalf("critical error: %v\n", err)
	case <-osSignals:
		log.Println("stopping by signal")
		grpcServer.GracefulStop()
	}
}

func setupServer(ath *authEng.Authenticator, inf *smtp.Client, db *sql.DB, exp time.Duration) proto.AuthServer {
	logger := logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.JSONFormatter{},
		Level:     logrus.DebugLevel,
	}
	entry := logrus.NewEntry(&logger)

	repo := postgres.NewRepository(db)

	regSrv := reg.NewService(repo, validation.NewReg(repo), inf)
	wrpReg := grpcCli.NewRegistraterWithTracing(grpcCli.NewRegistraterWithLogrus(regSrv, entry))

	authSrv := auth.NewService(exp, repo, ath)
	wrpAuth := grpcCli.NewAuthenticaterWithTracing(grpcCli.NewAuthenticaterWithLogrus(authSrv, entry))

	vrfSrv := verify.NewService(repo)
	wrpVrf := grpcCli.NewVerifierWithTracing(grpcCli.NewVerifierWithLogrus(vrfSrv, entry))

	rstSrv := reset.NewService(repo, inf, exp)
	wrpRst := grpcCli.NewPwdReseterWithTracing(grpcCli.NewPwdReseterWithLogrus(rstSrv, entry))

	chgSrv := change.NewService(repo, validation.NewChange())
	wrpChg := grpcCli.NewPwdChangerWithTracing(grpcCli.NewPwdChangerWithLogrus(chgSrv, entry))

	srv := grpcCli.NewServer(wrpReg, wrpAuth, wrpVrf, wrpRst, wrpChg)
	return srv
}

func setupDebugServer(addr string) *http.Server {
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, version)
	})

	s := http.Server{
		Addr:         addr,
		Handler:      http.DefaultServeMux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	return &s
}
