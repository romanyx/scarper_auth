package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net"
	"time"

	"github.com/romanyx/scraper_auth/internal/auth"
	grpcCli "github.com/romanyx/scraper_auth/internal/grpc"
	"github.com/romanyx/scraper_auth/internal/reg"
	"github.com/romanyx/scraper_auth/internal/storage/postgres"
	authEng "github.com/romanyx/scraper_auth/kit/auth"
	"github.com/romanyx/scraper_auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var (
		addr     = flag.String("addr", ":8080", "address of gRPC server")
		dsn      = flag.String("dsn", "", "postgres database DSN")
		tokenExp = flag.Duration("expire", time.Hour, "token live time")
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

	srv := setupServer(db, *tokenExp)
	s := grpc.NewServer()
	proto.RegisterAuthServer(s, srv)
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func setupServer(db *sql.DB, exp time.Duration) proto.AuthServer {
	var ath authEng.Authenticator

	repo := postgres.NewRepository(db)
	authSrv := auth.NewService(exp, repo, &ath)
	regSrv := reg.NewService(repo, &informer{})

	srv := grpcCli.NewServer(regSrv, authSrv)
	return srv
}

type informer struct{}

func (i *informer) Inform(context.Context, *reg.User) error {
	return nil
}
