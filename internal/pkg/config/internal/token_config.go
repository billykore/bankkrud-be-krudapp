package internal

import "time"

// Token config.
type Token struct {
	Secret    string
	HeaderKid string
	Duration  time.Duration
}
