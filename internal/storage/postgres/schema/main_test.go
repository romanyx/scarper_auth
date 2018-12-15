package schema

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest"
	"github.com/romanyx/scraper_auth/kit/docker"
)

var (
	db *sql.DB
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

	db = pgDocker.DB

	code := m.Run()

	db.Close()
	if err := pool.Purge(pgDocker.Resource); err != nil {
		log.Fatalf("could not purge postgres docker: %v", err)
	}

	os.Exit(code)
}
