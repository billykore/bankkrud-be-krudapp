package authentication

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=2,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type LoginResponse struct {
	Username        string `json:"username"`
	Token           string `json:"token"`
	ExpiredDuration int64  `json:"expired_duration"`
}

type LogoutRequest struct {
	Username string `param:"username" json:"username" validate:"required"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}
