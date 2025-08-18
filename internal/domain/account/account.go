// Package account contains account domain logic and entities.
package account

// Account represents a bank account entity.
type Account struct {
	AccountNumber string
	FullName      string
	Type          string
	Balance       int64
}
