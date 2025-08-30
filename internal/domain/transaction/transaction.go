// Package transaction contains transaction domain entities.
package transaction

const (
	// StatusPending represents a pending transaction status.
	StatusPending = "pending"
	// StatusInquirySuccess represents a successful transaction status.
	StatusInquirySuccess = "inq-success"
	// StatusSuccess represents a successful transaction status.
	StatusSuccess = "success"
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
	Notes                string
	Amount               int64
	Fee                  int64
}
