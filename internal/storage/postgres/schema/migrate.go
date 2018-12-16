package schema

import (
	"database/sql"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	bindata "github.com/golang-migrate/migrate/source/go_bindata"
	"github.com/pkg/errors"
)

const (
	source          = "go-bindata"
	database        = "postgres"
	migrationsTable = "versions"
)

//go:generate go-bindata -prefix migrations/ -pkg schema -o migrations.bindata.go migrations/

// Migrate migrates schema to given database
// connection.
func Migrate(db *sql.DB) error {
	m, err := newMigration(db)
	if err != nil {
		return errors.Wrap(err, "new migration")
	}

	if err := m.Up(); err != nil {
		return errors.Wrap(err, "migrate schema")
	}

	return nil
}

func newMigration(db *sql.DB) (*migrate.Migrate, error) {
	r := bindata.Resource(AssetNames(), Asset)

	s, err := bindata.WithInstance(r)
	if err != nil {
		return nil, errors.Wrap(err, "prepare source instance")
	}

	cfg := postgres.Config{
		MigrationsTable: migrationsTable,
	}
	d, err := postgres.WithInstance(db, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "prepare database instance")
	}

	m, err := migrate.NewWithInstance(source, s, database, d)
	if err != nil {
		return nil, errors.Wrap(err, "prepare migrate instance")
	}

	return m, nil
}
