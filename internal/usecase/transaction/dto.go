package transaction

import "time"

type TransactionDataResponse struct {
	UUID               string    `json:"uuid"`
	TransactionType    string    `json:"transaction_type"`
	SourceAccount      string    `json:"source_account"`
	DestinationAccount string    `json:"destination_account"`
	Status             string    `json:"status"`
	Notes              string    `json:"notes"`
	Amount             int64     `json:"amount"`
	Fee                int64     `json:"fee"`
	ProcessedAt        time.Time `json:"processed_at"`
}

type GetTransactionsRequest struct {
	TransactionType string `query:"transaction_type" json:"transaction_type" validate:"omitempty,only=transfer tapmoney"`
	SourceAccount   string `query:"source_account" json:"source_account" validate:"omitempty,number"`
	Status          string `query:"status" json:"status" validate:"omitempty,only=initiated pending failed completed"`
}

// Map converts the request to a map for repository queries.
func (req *GetTransactionsRequest) Map() map[string]any {
	m := make(map[string]any)
	if req.TransactionType != "" {
		m["transaction_type"] = req.TransactionType
	}
	if req.SourceAccount != "" {
		m["source_account"] = req.SourceAccount
	}
	if req.Status != "" {
		m["status"] = req.Status
	}
	return m
}

type GetTransactionRequest struct {
	UUID string `param:"uuid" json:"uuid" validate:"required,uuid"`
}
