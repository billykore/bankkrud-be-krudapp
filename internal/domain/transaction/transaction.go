// Package transaction contains transaction domain entities.
package transaction

import (
	"errors"
	"fmt"
	"time"
)

const (
	// StatusInitiated represents an initiated transaction status.
	StatusInitiated = "initiated"
	// StatusPending represents a pending transaction status.
	StatusPending = "pending"
	// StatusFailed represents a successful failed status.
	StatusFailed = "failed"
	// StatusCompleted represents a completed transaction status.
	StatusCompleted = "completed"
)

var (
	// ErrInvalidStateTransition is returned when a transaction state transition is invalid.
	ErrInvalidStateTransition = errors.New("invalid transaction state transition")
)

// Transaction represents a bank transaction entity.
type Transaction struct {
	UUID                 string
	TransactionReference string
	SourceAccount        string
	DestinationAccount   string
	TransactionType      string
	Status               string
	PaymentID            string
	Note                 string
	Amount               int64
	Fee                  int64
	Username             string
	ProcessedAt          time.Time
}

// Complete marks the transaction as completed.
func (t *Transaction) Complete(txRef string) error {
	if t.Status != StatusInitiated && t.Status != StatusPending {
		return fmt.Errorf("%w: cannot complete transaction from status %s", ErrInvalidStateTransition, t.Status)
	}
	t.Status = StatusCompleted
	t.TransactionReference = txRef
	t.ProcessedAt = time.Now()
	return nil
}

// Fail marks the transaction as failed.
func (t *Transaction) Fail() error {
	if t.Status == StatusCompleted {
		return fmt.Errorf("%w: cannot fail completed transaction", ErrInvalidStateTransition)
	}
	t.Status = StatusFailed
	t.ProcessedAt = time.Now()
	return nil
}

// IsProcessable checks if the transaction is in a state that can be processed.
func (t *Transaction) IsProcessable() bool {
	return t.Status == StatusInitiated || t.Status == StatusPending
}
