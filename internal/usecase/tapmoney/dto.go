package tapmoney

const SuccessfulMessage = "Payment successful"

type InitiateRequest struct {
	CardNumber    string `json:"card_number" validate:"required,min=16,max=19"`
	SourceAccount string `json:"source_account" validate:"required,number"`
	Amount        int64  `json:"amount" validate:"required,min=10000,max=1000000"`
}

type InitiateResponse struct {
	UUID          string `json:"uuid"`
	PaymentID     string `json:"payment_id"`
	Status        string `json:"status"`
	CardNumber    string `json:"card_number"`
	SourceAccount string `json:"source_account"`
	Amount        int64  `json:"amount"`
}

type ProcessRequest struct {
	UUID       string `param:"uuid" json:"uuid" validate:"required,uuid"`
	CardNumber string `json:"card_number" validate:"required,min=16,max=19"`
	Amount     int64  `json:"amount" validate:"required,min=10000,max=1000000"`
	Notes      string `json:"notes" validate:"max=255"`
}

type ProcessResponse struct {
	UUID       string `json:"uuid"`
	Message    string `json:"message"`
	Status     string `json:"status"`
	Amount     int64  `json:"amount"`
	CardNumber string `json:"card_number"`
	Notes      string `json:"notes"`
	Fee        int64  `json:"fee"`
}

type DetailRequest struct {
	UUID string `param:"uuid" json:"transaction_id" validate:"required,uuid"`
}

type DetailResponse struct {
	UUID       string `json:"uuid"`
	Message    string `json:"message"`
	Status     string `json:"status"`
	Amount     int64  `json:"amount"`
	CardNumber string `json:"card_number"`
	Notes      string `json:"notes"`
	Fee        int64  `json:"fee"`
}
