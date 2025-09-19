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

	// SaveToken saves a token for a user.
	SaveToken(ctx context.Context, username string, token Token) error

	// GetToken retrieves a token for a user.
	GetToken(ctx context.Context, username string) (Token, error)
}
