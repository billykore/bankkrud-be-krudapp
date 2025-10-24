package schedule

type GetSchedulesRequest struct {
	Offset          int    `query:"offset"`
	Limit           int    `query:"limit" json:"limit" validate:"omitempty,lte=50"`
	TransactionType string `query:"transaction_type"`
	Status          string `query:"status" json:"status" validate:"omitempty,oneof=active inactive"`
	Period          string `query:"period" json:"period" validate:"omitempty,oneof=daily weekly monthly"`
}

type GetScheduleByUUIDRequest struct {
	UUID string `param:"uuid" validate:"required,uuid"`
}

type GetScheduleResponse struct {
	UUID              string `json:"uuid"`
	Username          string `json:"username"`
	Name              string `json:"name"`
	TransactionType   string `json:"transaction_type"`
	SourceAccount     string `json:"source_account"`
	DestinationNumber string `json:"destination_number"`
	Amount            int64  `json:"amount"`
	Period            string `json:"period"`
	Day               string `json:"day"`
	Date              int    `json:"date"`
	Status            string `json:"status"`
}

type CreateScheduleRequest struct {
	Username          string `json:"username"`
	Name              string `json:"name"`
	TransactionType   string `json:"transaction_type"`
	SourceAccount     string `json:"source_account"`
	DestinationNumber string `json:"destination_number"`
	Amount            int64  `json:"amount"`
	Period            string `json:"period"`
	Day               string `json:"day"`
	Date              int    `json:"date"`
}

type CreateScheduleResponse struct {
	UUID string `json:"uuid"`
}

type UpdateScheduleStatusRequest struct {
	UUID   string `param:"uuid" validate:"required,uuid"`
	Status string `json:"status" validate:"required,oneof=active inactive"`
}

type UpdateScheduleStatusResponse struct {
	Message string `json:"message"`
}

type DeleteScheduleRequest struct {
	UUID string `param:"uuid" validate:"required,uuid"`
}

type DeleteScheduleResponse struct {
	Message string `json:"message"`
}
