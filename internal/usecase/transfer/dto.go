package transfer

type InitiateRequest struct {
	SourceAccount      string `json:"sourceAccount"`
	DestinationAccount string `json:"destinationAccount"`
	Amount             int64  `json:"amount"`
}

type InitiateResponse struct {
	TransactionID string `json:"transactionID"`
	Status        string `json:"status"`
}

type ProcessRequest struct {
	TransactionID      string `json:"transactionID"`
	SourceAccount      string `json:"sourceAccount"`
	DestinationAccount string `json:"destinationAccount"`
	Amount             int64  `json:"amount"`
}

type ProcessResponse struct {
	TransactionID string `json:"transactionID"`
	Status        string `json:"status"`
}
