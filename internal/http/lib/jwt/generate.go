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

func GenerateAccessToken(sub int64, role string) (string, error) {
	return GenerateToken(sub, role, accessTokenTTL, accessSecret)
}

func GenerateRefreshToken(sub int64, role string) (string, error) {
	return GenerateToken(sub, role, refreshTokenTTL, refreshSecret)
}

func GenerateAllTokens(sub int64, role string) (*entity.Token, error) {
	var err error
	tokens := &entity.Token{}
	tokens.AccessToken, err = GenerateAccessToken(sub, role)
	if err != nil {
		return nil, err
	}

	tokens.RefreshToken, err = GenerateRefreshToken(sub, role)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
