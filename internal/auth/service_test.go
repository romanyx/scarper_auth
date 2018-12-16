package auth

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/romanyx/scraper_auth/kit/auth"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func Test_Service(t *testing.T) {
	pw, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to generate password: %v", err)
	}

	tests := []struct {
		name      string
		tokenFunc func(context.Context, auth.Claims) (string, error)
		findFunc  func(context.Context, string, *User) error
		wantErr   bool
		expect    Token
	}{
		{
			name: "ok",
			tokenFunc: func(context.Context, auth.Claims) (string, error) {
				return "token", nil
			},
			findFunc: func(ctx context.Context, email string, u *User) error {
				u.PasswordHash = string(pw)
				return nil
			},
			expect: Token{
				Token: "token",
			},
		},
		{
			name: "not found",
			findFunc: func(ctx context.Context, email string, u *User) error {
				return ErrNotFound
			},
			wantErr: true,
		},
		{
			name: "wrong password",
			findFunc: func(ctx context.Context, email string, u *User) error {
				return nil
			},
			wantErr: true,
		},
		{
			name: "token gen",
			tokenFunc: func(context.Context, auth.Claims) (string, error) {
				return "", errors.New("mock err")
			},
			findFunc: func(ctx context.Context, email string, u *User) error {
				u.PasswordHash = string(pw)
				return nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := NewService(time.Hour, emailFinderFunc(tt.findFunc), tokenGeneratorFunc(tt.tokenFunc))

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var got Token
			err := s.Authenticate(ctx, "", "password", &got)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, got, tt.expect)
		})
	}
}

type tokenGeneratorFunc func(context.Context, auth.Claims) (string, error)

func (f tokenGeneratorFunc) GenerateToken(ctx context.Context, c auth.Claims) (string, error) {
	return f(ctx, c)
}

type emailFinderFunc func(context.Context, string, *User) error

func (f emailFinderFunc) FindByEmail(ctx context.Context, email string, u *User) error {
	return f(ctx, email, u)
}
