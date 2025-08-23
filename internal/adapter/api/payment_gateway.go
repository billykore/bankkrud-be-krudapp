package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.bankkrud.com/backend/svc/tapmoney/internal/domain/payment"
	"go.bankkrud.com/backend/svc/tapmoney/internal/pkg/config"
)

// PaymentGateway is a middleware payment gateway.
type PaymentGateway struct {
	client   *http.Client
	addr     string
	username string
	password string
}

func NewPaymentGateway(cfg *config.Configs, client *http.Client) *PaymentGateway {
	return &PaymentGateway{
		client:   client,
		addr:     cfg.DBD.Addr,
		username: cfg.DBD.Username,
		password: cfg.DBD.Password,
	}
}

func (pg *PaymentGateway) Inquiry(ctx context.Context, channel payment.Channel, bill payment.Bill) (payment.Payment, error) {
	url := pg.addr + "/invoke/raya.bill.payment.providers:inquiry"

	b, err := json.Marshal(PaymentGatewayInquiryRequest{
		ChannelID: channel.ID,
		CoreInfo: CoreInfoInquiryRequest{
			SourceAccount: bill.SourceAccount,
		},
		ThirdPartyInfo: ThirdPartyInfoInquiryRequest{
			UserName:   "uid-1231",
			BillNumber: bill.DestinationAccount,
			BillerCode: bill.BillerCode,
		},
	})
	if err != nil {
		return payment.Payment{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return payment.Payment{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(pg.username, pg.password)

	resp, err := pg.client.Do(req)
	if err != nil {
		return payment.Payment{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return payment.Payment{}, errors.New("failed to perform inquiry, status code: " + resp.Status)
	}

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return payment.Payment{}, err
	}

	var apiRes PaymentGatewayInquiryResponse
	err = json.Unmarshal(b, &apiRes)
	if err != nil {
		return payment.Payment{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if apiRes.ResponseCode != DBDSuccessResponseCode {
		return payment.Payment{}, fmt.Errorf("[%v]: %v", apiRes.ResponseCode, apiRes.ResponseDescription)
	}

	return payment.Payment{
		ID:      apiRes.TraceID,
		Status:  apiRes.ResponseDescription,
		Channel: channel,
		Bill:    bill,
	}, nil
}

func (pg *PaymentGateway) Payment(ctx context.Context, Bill payment.Bill) (payment.Payment, error) {
	//TODO implement me
	panic("implement me")
}
