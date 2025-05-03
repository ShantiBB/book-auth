package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"auth/internal/entity"
)

var (
	accessSecret    = []byte("access-secret-key")
	refreshSecret   = []byte("refresh-secret-key")
	accessTokenTTL  = 30 * time.Minute
	refreshTokenTTL = 30 * 24 * time.Hour
)

func GenerateAccessToken(sub int64, role string) (string, error) {
	claims := entity.AccessClaims{
		Sub:  sub,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}

func GenerateRefreshToken(sub int64, role string) (string, error) {
	claims := entity.RefreshClaims{
		Sub:  sub,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

func ParseRefreshToken(tokenStr string) (*entity.RefreshClaims, error) {
	tokenFunc := func(t *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	}
	token, err := jwt.ParseWithClaims(tokenStr, &entity.RefreshClaims{}, tokenFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*entity.RefreshClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
