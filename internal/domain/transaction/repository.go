package transaction

import "context"

// Repository defines a contract for data access and persistence operations.
type Repository interface {
	// Get retrieves a transaction entity by its UUID.
	GetByUUID(ctx context.Context, uuid string) (Transaction, error)

	// GetByParams retrieves a list of transaction entities based on the provided parameters.
	GetByParams(ctx context.Context, params map[string]any) ([]Transaction, error)

	// Create creates a transaction entity in the repository.
	Create(ctx context.Context, tx Transaction) error

	// Update updates an existing transaction entity in the repository.
	Update(ctx context.Context, tx Transaction) error
}
