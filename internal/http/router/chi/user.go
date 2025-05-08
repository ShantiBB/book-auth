package router

import (
	"github.com/go-chi/chi/v5"

	"auth/internal/http/handler"
	"auth/internal/http/lib/middleware"
)

func userRouter(h *handler.Handler) func(r chi.Router) {
	return func(r chi.Router) {
		r.Use(middleware.Auth)

		r.Post("/", h.CreateUser())
		r.Get("/", h.GetUserAll())
		r.Put("/{id}", h.UpdateUserByID())
		r.Delete("/{id}", h.DeleteUserByID())
		r.Get("/{id}", h.GetUserByID())
		r.Get("/me", h.GetUserMe())
	}
}
