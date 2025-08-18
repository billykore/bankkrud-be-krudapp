package tapmoney

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
