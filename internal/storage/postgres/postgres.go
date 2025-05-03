package postgres

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(log *slog.Logger) (*pgxpool.Pool, error) {
	ctx := context.Background()
	postgresURL := "postgres://postgres:1221@localhost:5432/users?sslmode=disable"

	cfg, err := pgxpool.ParseConfig(postgresURL)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		log.Error("failed to ping database", "error", err)
		return nil, err
	}

	log.Info("success connect to database")
	return pool, nil
}

func DBClose(pool *pgxpool.Pool, log *slog.Logger) error {
	log.Info("closing database connection")
	pool.Close()
	return nil
}
