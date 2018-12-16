package main

import (
	"flag"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest"
	"github.com/romanyx/scraper_auth/internal/storage/postgres/schema"
	"github.com/romanyx/scraper_auth/kit/docker"
	"github.com/romanyx/scraper_auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		os.Exit(m.Run())
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %v", err)
	}

	pgDocker, err := docker.NewPostgres(pool)
	if err != nil {
		log.Fatalf("prepare postgres with docker: %v", err)
	}

	schema.Migrate(pgDocker.DB)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := setupServer(pgDocker.DB, time.Hour)
	s := grpc.NewServer()
	proto.RegisterAuthServer(s, srv)
	reflection.Register(s)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	code := m.Run()

	pgDocker.DB.Close()
	if err := pool.Purge(pgDocker.Resource); err != nil {
		log.Fatalf("could not purge postgres docker: %v", err)
	}

	os.Exit(code)
}
