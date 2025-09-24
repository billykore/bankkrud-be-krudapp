package user

import (
	"context"
	"errors"
)

var ErrTokenNotFound = errors.New("token not found")

// Repository defines the interface for user data persistence.
type Repository interface {
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
