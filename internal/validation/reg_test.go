package validation

import (
	"context"
	"testing"

	"github.com/romanyx/scraper_auth/internal/reg"
	"github.com/stretchr/testify/assert"
)

func Test_Reg_Validate(t *testing.T) {
	tests := []struct {
		name       string
		f          reg.Form
		uniqueFunc func(context.Context, string) error
		wantErr    bool
		expect     Errors
	}{
		{
			name: "ok",
			f: reg.Form{
				AccountID:            "accountID",
				Email:                "unique@email.com",
				Password:             "secret",
				PasswordConfirmation: "secret",
			},
			uniqueFunc: func(context.Context, string) error {
				return nil
			},
		},
		{
			name: "not email",
			f: reg.Form{
				AccountID:            "accountID",
				Email:                "unique",
				Password:             "secret",
				PasswordConfirmation: "secret",
			},
			uniqueFunc: func(context.Context, string) error {
				return nil
			},
			wantErr: true,
			expect: Errors{
				Error{
					Field:   "email",
					Message: "must be a valid email address",
				},
			},
		},
		{
			name: "not unique",
			f: reg.Form{
				AccountID:            "accountID",
				Email:                "unique@email.com",
				Password:             "secret",
				PasswordConfirmation: "secret",
			},
			uniqueFunc: func(context.Context, string) error {
				return reg.ErrEmailExists
			},
			wantErr: true,
			expect: Errors{
				Error{
					Field:   "email",
					Message: "email already exists",
				},
			},
		},
		{
			name: "passwords mismatch",
			f: reg.Form{
				AccountID:            "accountID",
				Email:                "unique@email.com",
				Password:             "secret1",
				PasswordConfirmation: "secret",
			},
			uniqueFunc: func(context.Context, string) error {
				return nil
			},
			wantErr: true,
			expect: Errors{
				Error{
					Field:   "password_confirmation",
					Message: "mismatch",
				},
				Error{
					Field:   "password",
					Message: "mismatch",
				},
			},
		},
		{
			name: "passwords length",
			f: reg.Form{
				AccountID:            "accountID",
				Email:                "unique@email.com",
				Password:             "sec",
				PasswordConfirmation: "sec",
			},
			uniqueFunc: func(context.Context, string) error {
				return nil
			},
			wantErr: true,
			expect: Errors{
				Error{
					Field:   "password",
					Message: "the length must be between 6 and 32",
				},
				Error{
					Field:   "password_confirmation",
					Message: "the length must be between 6 and 32",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := NewReg(uniquerFunc(tt.uniqueFunc))
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			err := v.Validate(ctx, &tt.f)

			if tt.wantErr {
				assert.Error(t, err)
				es, ok := err.(Errors)
				assert.True(t, ok)
				assert.Equal(t, tt.expect, es)
				return
			}

			assert.Nil(t, err)
		})
	}
}

type uniquerFunc func(context.Context, string) error

func (f uniquerFunc) Unique(ctx context.Context, email string) error {
	return f(ctx, email)
}
