package reg

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	validationErrMsg = "you have validation errors"
)

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

// Informer send token to verify email.
type Informer interface {
	Inform(ctx context.Context, u *User) error
}

// Repository allows to work with database.
type Repository interface {
	Create(ctx context.Context, u *User) (func() error, func() error, error)
	Unique(ctx context.Context, email string) error
}

// Service holds everything required to registrate
// user.
type Service struct {
	Validater
	Repository
	Informer
}

// NewService factory prepares service for
// futher operations.
func NewService(r Repository, i Informer) *Service {
	s := Service{
		Repository: r,
		Informer:   i,
		Validater: ozzo{
			Repository: r,
		},
	}

	return &s
}

// Registrate registrates user.
func (s *Service) Registrate(ctx context.Context, f *Form) error {
	if err := s.Validate(ctx, f); err != nil {
		return errors.Wrap(err, "validate user")
	}

	u := User{
		AccountID: f.AccountID,
		Email:     f.Email,
	}
	pw, err := bcrypt.GenerateFromPassword([]byte(f.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "generating password hash")
	}
	u.PasswordHash = string(pw)

	commit, rollback, err := s.Create(ctx, &u)
	if err != nil {
		return errors.Wrap(err, "create user")
	}

	if err := s.Inform(ctx, &u); err != nil {
		rollback()
		return errors.Wrap(err, "inform")
	}

	return commit()
}

// Form is a registraton form.
type Form struct {
	Email                string
	AccountID            string
	Password             string
	PasswordConfirmation string
}

// User used to insert model.
type User struct {
	Email        string
	AccountID    string
	PasswordHash string
	Token        string
}

type ozzo struct {
	Repository
}

func (v ozzo) Validate(ctx context.Context, f *Form) error {
	return nil
}
