package authentication

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Email           string `json:"email"`
	PhoneNumber     string `json:"phoneNumber"`
	Username        string `json:"username"`
	FullName        string `json:"fullName"`
	Token           string `json:"token"`
	ExpiredDuration int64  `json:"expiredDuration"`
}
