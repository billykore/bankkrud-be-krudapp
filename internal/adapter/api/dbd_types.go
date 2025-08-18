package api

const (
	SuccessResponseCode = "00"
)

type DBDInquiryRequest struct {
	ChannelID      string                       `json:"ChannelId"`
	CoreInfo       CoreInfoInquiryRequest       `json:"CoreInfo"`
	ThirdPartyInfo ThirdPartyInfoInquiryRequest `json:"ThirdPartyInfo"`
}
type CoreInfoInquiryRequest struct {
	SourceAccount string `json:"SourceAccount"`
	Amount        string `json:"Amount"`
}
type ThirdPartyInfoInquiryRequest struct {
	BillNumber         string `json:"BillNumber"`
	BillerCode         string `json:"BillerCode"`
	DestinationAccount string `json:"DestinationAccount"`
	BankCode           string `json:"BankCode"`
	UserName           string `json:"UserName"`
	Month              string `json:"Month"`
}

type DBDInquiryResponse struct {
	ResponseCode        string                        `json:"ResponseCode"`
	ResponseDescription string                        `json:"ResponseDescription"`
	TraceID             string                        `json:"TraceId"`
	CoreInfo            CoreInfoInquiryResponse       `json:"CoreInfo"`
	ThirdPartyInfo      ThirdPartyInfoInquiryResponse `json:"ThirdPartyInfo"`
}

type CoreInfoInquiryResponse struct {
	StatusDescription string `json:"StatusDescription"`
}

type ThirdPartyInfoInquiryResponse struct {
	StatusDescription      string        `json:"StatusDescription"`
	CustomerName           string        `json:"CustomerName"`
	TarifDaya              string        `json:"TarifDaya"`
	LembarTagihan          string        `json:"LembarTagihan"`
	BulanTahun             string        `json:"BulanTahun"`
	StandMeter             string        `json:"Standmeter"`
	Nominal                string        `json:"Nominal"`
	ReferenceNumber        string        `json:"ReferenceNumber"`
	RegistrationNumber     string        `json:"RegistrationNumber"`
	NoMeter                string        `json:"NoMeter"`
	BillNumber             string        `json:"BillNumber"`
	IdentityInfoName       string        `json:"IdentityInfoName"`
	IdentityName           string        `json:"IdentityName"`
	Daya                   string        `json:"Daya"`
	Flag                   string        `json:"Flag"`
	Ut                     string        `json:"UT"`
	RegisNumber            string        `json:"RegisNumber"`
	InquiryData            string        `json:"InquiryData"`
	AdminFee               string        `json:"AdminFee"`
	DestinationAccountName string        `json:"DestinationAccountName"`
	TransactionType        string        `json:"TransactionType"`
	NPWP                   string        `json:"NPWP"`
	Bills                  []BillDetails `json:"billList"`
	BillAmount             string        `json:"BillAmount"`
	BillerReceipt          string        `json:"BillerReceipt"`
	BillingAmount          string        `json:"BillingAmount"`
	Amount                 string        `json:"Amount"`
	ServicePeriod          string        `json:"ServicePeriod"`
}

type BillDetails struct {
	ReferenceNumberBill1 string `json:"Reference Number Bill 1"`
	TotalBill1           string `json:"Total Bill 1"`
	ReferenceNumberBill2 string `json:"Reference Number Bill 2"`
	TotalBill2           string `json:"Total Bill 2"`
	ReferenceNumberBill3 string `json:"Reference Number Bill 3"`
	TotalBill3           string `json:"Total Bill 3"`
}
