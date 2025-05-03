package service

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	db   *pgxpool.Pool
	log  *slog.Logger
	repo Repository
}

type Repository interface {
	UserRepository
	TokenRepository
}

func New(db *pgxpool.Pool, log *slog.Logger, repo Repository) *Service {
	return &Service{db, log, repo}
}
