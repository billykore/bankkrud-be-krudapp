package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/payment"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/config"
)

// PaymentGateway is a middleware payment gateway.
type PaymentGateway struct{}

func NewPaymentGateway(cfg *config.Configs, client *http.Client) *PaymentGateway {
	return &PaymentGateway{}
}

func (pg *PaymentGateway) Inquiry(ctx context.Context, channel payment.Channel, bill payment.Bill) (payment.Payment, error) {
	return payment.Payment{
		ID:      uuid.New().String(),
		Status:  "success",
		Channel: channel,
		Bill:    bill,
	}, nil
}

func (pg *PaymentGateway) Payment(ctx context.Context, bill payment.Bill) (payment.Payment, error) {
	return payment.Payment{
		ID: uuid.New().String(),
	}, nil
}
