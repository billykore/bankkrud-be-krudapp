// Package payment defines the domain model for payments in a payment gateway system.
package payment

// Payment represents a payment in the payment gateway system.
type Payment struct {
	ID      string
	Status  string
	Channel Channel
	Bill    Bill
	Biller  Biller
}

// Channel represents a payment channel in the payment gateway system.
type Channel struct {
	ID string
}

// Biller represents a biller in the payment gateway system.
type Biller struct {
	Code string
	Name string
	Type string
}

// Bill represents a bill in the payment gateway system.
type Bill struct {
	BillNumber         string
	BillerCode         string
	Amount             int64
	SourceAccount      string
	DestinationAccount string
	Fee                int64
	FreeFee            bool
}
