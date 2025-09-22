package tapmoney

const SuccessfulMessage = "Payment successful"

type InquiryRequest struct {
	CardNumber    string `json:"card_number"`
	SourceAccount string `json:"source_account"`
	Amount        int64  `json:"amount"`
}

type InquiryResponse struct {
	TransactionID string `json:"transaction_id"`
	PaymentID     string `json:"payment_id"`
	Status        string `json:"status"`
	CardNumber    string `json:"card_number"`
	SourceAccount string `json:"source_account"`
	Amount        int64  `json:"amount"`
}

type PaymentRequest struct {
	TransactionID string `json:"transaction_id"`
	Amount        int64  `json:"amount"`
	Notes         string `json:"notes"`
}

type PaymentResponse struct {
	TransactionID string `json:"transaction_id"`
	Message       string `json:"message"`
	Status        string `json:"status"`
	Amount        int64  `json:"amount"`
	CardNumber    string `json:"card_number"`
	Notes         string `json:"notes"`
	Fee           int64  `json:"fee"`
}
