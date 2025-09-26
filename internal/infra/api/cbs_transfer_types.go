package api

type TransferRequest struct {
	CreditAccountNumber string `json:"noRekeningKredit"`
	DebitAccountNumber  string `json:"noRekeningDebit"`
	Amount              string `json:"nominal"`
	TransactionInfo     string `json:"keteranganTransaksi"`
	Fee                 string `json:"biaya"`
}

type TransferResponse struct {
	StatusCode           string   `json:"statusCode"`
	StatusDescription    string   `json:"statusDescription"`
	JournalSequence      string   `json:"journalSequence"`
	TransactionReference string   `json:"transactionReference"`
	ABMsg                []string `json:"abmsg"`
}
