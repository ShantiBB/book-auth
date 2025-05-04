package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"auth/internal/http/handler"
)

func New(r chi.Router, h *handler.Handler) {
	r.Use(middleware.CleanPath)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Recoverer)

	r.Route("/auth", authRouter(h))
	r.Route("/users", userRouter(h))
}
