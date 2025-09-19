package model

type User struct {
	Username     string `gorm:"unique"`
	Email        string `gorm:"unique"`
	PhoneNumber  string
	PasswordHash string
	FirstName    string
	LastName     string
}
