package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/shenikar/subscription-service/internal/config"
)

func Connect(cfg config.Config) (*pgx.Conn, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	return pgx.Connect(context.Background(), connStr)
}
