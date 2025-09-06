package api

// Response represents the response structure
// for HTTP responses in the TapMoney service.
type Response[T any] struct {
	Success bool           `json:"success"`
	Data    T              `json:"data,omitempty"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

// ErrorResponse represents the error structure
// for HTTP responses in the TapMoney service.
type ErrorResponse struct {
	Name    string `json:"name,omitempty"`
	Message string `json:"message,omitempty"`
}

// InquiryRequest represents the request structure
// for an inquiry operation in the TapMoney system.
type InquiryRequest struct {
	CardNumber    string `json:"cardNumber"`
	SourceAccount string `json:"sourceAccount"`
	Amount        int64  `json:"amount"`
}

// InquiryResponse represents the response structure
// for an inquiry operation in the TapMoney system.
type InquiryResponse struct {
	SequenceNumber string `json:"sequenceNumber"`
	Status         string `json:"description"`
	TraceId        string `json:"traceId"`
	Amount         int64  `json:"amount"`
	CardNumber     string `json:"cardNumber"`
	SourceAccount  string `json:"sourceAccount"`
}

// PaymentRequest represents the request structure
// for a payment operation in the TapMoney system.
type PaymentRequest struct {
	TransactionID string `json:"transactionID"`
	Amount        int64  `json:"amount"`
	Notes         string `json:"notes"`
}

// PaymentResponse represents the response structure
// for a payment operation in the TapMoney system.
type PaymentResponse struct {
	TransactionID string `json:"transactionID"`
	Message       string `json:"message"`
	Status        string `json:"status"`
	Amount        int64  `json:"amount"`
	CardNumber    string `json:"cardNumber"`
	Notes         string `json:"notes"`
	Fee           int64  `json:"fee"`
}
