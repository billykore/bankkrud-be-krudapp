package api

import "context"

// Service is the interface for the TapMoney service.
type Service interface {
	// Inquiry performs an inquiry operation for a TapMoney transaction.
	Inquiry(ctx context.Context, request InquiryRequest) (Response[InquiryResponse], error)
	// Payment performs a payment operation for a TapMoney transaction.
	Payment(ctx context.Context, req PaymentRequest) (Response[PaymentResponse], error)
}
