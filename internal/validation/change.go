package validation

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/romanyx/scraper_auth/internal/change"
)

// Change allows to validate user on sign up.
type Change struct {
	Uniquer
}

// NewChange factory allows to initialize change validater.
func NewChange() *Change {
	v := Change{}

	return &v
}

// Validate validates registration form.
func (v *Change) Validate(ctx context.Context, f *change.Form) error {
	var es Errors

	if err := validation.Validate(f.Password, validation.Length(6, 32)); err != nil {
		es = append(es, Error{
			Field:   "password",
			Message: err.Error(),
		})
	}

	if err := validation.Validate(f.PasswordConfirmation, validation.Length(6, 32)); err != nil {
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
