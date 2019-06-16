package validation

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/pkg/errors"
	"github.com/romanyx/scraper_auth/internal/reg"
)

// Uniquer check that email is exists in database.
type Uniquer interface {
	Unique(ctx context.Context, email string) error
}

// Reg allows to validate user on sign up.
type Reg struct {
	Uniquer
}

// NewReg factory allows to initialize reg validater.
func NewReg(u Uniquer) *Reg {
	r := Reg{
		Uniquer: u,
	}

	return &r
}

// Validate validates registration form.
func (v *Reg) Validate(ctx context.Context, f *reg.Form) error {
	var es Errors

	if err := validation.Validate(f.Email, validation.Required, is.Email); err != nil {
		es = append(es, Error{
			Field:   "email",
			Message: err.Error(),
		})
	}

	if err := validation.Validate(f.AccountID, validation.Required); err != nil {
		es = append(es, Error{
			Field:   "account_id",
			Message: err.Error(),
		})
	}

	if err := v.Uniquer.Unique(ctx, f.Email); err != nil {
		switch errors.Cause(err) {
		case reg.ErrEmailExists:
			es = append(es, Error{
				Field:   "email",
				Message: err.Error(),
			})
		default:
			return errors.Wrap(err, "unique")
		}
	}

	fmt.Println(f.Password, f.PasswordConfirmation)

	if err := validation.Validate(f.Password, validation.Required, validation.Length(6, 32)); err != nil {
		es = append(es, Error{
			Field:   "password",
			Message: err.Error(),
		})
	}

	if err := validation.Validate(f.PasswordConfirmation, validation.Required, validation.Length(6, 32)); err != nil {
		es = append(es, Error{
			Field:   "password_confirmation",
			Message: err.Error(),
		})
	}

	if f.Password != f.PasswordConfirmation {
		es = append(es, Error{
			Field:   "password_confirmation",
			Message: "mismatch",
		}, Error{
			Field:   "password",
			Message: "mismatch",
		})
	}

	if len(es) > 0 {
		return es
	}

	return nil
}
