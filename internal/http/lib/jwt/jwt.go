package jwt

import (
	"time"

	"auth/internal/entity"
)

var (
	accessSecret    = []byte("access-secret-key")
	refreshSecret   = []byte("refresh-secret-key")
	accessTokenTTL  = 30 * time.Minute
	refreshTokenTTL = 30 * 24 * time.Hour
)

func GenerateAccessToken(sub int64, role string) (string, error) {
	return GenerateToken(sub, role, accessTokenTTL, accessSecret)
}

func GenerateRefreshToken(sub int64, role string) (string, error) {
	return GenerateToken(sub, role, refreshTokenTTL, refreshSecret)
}

func GetClaimsAccessToken(tokenStr string) (*entity.Claims, error) {
	return parseToken(tokenStr, accessSecret)
}

func GetClaimsRefreshToken(tokenStr string) (*entity.Claims, error) {
	return parseToken(tokenStr, refreshSecret)
}
