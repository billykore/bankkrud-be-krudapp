package user

// AuthService is an interface for user authentication and authorization.
type AuthService interface {
	// GenerateToken generates a token for the given user.
	GenerateToken(user User) (Token, error)

	// ValidatePassword validates the password for the given user.
	ValidatePassword(requestPassword, userPassword string) error
}
