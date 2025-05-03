package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"auth/internal/entity"
)

func GenerateToken(sub int64, role string, ttl time.Duration, secret []byte) (string, error) {
	claims := entity.Claims{
		Sub:  sub,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
