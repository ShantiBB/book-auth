package router

import (
	_ "auth/docs"
	"auth/internal/http/handler"
	localMW "auth/internal/http/lib/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func New(r chi.Router, h *handler.Handler) {
	r.Use(middleware.CleanPath)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Recoverer)
	r.Use(localMW.ContentTypeJSON)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/auth", authRouter(h))
	r.Route("/users", userRouter(h))
}
