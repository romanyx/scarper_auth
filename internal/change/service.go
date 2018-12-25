package change

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/romanyx/scraper_auth/internal/user"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=service.go -package=change -destination=service.mock.go

var (
	// ErrNotFound returns when token not found.
	ErrNotFound = errors.New("token not found")
	// ErrTokenExpired returns when token is expired.
	ErrTokenExpired = errors.New("token expired")
)

// Form contains data
// required to change password
type Form struct {
	Password, PasswordConfirmation string
}

// Token contains data
// about reset.
type Token struct {
	ExpiredAt time.Time
	UserID    int32
}

var (
	// ErrEmailExists returns when given email is already
	// present in database.
	ErrEmailExists = errors.New("email not found")
)

// Validater validates user fields.
type Validater interface {
	Validate(ctx context.Context, f *Form) error
}

// Repository allows to work with database.
type Repository interface {
	FindResetToken(ctx context.Context, token string, r *Token) error
	ChangePassword(ctx context.Context, id int32, passwordHash string) error
	Find(ctx context.Context, id int32, u *user.User) error
}

// Service holds data required to verify user.
type Service struct {
	Repository
	Validater
}

// NewService factory returns ready to user
// service.
func NewService(r Repository, v Validater) *Service {
	s := Service{
		Repository: r,
		Validater:  v,
	}

	return &s
}

// Change allows to change user password
func (s *Service) Change(ctx context.Context, token string, form *Form, u *user.User) error {
	var t Token
	if err := s.Repository.FindResetToken(ctx, token, &t); err != nil {
		return errors.Wrap(err, "find reset token")
	}

	if time.Now().UTC().After(t.ExpiredAt) {
		return ErrTokenExpired
	}

	if err := s.Validater.Validate(ctx, form); err != nil {
		return errors.Wrap(err, "validation")
	}

	pw, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "generating password hash")
	}

	if err := s.Repository.ChangePassword(ctx, t.UserID, string(pw)); err != nil {
		return errors.Wrap(err, "change password")
	}

	if err := s.Repository.Find(ctx, t.UserID, u); err != nil {
		return errors.Wrap(err, "find by id")
	}

	return nil
}
