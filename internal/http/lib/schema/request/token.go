package request

type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Refresh struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
