package tapmoney

const SuccessfulMessage = "Payment successful"

type InquiryRequest struct {
	CardNumber string `json:"cardNumber"`
	PocketID   int64  `json:"pocketID"`
	Amount     int64  `json:"amount"`
}

type InquiryResponse struct {
	TransactionID  string `json:"transactionID"`
	PaymentID      string `json:"paymentID"`
	SequenceNumber string `json:"sequenceNumber"`
	Status         string `json:"description"`
	Amount         int64  `json:"amount"`
	CardNumber     string `json:"cardNumber"`
	PocketID       int64  `json:"pocketID"`
}

type PaymentRequest struct {
	TransactionID string `json:"transactionID"`
	Amount        int64  `json:"amount"`
	Notes         string `json:"notes"`
}

type PaymentResponse struct {
	TransactionID string `json:"transactionID"`
	Message       string `json:"message"`
	Status        string `json:"status"`
	Amount        int64  `json:"amount"`
	CardNumber    string `json:"cardNumber"`
	Notes         string `json:"notes"`
	Fee           int64  `json:"fee"`
}
