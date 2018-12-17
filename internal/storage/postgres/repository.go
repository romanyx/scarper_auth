package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/romanyx/scraper_auth/internal/auth"
	"github.com/romanyx/scraper_auth/internal/reg"
	"github.com/romanyx/scraper_auth/internal/user"
	"github.com/romanyx/scraper_auth/internal/verify"
)

const (
	driverName = "postgres"
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
func (r *Repository) Create(ctx context.Context, u *user.NewUser) (func() error, func() error, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, nil, errors.Wrap(err, "begin tx")
	}

	stmt, err := tx.PrepareNamed(createQuery)
	if err != nil {
		return nil, nil, errors.Wrap(err, "prepare named")
	}

	if _, err := stmt.ExecContext(ctx, map[string]interface{}{
		"account_id":    u.AccountID,
		"email":         u.Email,
		"status":        u.Status,
		"password_hash": u.PasswordHash,
		"token":         u.Token,
	}); err != nil {
		return nil, nil, errors.Wrap(err, "exec context")
	}

	return tx.Commit, tx.Rollback, nil
}

const findQuery = "SELECT id, account_id, email, status, token, password_hash, created_at, updated_at FROM users WHERE account_id=:account_id"

// Find finds user by email.
func (r *Repository) Find(ctx context.Context, accountID string, u *user.User) error {
	stmt, err := r.db.PrepareNamed(findQuery)
	if err != nil {
		return errors.Wrap(err, "prepare named")
	}

	if err := stmt.QueryRowContext(ctx, map[string]interface{}{
		"account_id": accountID,
	}).Scan(&u.ID, &u.AccountID, &u.Email, &u.Status, &u.Token, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return auth.ErrNotFound
		}
		return errors.Wrap(err, "query row scan")
	}

	return nil
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
func (r *Repository) FindByEmail(ctx context.Context, email string, u *user.User) error {
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

	if s == user.StatusNew {
		return auth.ErrNotVerified
	}

	return nil
}

const tokenFindQuery = "SELECT id, account_id FROM users WHERE token=:token"

// FindByToken finds user by token.
func (r *Repository) FindByToken(ctx context.Context, token string, u *user.User) error {
	stmt, err := r.db.PrepareNamed(tokenFindQuery)
	if err != nil {
		return errors.Wrap(err, "prepare named")
	}

	if err := stmt.QueryRowContext(ctx, map[string]interface{}{
		"token": token,
	}).Scan(&u.ID, &u.AccountID); err != nil {
		if err == sql.ErrNoRows {
			return verify.ErrNotFound
		}
		return errors.Wrap(err, "query row scan")
	}

	return nil
}

const verifyQuery = `UPDATE users SET token=null, status='VERIFIED', updated_at=now() WHERE id=:id`

// Verify verifies user.
func (r *Repository) Verify(ctx context.Context, id int32) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}

	stmt, err := tx.PrepareNamed(verifyQuery)
	if err != nil {
		return errors.Wrap(err, "prepare named")
	}

	if _, err := stmt.ExecContext(ctx, map[string]interface{}{
		"id": id,
	}); err != nil {
		return errors.Wrap(err, "exec context")
	}

	return nil
}
