// Package transaction contains transaction domain entities.
package transaction

const (
	// InquirySuccess represents a successful transaction status.
	InquirySuccess = "inq-success"
)

// Transaction represents a bank transaction entity.
type Transaction struct {
	UUID                 string
	TransactionReference string
	SourceAccount        string
	DestinationAccount   string
	Amount               int64
	Fee                  int64
	Status               string
	PaymentID            string
}
