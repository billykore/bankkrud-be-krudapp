package middleware

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/adapter/http/response"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/user"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/config"
)

// AuthenticateUser returns a middleware function that validates token from headers
// and extract user information.
func AuthenticateUser(cfg *config.Configs) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		ContextKey:     string(user.ContextKey),
		SigningKey:     []byte(cfg.Token.Secret),
		SuccessHandler: successHandler,
		ErrorHandler:   errorHandler,
	})
}

// successHandler extract user information from token
// and save the information in the request context.
func successHandler(ctx echo.Context) {
	t := ctx.Get(string(user.ContextKey)).(*jwt.Token)
	usr := userFromToken(t)
	c := ctx.Request().Context()
	c = ContextWithUser(c, usr)
	ctx.SetRequest(ctx.Request().WithContext(c))
}

// errorHandler returns an unauthorized response if there is an authentication error.
func errorHandler(ctx echo.Context, err error) error {
	return ctx.JSON(response.Unauthorized(newAuthenticationError(err.Error())))
}

// authenticationError represents an authentication error.
type authenticationError struct {
	Message string `json:"message"`
}

func newAuthenticationError(message string) *authenticationError {
	return &authenticationError{Message: message}
}

func (e *authenticationError) Error() string {
	return e.Message
}

// ContextWithUser set user data to the ctx context.
func ContextWithUser(ctx context.Context, u user.User) context.Context {
	return context.WithValue(ctx, user.ContextKey, u)
}

// userFromToken returns user information from JWT token.
func userFromToken(token *jwt.Token) user.User {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return user.User{}
	}
	lastLogin, err := time.Parse(time.RFC3339, claims["last_login"].(string))
	if err != nil {
		lastLogin = time.Time{}
	}
	return user.User{
		Username:    claims["sub"].(string),
		Email:       claims["email"].(string),
		PhoneNumber: claims["phone_number"].(string),
		LastLogin:   lastLogin,
	}
}
