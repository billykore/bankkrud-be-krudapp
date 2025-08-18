package cbs

import "context"

// Service is the interface for the core banking system service API.
type Service interface {
	// GetStatus retrieves the current status of the core banking system.
	GetStatus(ctx context.Context) (Status, error)
}
