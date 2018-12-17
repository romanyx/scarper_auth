package change

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/romanyx/scraper_auth/internal/user"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=service.go -package=change -destination=service.mock.go

const (
	validationErrMsg = "you have validation errors"
)

var (
	// ErrNotFound returns when token not found.
	ErrNotFound = errors.New("token not found")
	// ErrTokenExpired returns when token is expired.
	ErrTokenExpired = errors.New("token expired")

	uuidMx = sync.Mutex{}
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

// ValidationErrors holds validation errors
// list.
type ValidationErrors []ValidationError

// Error implements error interface to return
// slice as error for futher manipulations.
func (v ValidationErrors) Error() string {
	return validationErrMsg
}

// ValidationError holds field and message
// of validation exception.
type ValidationError struct {
	Field, Message string
}

// Validater validates user fields.
type Validater interface {
	Validate(ctx context.Context, f *Form) error
}

// Repository allows to work with database.
type Repository interface {
	FindResetToken(ctx context.Context, token string, r *Token) error
	ChangePassword(ctx context.Context, id int32, passwordHash string) error
	FindByID(ctx context.Context, id int32, u *user.User) error
}

// Service holds data required to verify user.
type Service struct {
	Repository
	Validater
}

// NewService factory returns ready to user
// service.
func NewService(r Repository) *Service {
	s := Service{
		Repository: r,
		Validater:  &ozzo{},
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

	if err := s.Repository.FindByID(ctx, t.UserID, u); err != nil {
		return errors.Wrap(err, "find by id")
	}

	return nil
}

func uuidStr() string {
	uuidMx.Lock()
	defer uuidMx.Unlock()
	return uuid.NewV4().String()
}

type ozzo struct{}

func (v ozzo) Validate(ctx context.Context, form *Form) error {
	return nil

}
