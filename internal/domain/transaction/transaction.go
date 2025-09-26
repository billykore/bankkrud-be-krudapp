// Package transaction contains transaction domain entities.
package transaction

import "time"

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
