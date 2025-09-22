package model

import (
	"time"
)

type Token struct {
	Value     string    `json:"value"`
	ExpiresAt time.Time `json:"expires_at"`
}
