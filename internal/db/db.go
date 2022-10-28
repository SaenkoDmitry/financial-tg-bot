package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/config"
	"log"
)

const (
	DSN = "postgres://%s:%s@%s:%s/%s?sslmode=disable"
)

func getDSN(cfg *config.Service) string {
	user, password, dbName, port, host := cfg.PostgresUser(), cfg.PostgresPassword(),
		cfg.PostgresDB(), cfg.PostgresPort(), cfg.PostgresHost()
	return fmt.Sprintf(DSN, user, password, host, port, dbName)
}

func InitPool(cfg *config.Service) (*pgxpool.Pool, error) {
	dsn := getDSN(cfg)
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Printf("Unable to parse DATABASE_URL=%s with error: %s", dsn, err.Error())
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Printf("Unable to create connection pool: %s", err.Error())
		return nil, err
	}

	return pool, nil
}
