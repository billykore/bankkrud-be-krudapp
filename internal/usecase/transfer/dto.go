package transfer

type InitiateRequest struct {
	SourceAccount      string `json:"source_account" validate:"required,number"`
	DestinationAccount string `json:"destination_account" validate:"required,number"`
	Amount             int64  `json:"amount" validate:"required,gte=1000,lte=50000000"`
}

type InitiateResponse struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}

type ProcessRequest struct {
	TransactionID      string `json:"transaction_id" validate:"required,uuid"`
	SourceAccount      string `json:"source_account" validate:"required,number"`
	DestinationAccount string `json:"destination_account" validate:"required,number"`
	Amount             int64  `json:"amount" validate:"required"`
}

type ProcessResponse struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}
