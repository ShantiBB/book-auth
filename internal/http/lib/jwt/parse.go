package jwt

import (
	"github.com/golang-jwt/jwt/v5"

	"auth/internal/entity"
)

func parseToken(tokenStr string, secret []byte) (*entity.Claims, error) {
	tokenFunc := func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	}
	token, err := jwt.ParseWithClaims(tokenStr, &entity.Claims{}, tokenFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*entity.Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

func GetClaimsAccessToken(tokenStr string) (*entity.Claims, error) {
	return parseToken(tokenStr, accessSecret)
}

func GetClaimsRefreshToken(tokenStr string) (*entity.Claims, error) {
	return parseToken(tokenStr, refreshSecret)
}
