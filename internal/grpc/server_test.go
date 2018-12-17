package gprc

import (
	"context"
	"errors"
	"testing"

	"github.com/romanyx/scraper_auth/internal/auth"
	"github.com/romanyx/scraper_auth/internal/change"
	"github.com/romanyx/scraper_auth/internal/reg"
	"github.com/romanyx/scraper_auth/internal/reset"
	"github.com/romanyx/scraper_auth/internal/user"
	"github.com/romanyx/scraper_auth/internal/verify"
	"github.com/romanyx/scraper_auth/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_Server_SignIn(t *testing.T) {
	tests := []struct {
		name     string
		authFunc func(ctx context.Context, email, password string, t *auth.Token) error
		code     codes.Code
		expect   *proto.SignInResponse
	}{
		{
			name: "ok",

			authFunc: func(ctx context.Context, email, password string, t *auth.Token) error {
				t.Token = "token"
				return nil
			},
			code: codes.OK,
			expect: &proto.SignInResponse{
				Token: "token",
			},
		},
		{
			name: "not found",

			authFunc: func(ctx context.Context, email, password string, t *auth.Token) error {
				return auth.ErrNotFound
			},
			code: codes.InvalidArgument,
		},
		{
			name: "wrong password",

			authFunc: func(ctx context.Context, email, password string, t *auth.Token) error {
				return auth.ErrWrongPassword
			},
			code: codes.InvalidArgument,
		},
		{
			name: "not verified",

			authFunc: func(ctx context.Context, email, password string, t *auth.Token) error {
				return auth.ErrNotVerified
			},
			code: codes.InvalidArgument,
		},
		{
			name: "internal",

			authFunc: func(ctx context.Context, email, password string, t *auth.Token) error {
				return errors.New("mock error")
			},
			code: codes.Internal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := NewServer(nil, authenticaterFunc(tt.authFunc), nil, nil, nil)

			req := proto.SignInRequest{
				Email:    "email",
				Password: "password",
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			got, err := s.SignIn(ctx, &req)
			assert.Equal(t, got, tt.expect)

			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, st.Code(), tt.code)
		})
	}
}

type authenticaterFunc func(context.Context, string, string, *auth.Token) error

func (f authenticaterFunc) Authenticate(ctx context.Context, email string, password string, t *auth.Token) error {
	return f(ctx, email, password, t)
}

func Test_Server_SignUp(t *testing.T) {
	tests := []struct {
		name    string
		regFunc func(context.Context, *reg.Form, *user.User) error
		code    codes.Code
	}{
		{
			name: "ok",
			regFunc: func(context.Context, *reg.Form, *user.User) error {
				return nil
			},
			code: codes.OK,
		},
		{
			name: "validation error",
			regFunc: func(context.Context, *reg.Form, *user.User) error {
				return make(reg.ValidationErrors, 0)
			},
			code: codes.InvalidArgument,
		},
		{
			name: "internal",
			regFunc: func(context.Context, *reg.Form, *user.User) error {
				return errors.New("mock error")
			},
			code: codes.Internal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := NewServer(registraterFunc(tt.regFunc), nil, nil, nil, nil)

			req := proto.SignUpRequest{
				Email:    "email",
				Password: "password",
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			_, err := s.SignUp(ctx, &req)
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, st.Code(), tt.code)
		})
	}
}

type registraterFunc func(context.Context, *reg.Form, *user.User) error

func (rf registraterFunc) Registrate(c context.Context, f *reg.Form, u *user.User) error {
	return rf(c, f, u)
}

func Test_Server_Verify(t *testing.T) {
	tests := []struct {
		name       string
		verifyFunc func(context.Context, string, *user.User) error
		code       codes.Code
	}{
		{
			name: "ok",
			verifyFunc: func(context.Context, string, *user.User) error {
				return nil
			},
			code: codes.OK,
		},
		{
			name: "not found",
			verifyFunc: func(context.Context, string, *user.User) error {
				return verify.ErrNotFound
			},
			code: codes.NotFound,
		},
		{
			name: "internal",
			verifyFunc: func(context.Context, string, *user.User) error {
				return errors.New("mock error")
			},
			code: codes.Internal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := NewServer(nil, nil, verifierFunc(tt.verifyFunc), nil, nil)

			req := proto.EmailVerifyRequest{
				Token: "token",
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			_, err := s.Verify(ctx, &req)
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, st.Code(), tt.code)
		})
	}
}

type verifierFunc func(context.Context, string, *user.User) error

func (f verifierFunc) Verify(ctx context.Context, token string, u *user.User) error {
	return f(ctx, token, u)
}

func Test_Server_Reset(t *testing.T) {
	tests := []struct {
		name      string
		resetFunc func(context.Context, string) error
		code      codes.Code
	}{
		{
			name: "ok",
			resetFunc: func(context.Context, string) error {
				return nil
			},
			code: codes.OK,
		},
		{
			name: "not found",
			resetFunc: func(context.Context, string) error {
				return reset.ErrNotFound
			},
			code: codes.NotFound,
		},
		{
			name: "internal",
			resetFunc: func(context.Context, string) error {
				return errors.New("mock error")
			},
			code: codes.Internal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := NewServer(nil, nil, nil, pwdReseterFunc(tt.resetFunc), nil)

			req := proto.PasswordResetRequest{
				Email: "john@example.com",
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			_, err := s.Reset(ctx, &req)
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, st.Code(), tt.code)
		})
	}
}

type pwdReseterFunc func(context.Context, string) error

func (f pwdReseterFunc) Reset(ctx context.Context, email string) error {
	return f(ctx, email)
}

func Test_Server_Change(t *testing.T) {
	tests := []struct {
		name       string
		changeFunc func(context.Context, string, *change.Form, *user.User) error
		code       codes.Code
	}{
		{
			name: "ok",
			changeFunc: func(context.Context, string, *change.Form, *user.User) error {
				return nil
			},
			code: codes.OK,
		},
		{
			name: "not found",
			changeFunc: func(context.Context, string, *change.Form, *user.User) error {
				return change.ErrNotFound
			},
			code: codes.NotFound,
		},
		{
			name: "expired",
			changeFunc: func(context.Context, string, *change.Form, *user.User) error {
				return change.ErrTokenExpired
			},
			code: codes.InvalidArgument,
		},
		{
			name: "validation errors",
			changeFunc: func(context.Context, string, *change.Form, *user.User) error {
				return make(change.ValidationErrors, 0)
			},
			code: codes.InvalidArgument,
		},
		{
			name: "internal",
			changeFunc: func(context.Context, string, *change.Form, *user.User) error {
				return errors.New("mock error")
			},
			code: codes.Internal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := NewServer(nil, nil, nil, nil, pwdChangerFunc(tt.changeFunc))

			req := proto.PasswordChangeRequest{
				Token:                "token",
				Password:             "password",
				PasswordConfirmation: "password_confirmation",
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			_, err := s.Change(ctx, &req)
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, st.Code(), tt.code)
		})
	}
}

type pwdChangerFunc func(context.Context, string, *change.Form, *user.User) error

func (f pwdChangerFunc) Change(ctx context.Context, token string, form *change.Form, u *user.User) error {
	return f(ctx, token, form, u)
}
