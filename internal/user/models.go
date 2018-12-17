package user

import "time"

const (
	// StatusNew initial user status.
	StatusNew = "NEW"
	// StatusVerified status for verifiend users.
	StatusVerified = "VERIFIED"
)

// User contains all user fields.
type User struct {
	ID           int32
	AccountID    string
	Email        string
	Status       string
	PasswordHash string
	Token        *string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser contains information needed to
// create new user.
type NewUser struct {
	Email        string
	AccountID    string
	Status       string
	PasswordHash string
	Token        string
}
