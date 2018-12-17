package main

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"github.com/romanyx/scraper_auth/internal/auth"
	grpcCli "github.com/romanyx/scraper_auth/internal/grpc"
	"github.com/romanyx/scraper_auth/internal/reg"
	"github.com/romanyx/scraper_auth/internal/storage/postgres"
	"github.com/romanyx/scraper_auth/internal/user"
	"github.com/romanyx/scraper_auth/internal/verify"
	authEng "github.com/romanyx/scraper_auth/kit/auth"
	"github.com/romanyx/scraper_auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	alg = "RS256"
)

func main() {
	var (
		addr           = flag.String("addr", ":8080", "address of gRPC server")
		dsn            = flag.String("dsn", "", "postgres database DSN")
		tokenExp       = flag.Duration("expire", time.Hour, "token live time")
		privateKeyFile = flag.String("key", "", "private key file path")
		keyID          = flag.String("id", "", "private key id")
	)

	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		log.Fatalf("failed to connect db: %v\n", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	keyContents, err := ioutil.ReadFile(*privateKeyFile)
	if err != nil {
		log.Fatalf("main : Reading auth private key : %v", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyContents)
	if err != nil {
		log.Fatalf("main : Parsing auth private key : %v", err)
	}

	publicKeyLookup := authEng.NewSingleKeyFunc(*keyID, key.Public().(*rsa.PublicKey))
	authenticator, err := authEng.NewAuthenticator(key, *keyID, alg, publicKeyLookup)
	if err != nil {
		log.Fatalf("main : Constructing authenticator : %v", err)
	}

	srv := setupServer(authenticator, db, *tokenExp)
	s := grpc.NewServer()
	proto.RegisterAuthServer(s, srv)
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func setupServer(ath *authEng.Authenticator, db *sql.DB, exp time.Duration) proto.AuthServer {
	repo := postgres.NewRepository(db)
	authSrv := auth.NewService(exp, repo, ath)
	regSrv := reg.NewService(repo, &informer{})
	vrfSrv := verify.NewService(repo)

	srv := grpcCli.NewServer(regSrv, authSrv, vrfSrv)
	return srv
}

type informer struct{}

func (i *informer) Inform(context.Context, *user.User) error {
	return nil
}
