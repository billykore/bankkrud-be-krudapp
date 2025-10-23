package model

import "gorm.io/gorm"

// Schedule represents a schedule record in the database.
type Schedule struct {
	gorm.Model
	UUID              string `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex"`
	Username          string
	Name              string
	TransactionType   string
	SourceAccount     string
	DestinationNumber string
	Amount            int64
	Period            string
	Status            string
	CronExpression    string
}
