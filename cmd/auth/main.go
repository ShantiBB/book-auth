package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"auth/internal/http/handler"
	router "auth/internal/http/router/chi"
	repository "auth/internal/repository/postgres"
	"auth/internal/service"
	"auth/internal/storage/postgres"

	"github.com/go-chi/chi/v5"
)

// @title Auth API
// @version 1.0
// @description This is the authentication API service.

// @host localhost:8085
// @BasePath /

// @schemes http
func main() {
	logHandler := slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	)

	log := slog.New(logHandler)
	log.Info("init logger")

	db, err := postgres.NewPool(log)
	if err != nil {
		os.Exit(1)
	}

	defer func() {
		err = postgres.DBClose(db, log)
		if err != nil {
			os.Exit(1)
		}

		log.Info("success close database")
	}()

	postgresRepos := repository.New(db)
	services := service.New(db, log, postgresRepos)
	handlers := handler.New(db, log, services)

	chiRouter := chi.NewRouter()
	router.New(chiRouter, handlers)

	log.Info("start auth service", "address", "localhost:8085")
	server := &http.Server{
		Addr:         "localhost:8085",
		Handler:      chiRouter,
		ReadTimeout:  4 * time.Second,
		WriteTimeout: 4 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err = server.ListenAndServe(); err != nil {
		log.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	log.Error("server stopped")
}
