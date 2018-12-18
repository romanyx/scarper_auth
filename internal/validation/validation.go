package validation

const (
	errMsg = "you have validation errors"
)

// Errors holds validation errors
// list.
type Errors []Error

// Error implements error interface to return
// slice as error for futher manipulations.
func (v Errors) Error() string {
	return errMsg
}

// Error holds field and message
// of validation exception.
type Error struct {
	Field, Message string
}
