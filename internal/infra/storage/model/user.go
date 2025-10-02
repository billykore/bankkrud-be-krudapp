package model

import "time"

type User struct {
	Username     string `gorm:"unique"`
	Email        string `gorm:"unique"`
	CIF          string `gorm:"unique"`
	PhoneNumber  string `gorm:"unique"`
	PasswordHash string
	FirstName    string
	LastName     string
	Address      string
	DateOfBirth  time.Time
	LastLogin    time.Time
	Status       string
}
