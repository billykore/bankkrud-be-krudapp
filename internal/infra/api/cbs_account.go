package api

import (
	"context"
	"math/rand"
	"strings"

	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/account"
)

const digits = "0123456789"

// CBSAccountAPI is the core banking system service API for getting account information.
type CBSAccountAPI struct{}

// NewAccountAPI creates a new instance of the AccountAPI.
func NewCBSAccountAPI() *CBSAccountAPI {
	return &CBSAccountAPI{}
}

func (api *CBSAccountAPI) Get(ctx context.Context, accountNumber string) (account.Account, error) {
	return account.Account{
		AccountNumber: accountNumber,
		FullName:      "John Doe",
		Type:          "Savings",
		Balance:       10000000,
	}, nil
}

func (api *CBSAccountAPI) Create(ctx context.Context, username string) (account.Account, error) {
	var cifBuilder strings.Builder
	defer cifBuilder.Reset()
	for _ = range 10 {
		cifBuilder.WriteByte(digits[rand.Intn(10)])
	}
	return account.Account{
		CIF: cifBuilder.String(),
	}, nil
}
