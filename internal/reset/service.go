package reset

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/romanyx/scraper_auth/internal/user"
	"github.com/twinj/uuid"
)

//go:generate mockgen -source=service.go -package=reset -destination=service.mock.go

var (
	// ErrNotFound returns when user for given email
	// not found.
	ErrNotFound = errors.New("user not foud")

	uuidMx = sync.Mutex{}
)

// Informer send token to verify email.
type Informer interface {
	Inform(ctx context.Context, u *user.User) error
}

// Repository allows to work with database.
type Repository interface {
	FindByEmail(ctx context.Context, email string, u *user.User) error
	Reset(ctx context.Context, userID int32, token string, exp time.Time) (func() error, func() error, error)
}

// Service holds data requird to verify user.
type Service struct {
	ExpiredAfter time.Duration
	Repository
	Informer
}

// NewService factory returns ready to user
// service.
func NewService(r Repository, i Informer, exp time.Duration) *Service {
	s := Service{
		Repository:   r,
		ExpiredAfter: exp,
		Informer:     i,
	}

	return &s
}

// Reset generates password reset for user.
func (s *Service) Reset(ctx context.Context, email string) error {
	var u user.User
	if err := s.Repository.FindByEmail(ctx, email, &u); err != nil {
		return errors.Wrap(err, "find by email")
	}

	token := uuidStr()
	exp := time.Now().Add(s.ExpiredAfter).UTC()
	commit, rollback, err := s.Repository.Reset(ctx, u.ID, token, exp)
	if err != nil {
		return errors.Wrap(err, "reset token")
	}

	if err := s.Informer.Inform(ctx, &u); err != nil {
		rollback()
		return errors.Wrap(err, "inform user")
	}

	return commit()
}

func uuidStr() string {
	uuidMx.Lock()
	defer uuidMx.Unlock()
	return uuid.NewV4().String()
}
