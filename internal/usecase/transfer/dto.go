package transfer

import "time"

type InitiateRequest struct {
	SourceAccount      string `json:"source_account" validate:"required,number"`
	DestinationAccount string `json:"destination_account" validate:"required,number"`
	Amount             int64  `json:"amount" validate:"required,gte=1000,lte=50000000"`
	Note               string `json:"note"`
}

type InitiateResponse struct {
	UUID   string `json:"uuid"`
	Status string `json:"status"`
}

type ProcessRequest struct {
	UUID               string `param:"uuid" json:"uuid" validate:"required,uuid"`
	SourceAccount      string `json:"source_account" validate:"required,number"`
	DestinationAccount string `json:"destination_account" validate:"required,number"`
	Amount             int64  `json:"amount" validate:"required"`
}

type ProcessResponse struct {
	UUID   string `json:"uuid"`
	Status string `json:"status"`
}

type DetailRequest struct {
	UUID string `param:"uuid" json:"uuid" validate:"required,uuid"`
}

type DetailResponse struct {
	UUID               string    `json:"uuid"`
	Status             string    `json:"status"`
	Amount             int64     `json:"amount"`
	Fee                int64     `json:"fee"`
	SourceAccount      string    `json:"source_account"`
	DestinationAccount string    `json:"destination_account"`
	Note               string    `json:"note"`
	ProcessedAt        time.Time `json:"processed_at"`
}
