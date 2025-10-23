package schedule

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("schedule not found")
)

type Filter map[string]string

type Query struct {
	Limit  int
	Offset int
	Filter Filter
}

// Repository defines the interface for schedule data operations.
type Repository interface {
	// GetAll retrieves all schedules within the given query.
	GetAll(ctx context.Context, query Query) ([]Schedule, error)

	// GetByUUID retrieves a schedule by its UUID.
	GetByUUID(ctx context.Context, uuid string) (Schedule, error)

	// Create adds a new schedule.
	Create(ctx context.Context, schedule Schedule) error

	//UpdateStatus modifies the status of an existing schedule.
	UpdateStatus(ctx context.Context, uuid string, status string) error

	// Delete removes a schedule by its UUID.
	Delete(ctx context.Context, uuid string) error
}
