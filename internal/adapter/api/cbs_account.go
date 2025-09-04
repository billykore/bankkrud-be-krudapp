package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/account"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/config"
)

const TransactionTypeInquiry = "inquiry"

// AccountAPI is the core banking system service API for getting account information.
type AccountAPI struct {
	cbsAuth  *CBSAuth
	client   *http.Client
	addr     string
	username string
	password string
}

// NewAccountAPI creates a new instance of the AccountAPI.
func NewAccountAPI(cfg *config.Configs, client *http.Client, cbsAuth *CBSAuth) *AccountAPI {
	return &AccountAPI{
		client:   client,
		cbsAuth:  cbsAuth,
		addr:     cfg.CBS.Addr,
		username: cfg.CBS.Username,
		password: cfg.CBS.Password,
	}
}

func (ca *AccountAPI) Get(ctx context.Context, accountNumber string) (account.Account, error) {
	url := ca.addr + "/api/transaction"

	b, err := json.Marshal(GetAccountRequest{
		AccountNumber:   accountNumber,
		TransactionType: TransactionTypeInquiry,
	})
	if err != nil {
		return account.Account{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return account.Account{}, err
	}

	token, err := ca.cbsAuth.getToken(ctx)
	if err != nil {
		return account.Account{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := ca.client.Do(req)
	if err != nil {
		return account.Account{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return account.Account{}, errors.New("failed to perform get status, status code: " + resp.Status)
	}

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return account.Account{}, err
	}
	var apiRes GetAccountResponse
	err = json.Unmarshal(b, &apiRes)
	if err != nil {
		return account.Account{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if apiRes.StatusCode != CBSSuccessResponseCode {
		return account.Account{}, fmt.Errorf("[%v]: %v", apiRes.StatusCode, apiRes.StatusDescription)
	}

	return account.Account{
		AccountNumber: apiRes.Data.AccountNumber,
		FullName:      apiRes.Data.Name,
		Type:          apiRes.Data.AccountType,
		Balance:       apiRes.Data.GetBalance(),
	}, nil
}
