package account

import "context"

// Service defines a contract for account data access and persistence operations.
type Service interface {
	// Get retrieves an account from the repository by its account number.
	Get(ctx context.Context, accountNumber string) (Account, error)

	// Transfer moves amount from one account to another
	// and returns an error if the operation fails.
	Transfer(ctx context.Context, srcAccountNumber string, destAccountNumber string, amount int64) error
}
