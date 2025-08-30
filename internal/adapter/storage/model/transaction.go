package model

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	UUID                 string
	SourceAccount        string
	DestinationAccount   string
	TransactionType      string
	TransactionReference string
	Status               string
	Note                 string
	Amount               int64
	Fee                  int64
}
