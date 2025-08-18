package pocket

import "context"

type Repository interface {
	// Get retrieves a pocket with opened status from the repository by its ID.
	Get(ctx context.Context, id int64) (Pocket, error)
}
