package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func (s *AuthService) ValidatePassword(requestPassword, userPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(requestPassword))
}

func (s *AuthService) GenerateToken(u user.User) (user.Token, error) {
	exp := time.Now().Add(s.tokenDuration)

	claims := jwt.MapClaims{
		"iss":           "api.bankkrud.com",
		"aud":           jwt.ClaimStrings{"app.bankkrud.com"},
		"exp":           jwt.NewNumericDate(exp),
		"nbf":           jwt.NewNumericDate(time.Now()),
		"iat":           jwt.NewNumericDate(time.Now()),
		"jti":           uuid.New().String(),
		"sub":           u.Username,
		"cif":           u.CIF,
		"email":         u.Email,
		"phone_number":  u.PhoneNumber,
		"full_name":     u.FullName(),
		"address":       u.Address,
		"date_of_birth": u.DateOfBirth,
		"last_login":    u.LastLogin,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	strToken, err := token.SignedString([]byte(s.tokenSecret))
	if err != nil {
		return user.Token{}, err
	}
	return user.Token{
		Value:     strToken,
		ExpiresAt: exp,
	}, nil
}
