package reg

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/romanyx/scraper_auth/internal/user"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=service.go -package=reg -destination=service.mock.go

const (
	validationErrMsg = "you have validation errors"
)

var (
	// ErrEmailExists returns when given email is already
	// present in database.
	ErrEmailExists = errors.New("email not found")

	uuidMx = sync.Mutex{}
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
	Inform(ctx context.Context, u *user.User) error
}

// Repository allows to work with database.
type Repository interface {
	Create(ctx context.Context, u *user.NewUser) (func() error, func() error, error)
	Find(ctx context.Context, accountID string, u *user.User) error
	Unique(ctx context.Context, email string) error
}

// Service holds everything required to registrate
// user.
type Service struct {
	Validater
	Repository
	Informer
	*sync.Mutex
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
		Mutex: &sync.Mutex{},
	}

	return &s
}

// Registrate registrates user.
func (s *Service) Registrate(ctx context.Context, f *Form, usr *user.User) error {
	if err := s.Validate(ctx, f); err != nil {
		return errors.Wrap(err, "validate user")
	}

	u := user.NewUser{
		AccountID: f.AccountID,
		Email:     f.Email,
	}
	pw, err := bcrypt.GenerateFromPassword([]byte(f.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "generating password hash")
	}
	u.PasswordHash = string(pw)
	u.Status = user.StatusNew
	u.Token = uuidStr()

	commit, rollback, err := s.Create(ctx, &u)
	if err != nil {
		return errors.Wrap(err, "create user")
	}

	if err := s.Find(ctx, f.AccountID, usr); err != nil {
		rollback()
		return errors.Wrap(err, "find user")
	}

	if err := s.Inform(ctx, usr); err != nil {
		rollback()
		return errors.Wrap(err, "inform")
	}

	return commit()
}

func uuidStr() string {
	uuidMx.Lock()
	defer uuidMx.Unlock()
	return uuid.NewV4().String()
}

// Form is a registraton form.
type Form struct {
	Email                string
	AccountID            string
	Password             string
	PasswordConfirmation string
}

type ozzo struct {
	Repository
}

func (v ozzo) Validate(ctx context.Context, f *Form) error {
	return nil
}
