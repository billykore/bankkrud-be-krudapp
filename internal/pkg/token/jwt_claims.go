package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	FullName    string    `json:"full_name"`
	LastLogin   time.Time `json:"last_login"`
	jwt.RegisteredClaims
}
