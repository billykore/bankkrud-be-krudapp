package user

import (
	"context"
	"errors"
	"time"
)

// User represents a user in the system.
type User struct {
	UUID        string
	Username    string
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	Password    string
	CIF         string
	Status      string
	Address     string
	DateOfBirth time.Time
	LastLogin   time.Time
}

// FullName returns the full name of the user.
func (u *User) FullName() string {
	if u.FirstName == "" || u.LastName == "" {
		return ""
	}
	return u.FirstName + " " + u.LastName
}

// IsActive checks if the user is active.
func (u *User) IsActive() bool {
	return u.Status == "active"
}

// IsInactive checks if the user is inactive.
func (u *User) IsInactive() bool {
	return u.Status == "inactive"
}

type Token struct {
	Value     string
	ExpiresAt time.Time
}

func (t *Token) ExpiredDuration() int64 {
	return int64(time.Until(t.ExpiresAt).Seconds())
}

func (t *Token) Expired() bool {
	return time.Now().After(t.ExpiresAt)
}

type ContextKeyType string

// ContextKey represents the key for storing user data in the context.
const ContextKey ContextKeyType = "user"

// ErrUserFromContextFailed is returned when the user cannot be retrieved from the context.
var ErrUserFromContextFailed = errors.New("failed to get user from context")

// FromContext retrieves the user from the context.
func FromContext(ctx context.Context) (User, error) {
	user, ok := ctx.Value(ContextKey).(User)
	if !ok {
		return User{}, ErrUserFromContextFailed
	}
	return user, nil
}
