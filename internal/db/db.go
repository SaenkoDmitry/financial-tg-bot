package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/config"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/logger"
	"go.uber.org/zap"
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
		logger.Error("unable to init database connection pool", zap.String("url", dsn), zap.Error(err))
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		logger.Error("unable to create database connection pool", zap.String("url", dsn), zap.Error(err))
		return nil, err
	}

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	db := stdlib.OpenDBFromPool(pool)
	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	return pool, nil
}
