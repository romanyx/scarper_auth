package auth

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/romanyx/scraper_auth/internal/user"
	"github.com/romanyx/scraper_auth/kit/auth"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrNotFound returns when given email is not
	// found in database.
	ErrNotFound = errors.New("email not found")
	// ErrWrongPassword returns when given password
	// is not equal to it's hash in database.
	ErrWrongPassword = errors.New("wrong password")
	// ErrNotVerified returns when user in database is not
	// yet verified.
	ErrNotVerified = errors.New("not verified")
)

// TokenGenerator is the behavior we need in our
// Authenticate to generate tokens for authenticated users.
type TokenGenerator interface {
	GenerateToken(ctx context.Context, claims auth.Claims) (string, error)
}

// EmailFinder allows to find user by it's email.
type EmailFinder interface {
	FindByEmail(ctx context.Context, email string, u *user.User) error
}

// Service holds required data for user
// authentication.
type Service struct {
	ExpireAfter time.Duration
	EmailFinder
	TokenGenerator
}

// NewService factory created ready to service.
func NewService(exp time.Duration, f EmailFinder, tg TokenGenerator) *Service {
	s := Service{
		ExpireAfter:    exp,
		EmailFinder:    f,
		TokenGenerator: tg,
	}

	return &s
}

// Authenticate allows to authenticate user by gine email and password
// and set t Token value as generated token.
func (s *Service) Authenticate(ctx context.Context, email, password string, t *Token) error {
	var u user.User
	if err := s.FindByEmail(ctx, email, &u); err != nil {
		return errors.Wrap(err, "find user by email")
	}

	// Compare the provided password with the saved hash. Use the bcrypt
	// comparison function so it is cryptographically secure.
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return ErrWrongPassword
	}

	// If we are this far the request is valid. Create some claims for the user
	// and generate their token.
	claims := auth.NewClaims(u.AccountID, time.Now(), s.ExpireAfter)

	tknStr, err := s.GenerateToken(ctx, claims)
	if err != nil {
		return errors.Wrap(err, "generate token")
	}
	t.Token = tknStr

	return nil
}

// Token holds token data.
type Token struct {
	Token string
}
