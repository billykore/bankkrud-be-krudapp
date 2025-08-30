package model

import "gorm.io/gorm"

const (
	SakuStatusOpened = "Opened"
)

type Pocket struct {
	gorm.Model
	UUID          string
	UserID        string
	Name          string
	AccountNumber string
	Status        string
	SavingType    int
}
