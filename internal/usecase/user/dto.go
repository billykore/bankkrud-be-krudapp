package user

import "time"

type GetByUsernameRequest struct {
	Fields string `json:"fields" query:"fields" validate:"omitempty,only=cif email phone_number address date_of_birth last_login"`
}

type GetByUsernameResponse struct {
	Username    string    `json:"username"`
	FullName    string    `json:"full_name"`
	CIF         string    `json:"cif,omitempty"`
	Email       string    `json:"email,omitempty"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Address     string    `json:"address,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitzero"`
	LastLogin   time.Time `json:"last_login,omitzero"`
}
