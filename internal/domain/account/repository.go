package account

import "context"

// Repository defines a contract for account data access and persistence operations.
type Repository interface {
	// Get retrieves an account from the repository by its account number.
	Get(ctx context.Context, accountNumber string) (Account, error)
}
