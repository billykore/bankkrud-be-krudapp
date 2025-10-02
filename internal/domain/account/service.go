package account

import "context"

// Repository defines a contract for account data access and persistence operations.
// Repository can be an API, database, or any other service that provides account data.
type Repository interface {
	// Get retrieves an account from the repository by its account number.
	Get(ctx context.Context, accountNumber string) (Account, error)

	// Create creates a new account in the repository.
	Create(ctx context.Context, username string) (Account, error)
}
