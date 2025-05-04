package request

type Register struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,passwd"`
	Email    string `json:"email" validate:"omitempty,email"`
}

type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Refresh struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
