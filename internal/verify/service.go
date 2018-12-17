package verify

import (
	"context"

	"github.com/pkg/errors"
	"github.com/romanyx/scraper_auth/internal/user"
)

//go:generate mockgen -source=service.go -package=verify -destination=service.mock.go

var (
	// ErrNotFound returns when token not found in
	// database.
	ErrNotFound = errors.New("token not found")
)

// Repository allows to find user by its token,
// and verify it.
type Repository interface {
	FindByToken(context.Context, string, *user.User) error
	Verify(context.Context, int32) error
	Find(context.Context, string, *user.User) error
}

// Service holds data requird to verify user.
type Service struct {
	Repository
}

// NewService factory returns ready to user
// service.
func NewService(r Repository) *Service {
	s := Service{
		Repository: r,
	}

	return &s
}

// Verify verifies user.
func (s *Service) Verify(ctx context.Context, token string, u *user.User) error {
	if err := s.FindByToken(ctx, token, u); err != nil {
		return errors.Wrap(err, "find user by token")
	}

	if err := s.Repository.Verify(ctx, u.ID); err != nil {
		return errors.Wrap(err, "verify user")
	}

	if err := s.Repository.Find(ctx, u.AccountID, u); err != nil {
		return errors.Wrap(err, "find user")
	}

	return nil
}
