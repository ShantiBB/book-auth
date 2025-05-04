package jwt

import (
	"time"
)

var (
	accessSecret    = []byte("access-secret-key")
	refreshSecret   = []byte("refresh-secret-key")
	accessTokenTTL  = 30 * time.Minute
	refreshTokenTTL = 30 * 24 * time.Hour
)
