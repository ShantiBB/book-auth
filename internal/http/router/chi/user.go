package router

import (
	"github.com/go-chi/chi/v5"

	"auth/internal/http/handler"
)

func userRouter(h *handler.Handler) func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", h.CreateUser())
		r.Get("/{id}", h.GetUserByID())
		r.Get("/", h.GetUserAll())
		r.Put("/{id}", h.UpdateUserByID())
		r.Delete("/{id}", h.DeleteUserByID())
	}
}
