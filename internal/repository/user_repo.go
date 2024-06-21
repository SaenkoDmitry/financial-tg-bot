package repository

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/logger"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (c *UserRepository) GetUserCurrency(ctx context.Context, userID int64) (currency string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "db:GetUserCurrency")
	defer span.Finish()

	// language=SQL
	sql := `SELECT currency_id FROM financial_bot.user WHERE id = $1`
	span.SetTag("sql", sql)
	row := c.pool.QueryRow(ctx, sql, userID)
	if err = row.Scan(&currency); err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot extract currency", zap.Int64("userID", userID), zap.Error(err))
		return constants.ServerCurrency, err
	}
	return
}

func (c *UserRepository) SetUserCurrency(ctx context.Context, userID int64, newCurrency string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "db:SetUserCurrency")
	defer span.Finish()

	// language=SQL
	sql := `INSERT INTO financial_bot.user (id, currency_id) 
			VALUES ($1, $2) ON CONFLICT (id) 
			DO UPDATE SET currency_id = EXCLUDED.currency_id`
	span.SetTag("sql", sql)
	_, err := c.pool.Exec(ctx, sql, userID, newCurrency)
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot set currency", zap.Int64("userID", userID), zap.Error(err))
		return err
	}
	return nil
}

func (c *UserRepository) GetCurrenciesFilteredByUser(ctx context.Context, userID int64) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "db:GetCurrenciesFilteredByUser")
	defer span.Finish()

	// language=SQL
	sql := `SELECT id FROM financial_bot.currency
			WHERE id NOT IN (SELECT u.currency_id 
			FROM financial_bot.user u where u.id = $1)`
	span.SetTag("sql", sql)
	rows, err := c.pool.Query(ctx, sql, userID)
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot select available currencies", zap.Int64("userID", userID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	currencies := make([]string, 0)
	for rows.Next() {
		var temp string
		err = rows.Scan(&temp)
		if err != nil {
			span.SetTag("error", err.Error())
			logger.Error("cannot scan available currencies", zap.Int64("userID", userID), zap.Error(err))
			return nil, err
		}
		currencies = append(currencies, temp)
	}
	return currencies, nil
}
