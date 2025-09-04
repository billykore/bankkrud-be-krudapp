package pocket

import "context"

// Repository defines the interface for pocket repository operations.
type Repository interface {
	// GetByAccountNumber Get retrieves a pocket with opened status
	// from the repository by its account number.
	GetByAccountNumber(ctx context.Context, accountNumber string) (Pocket, error)
}
