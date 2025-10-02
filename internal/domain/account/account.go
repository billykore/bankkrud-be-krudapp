// Package account contains account domain logic and entities.
package account

// Account represents a bank account entity.
type Account struct {
	CIF           string
	AccountNumber string
	FullName      string
	Type          string
	Balance       int64
}

func (acc Account) CanTransfer(amount int64) bool {
	return acc.Balance >= amount
}
