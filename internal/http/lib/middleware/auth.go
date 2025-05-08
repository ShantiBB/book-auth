package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/render"

	"auth/internal/http/lib/jwt"
	"auth/internal/http/lib/schema/response"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const bearerPrefix = "Bearer "

		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, response.Error("missing or invalid Authorization header"))
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, bearerPrefix)

		claims, err := jwt.GetClaimsAccessToken(tokenStr)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, response.Error("invalid token"))
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.Sub)
		ctx = context.WithValue(ctx, "userRole", claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
