package api

import "strconv"

type GetAccountRequest struct {
	TransactionType string `json:"tipeTransaksi"`
	AccountNumber   string `json:"noRekening"`
}

type GetAccountResponse struct {
	StatusCode           string         `json:"statusCode"`
	StatusDescription    string         `json:"statusDescription"`
	ErrorCode            string         `json:"errorCode"`
	JournalSequence      string         `json:"journalSequence"`
	TransactionReference string         `json:"transactionReference"`
	Data                 AccountDetails `json:"data"`
}

type AccountDetails struct {
	AccountNumber    string `json:"noRekening"`
	AccountType      string `json:"tipeRekening"`
	Name             string `json:"nama"`
	Currency         string `json:"mataUang"`
	Status           string `json:"status"`
	BlockedStatus    string `json:"blokir"`
	Balance          string `json:"saldo"`
	MinBalance       string `json:"saldoMinimum"`
	AvailableBalance string `json:"saldoTersedia"`
	CID              string `json:"cif"`
	ProductType      string `json:"tipeProduk"`
}

// GetBalance returns int64 balance.
func (details *AccountDetails) GetBalance() int64 {
	balance, _ := strconv.ParseInt(details.Balance, 10, 64)
	return balance
}
