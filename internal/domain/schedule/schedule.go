package schedule

const (
	MonthlyPeriod = "monthly"
	WeeklyPeriod  = "weekly"
	DailyPeriod   = "daily"
)

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)

// Schedule represents a scheduled transaction.
type Schedule struct {
	ID                int
	UUID              string
	Username          string
	Name              string
	TransactionType   string
	SourceAccount     string
	DestinationNumber string
	Amount            int64
	Period            string
	Status            string
	Day               string
	Date              int
}
