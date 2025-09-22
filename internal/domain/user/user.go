package user

import "time"

// User represents a user in the system.
type User struct {
	UUID        string
	Username    string
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	Password    string
	DateOfBirth string
	CIF         string
	Status      string
	Address     string
	LastLogin   time.Time
}

// FullName returns the full name of the user.
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// IsActive checks if the user is active.
func (u *User) IsActive() bool {
	return u.Status == "active"
}

// IsInactive checks if the user is inactive.
func (u *User) IsInactive() bool {
	return u.Status == "inactive"
}

type Token struct {
	Value     string
	ExpiresAt time.Time
}

func (t *Token) ExpiredDuration() int64 {
	return int64(time.Until(t.ExpiresAt).Seconds())
}

func (t *Token) Expired() bool {
	return time.Now().After(t.ExpiresAt)
}
