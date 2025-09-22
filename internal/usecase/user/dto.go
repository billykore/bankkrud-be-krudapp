package user

type GetByUsernameRequest struct {
	Username string `json:"username" param:"username"`
}

type GetByUsernameResponse struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	CIF         string `json:"cif"`
	Address     string `json:"address"`
}
