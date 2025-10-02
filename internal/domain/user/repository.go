package user

import (
	"context"
	"errors"
)

var (
	// ErrDuplicateUserData is returned when user data already exists.
	ErrDuplicateUserData = errors.New("duplicate user data")

	// ErrUserNotFound is returned when a user is not found.
	ErrUserNotFound = errors.New("user not found")

	// ErrTokenNotFound is returned when a token is not found.
	ErrTokenNotFound = errors.New("token not found")
)

// Repository defines the interface for user data persistence.
type Repository interface {
	// CreateUser creates a new user.
	Create(ctx context.Context, user User) error

	// GetByUsername retrieves a user by their username.
	GetByUsername(ctx context.Context, username string) (User, error)

	// GetFieldsByUsername retrieves a user's fields by their username.
	GetFieldsByUsername(ctx context.Context, username string, fields ...string) (User, error)

	// SaveToken saves a token for a user.
	SaveToken(ctx context.Context, username string, token Token) error

	// GetToken retrieves a token for a user.
	GetToken(ctx context.Context, username string) (Token, error)

	// DeleteToken deletes a token for a user.
	DeleteToken(ctx context.Context, username string) error
}
