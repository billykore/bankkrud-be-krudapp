package transfer

type Transfer struct {
	SourceAccount        string
	DestinationAccount   string
	Amount               int64
	Fee                  int64
	Status               string
	Notes                string
	TransactionID        string
	TransactionReference string
}
