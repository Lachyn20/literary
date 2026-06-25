package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/hemra-siirow/literary/internal/infrastructure/config"
)

func NewPool(cfg *config.Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)
	return pgxpool.New(context.Background(), connString)
}
