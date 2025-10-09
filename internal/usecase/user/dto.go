package user

import (
	"encoding/json"
	"time"
)

type CreateRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=100"`
	Password    string `json:"password" validate:"required,min=8,max=100"`
	FirstName   string `json:"first_name" validate:"required,min=3,max=100"`
	LastName    string `json:"last_name" validate:"required,min=3,max=100"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required,phonenumber"`
	Address     string `json:"address" validate:"required,min=3,max=200"`
	DateOfBirth string `json:"date_of_birth" validate:"required"`
}

type CreateResponse struct {
	Message string `json:"message"`
}

type GetByUsernameRequest struct {
	Fields string `json:"fields" query:"fields" validate:"omitempty,only=cif email phone_number address date_of_birth last_login"`
}

func (req *GetByUsernameRequest) MarshalBinary() ([]byte, error) {
	return json.Marshal(req)
}

func (req *GetByUsernameRequest) String() string {
	b, err := req.MarshalBinary()
	if err != nil {
		return ""
	}
	return string(b)
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

func (res *GetByUsernameResponse) MarshalBinary() ([]byte, error) {
	return json.Marshal(res)
}

func (res *GetByUsernameResponse) String() string {
	b, err := res.MarshalBinary()
	if err != nil {
		return ""
	}
	return string(b)
}
