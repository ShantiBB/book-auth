package handler

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	db  *pgxpool.Pool
	log *slog.Logger
	svc Service
}

type Service interface {
	UserService
	TokenService
}

func New(db *pgxpool.Pool, log *slog.Logger, svc Service) *Handler {
	return &Handler{db, log, svc}
}
