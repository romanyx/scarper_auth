package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/romanyx/scraper_auth/internal/auth"
	"github.com/romanyx/scraper_auth/internal/reg"
	"github.com/twinj/uuid"
)

const (
	driverName = "postgres"

	statusNew = "NEW"
)

// Repository holds crud actions.
type Repository struct {
	db *sqlx.DB
}

// NewRepository returns ready to work repository.
func NewRepository(db *sql.DB) *Repository {
	r := Repository{
		db: sqlx.NewDb(db, driverName),
	}

	return &r
}

const createQuery = `INSERT INTO users (email, status, token, account_id, password_hash) VALUES (:email, :status, :token, :account_id, :password_hash)`

// Create insert user into database.
func (r *Repository) Create(ctx context.Context, u *reg.User) (func() error, func() error, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, nil, errors.Wrap(err, "begin tx")
	}

	stmt, err := tx.PrepareNamed(createQuery)
	if err != nil {
		return nil, nil, errors.Wrap(err, "prepare named")
	}

	tk := uuid.NewV4()
	u.Token = tk.String()

	if _, err := stmt.ExecContext(ctx, map[string]interface{}{
		"account_id":    u.AccountID,
		"email":         u.Email,
		"status":        statusNew,
		"password_hash": u.PasswordHash,
		"token":         u.Token,
	}); err != nil {
		return nil, nil, errors.Wrap(err, "exec context")
	}

	return tx.Commit, tx.Rollback, nil
}

const uniqueQuery = "SELECT COUNT(*) FROM users WHERE email=:email"

// Unique checks that email is unique.
func (r *Repository) Unique(ctx context.Context, email string) error {
	stmt, err := r.db.PrepareNamed(uniqueQuery)
	if err != nil {
		return errors.Wrap(err, "prepare named")
	}

	var c int
	if err := stmt.QueryRowContext(ctx, map[string]interface{}{
		"email": email,
	}).Scan(c); err != nil {
		return reg.ErrEmailExists
	}

	return nil
}

const emailFindQuery = "SELECT password_hash, account_id, status FROM users WHERE email=:email"

// FindByEmail finds user by email.
func (r *Repository) FindByEmail(ctx context.Context, email string, u *auth.User) error {
	stmt, err := r.db.PrepareNamed(emailFindQuery)
	if err != nil {
		return errors.Wrap(err, "prepare named")
	}

	var s string
	if err := stmt.QueryRowContext(ctx, map[string]interface{}{
		"email": email,
	}).Scan(&u.PasswordHash, &u.AccountID, &s); err != nil {
		if err == sql.ErrNoRows {
			return auth.ErrNotFound
		}
		return errors.Wrap(err, "query row scan")
	}

	if s == statusNew {
		return auth.ErrNotVerified
	}

	return nil
}
