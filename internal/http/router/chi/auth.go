package router

import (
	"github.com/go-chi/chi/v5"

	"auth/internal/http/handler"
)

func authRouter(h *handler.Handler) func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/login", h.Login())
		r.Post("/refresh", h.Refresh())
	}
}
