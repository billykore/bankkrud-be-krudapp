package authentication

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=2,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type LoginResponse struct {
	Email           string `json:"email"`
	PhoneNumber     string `json:"phone_number"`
	Username        string `json:"username"`
	FullName        string `json:"full_name"`
	Token           string `json:"token"`
	ExpiredDuration int64  `json:"expired_duration"`
}
