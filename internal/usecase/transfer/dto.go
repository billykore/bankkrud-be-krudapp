package transfer

type InitiateRequest struct {
	SourceAccount      string `json:"source_account"`
	DestinationAccount string `json:"destination_account"`
	Amount             int64  `json:"amount"`
}

type InitiateResponse struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}

type ProcessRequest struct {
	TransactionID      string `json:"transaction_id"`
	SourceAccount      string `json:"source_account"`
	DestinationAccount string `json:"destination_account"`
	Amount             int64  `json:"amount"`
}

type ProcessResponse struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}
