package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
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

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
