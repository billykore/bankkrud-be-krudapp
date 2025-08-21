package payment

import "context"

// Service defines the interface for payment operations.
type Service interface {
	// Inquiry performs an inquiry operation for a payment.
	Inquiry(ctx context.Context, channel Channel, bill Bill) (Payment, error)

	// Payment performs a payment operation for a payment.
	Payment(ctx context.Context, Bill Bill) (Payment, error)
}
