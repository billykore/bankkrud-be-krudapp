package user

import "time"

type GetByUsernameRequest struct {
	Username string `json:"username" param:"username" validate:"required"`
}

type GetByUsernameResponse struct {
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	FullName    string    `json:"full_name"`
	PhoneNumber string    `json:"phone_number"`
	CIF         string    `json:"cif"`
	Address     string    `json:"address"`
	LastLogin   time.Time `json:"last_login"`
}
