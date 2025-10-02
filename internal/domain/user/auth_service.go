package user

// AuthService is an interface for user authentication and authorization.
type AuthService interface {
	// GenerateToken generates a token for the given user.
	GenerateToken(user User) (Token, error)

	// HashPassword hashes the given password.
	HashPassword(password string) (string, error)

	// ValidatePassword validates the password for the given user.
	ValidatePassword(requestPassword, userPassword string) error
}
