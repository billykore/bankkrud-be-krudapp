package authentication

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Email           string `json:"email"`
	PhoneNumber     string `json:"phone_number"`
	Username        string `json:"username"`
	FullName        string `json:"full_name"`
	Token           string `json:"token"`
	ExpiredDuration int64  `json:"expired_duration"`
}
