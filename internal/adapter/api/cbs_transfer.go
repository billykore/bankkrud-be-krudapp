package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/transfer"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/config"
)

type TransferAPI struct {
	client   *http.Client
	cbsAuth  *CBSAuth
	addr     string
	username string
	password string
}

func NewTransferAPI(cfg *config.Configs, client *http.Client, cbsAuth *CBSAuth) *TransferAPI {
	return &TransferAPI{
		client:   client,
		cbsAuth:  cbsAuth,
		addr:     cfg.CBS.Addr,
		username: cfg.CBS.Username,
		password: cfg.CBS.Password,
	}
}

func (ta *TransferAPI) Transfer(ctx context.Context, srcAccountNumber, destAccountNumber string, amount int64, remark string) (transfer.Transfer, error) {
	url := ta.addr + "/api/transaction"

	b, err := json.Marshal(TransferRequest{
		CreditAccountNumber: srcAccountNumber,
		DebitAccountNumber:  destAccountNumber,
		Amount:              strconv.FormatInt(amount, 10),
		TransactionInfo:     remark,
		Fee:                 transferFee,
	})
	if err != nil {
		return transfer.Transfer{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return transfer.Transfer{}, err
	}

	token, err := ta.cbsAuth.GetToken(ctx)
	if err != nil {
		return transfer.Transfer{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := ta.client.Do(req)
	if err != nil {
		return transfer.Transfer{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return transfer.Transfer{}, errors.New("failed to perform get status, status code: " + resp.Status)
	}

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return transfer.Transfer{}, err
	}

	var apiRes TransferResponse
	err = json.Unmarshal(b, &apiRes)
	if err != nil {
		return transfer.Transfer{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if apiRes.StatusCode != CBSSuccessResponseCode {
		return transfer.Transfer{}, fmt.Errorf("[%v]: %v", apiRes.StatusCode, apiRes.StatusDescription)
	}

	return transfer.Transfer{
		SourceAccount:        srcAccountNumber,
		DestinationAccount:   destAccountNumber,
		Amount:               amount,
		Status:               "success",
		TransactionReference: apiRes.TransactionReference,
	}, nil
}
