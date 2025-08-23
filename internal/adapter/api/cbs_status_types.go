package api

const CBSSuccessResponseCode = "00"

type CBSGetStatusResponse struct {
	StatusCode        string           `json:"statusCode"`
	StatusDescription string           `json:"statusDescription"`
	Data              CBSGetStatusData `json:"data"`
}

type CBSGetStatusData struct {
	SystemDate    string `json:"systemDate"`
	EODStatus     string `json:"eodStatus"`
	StandInStatus string `json:"standinStatus"`
}

// IsEOD returns true if the EOD status is "STARTED".
func (data *CBSGetStatusData) IsEOD() bool {
	return data.EODStatus == "STARTED"
}

// IsStandIn returns true if the StandIn status is "Y".
func (data *CBSGetStatusData) IsStandIn() bool {
	return data.StandInStatus == "Y"
}
