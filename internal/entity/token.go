package entity

import "github.com/golang-jwt/jwt/v5"

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	Sub  int64  `json:"sub"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}
