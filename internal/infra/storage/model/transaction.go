package model

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UUID                 string `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex"`
	SourceAccount        string
	DestinationAccount   string
	TransactionType      string
	TransactionReference string
	Status               string
	Note                 string
	Amount               int64
	Fee                  int64
	UserUsername         string
}
