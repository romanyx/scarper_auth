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

var (
	// ErrEmailExists returns when given email is already
	// present in database.
	ErrEmailExists = errors.New("email already exists")

	uuidMx = sync.Mutex{}
)

// Validater validates user fields.
type Validater interface {
	Validate(ctx context.Context, f *Form) error
}

// Informer send token to verify email.
type Informer interface {
	Verify(ctx context.Context, u *user.User) error
}

// Repository allows to work with database.
type Repository interface {
	Create(ctx context.Context, u *user.NewUser, usr *user.User) (func() error, func() error, error)
	FindByAccountID(ctx context.Context, accountID string, u *user.User) error
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
func NewService(r Repository, v Validater, i Informer) *Service {
	s := Service{
		Repository: r,
		Informer:   i,
		Validater:  v,
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

	commit, rollback, err := s.Create(ctx, &u, usr)
	if err != nil {
		return errors.Wrap(err, "create user")
	}

	if err := s.Verify(ctx, usr); err != nil {
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
