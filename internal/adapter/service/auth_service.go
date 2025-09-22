package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/user"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	tokenSecret   string
	tokenDuration time.Duration
}

func NewAuthService(cfg *config.Configs) *AuthService {
	return &AuthService{
		tokenSecret:   cfg.Token.Secret,
		tokenDuration: cfg.Token.Duration,
	}
}

func (as *AuthService) ValidatePassword(requestPassword, userPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(requestPassword))
}

func (as *AuthService) GenerateToken(u user.User) (user.Token, error) {
	exp := time.Now().Add(as.tokenDuration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": u.Username,
		"exp": exp,
	})
	strToken, err := token.SignedString([]byte(as.tokenSecret))
	if err != nil {
		return user.Token{}, err
	}
	return user.Token{
		Value:     strToken,
		ExpiresAt: exp,
	}, nil
}
