package validation

import (
	"context"
	"testing"

	"github.com/romanyx/scraper_auth/internal/change"
	"github.com/stretchr/testify/assert"
)

func Test_Change_Validate(t *testing.T) {
	tests := []struct {
		name    string
		f       change.Form
		wantErr bool
		expect  Errors
	}{
		{
			name: "ok",
			f: change.Form{
				Password:             "secret",
				PasswordConfirmation: "secret",
			},
		},
		{
			name: "passwords mismatch",
			f: change.Form{
				Password:             "secret1",
				PasswordConfirmation: "secret",
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
			f: change.Form{
				Password:             "sec",
				PasswordConfirmation: "sec",
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

			v := NewChange()
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
