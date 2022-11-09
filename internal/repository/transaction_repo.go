package repository

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"

	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/logger"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model"
)

type TransactionRepository struct {
	m    map[int64]map[string][]model.Transaction // map[userID]map[category][]transaction
	pool *pgxpool.Pool
}

func NewTransactionRepository(pool *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{
		m:    make(map[int64]map[string][]model.Transaction),
		pool: pool,
	}
}

func (c *TransactionRepository) AddOperation(ctx context.Context, userID int64, categoryID string, amount decimal.Decimal, createdAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "db:AddOperation")
	defer span.Finish()

	// language=SQL
	sql := `INSERT INTO financial_bot.transaction 
			(user_id, category_id, amount, created_at) 
			VALUES($1, $2, $3, $4) RETURNING id`
	span.SetTag("sql", sql)
	row := c.pool.QueryRow(ctx, sql, userID, categoryID, amount, time.Now())
	var transactionID int64
	err := row.Scan(&transactionID)
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot add new operation",
			zap.Int64("userID", userID),
			zap.String("categoryID", categoryID),
			zap.Error(err))
		return err
	}
	return nil
}

func (c *TransactionRepository) CalcAmountByPeriod(ctx context.Context, userID int64, moment time.Time, currencyID string) (map[string]decimal.Decimal, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "db:CalcAmountByPeriod")
	defer span.Finish()

	// language=SQL
	sql := `SELECT 
    		t.category_id, SUM(t.amount * COALESCE(r.multiplier, 1)) AS amount
    	FROM financial_bot.transaction t
    		LEFT JOIN financial_bot.rate r on t.created_at::date = r.on_date AND r.currency_id = $3
			WHERE t.user_id = $1 AND t.created_at > $2
    		GROUP BY t.category_id`
	span.SetTag("sql", sql)
	rows, err := c.pool.Query(ctx, sql, userID, moment, currencyID)
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot extract transactions", zap.Int64("userID", userID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	expenses := make(map[string]decimal.Decimal)
	for rows.Next() {
		var categoryID string
		var amount decimal.Decimal
		err = rows.Scan(&categoryID, &amount)
		if err != nil {
			span.SetTag("error", err.Error())
			logger.Error("cannot scan transactions", zap.Int64("userID", userID), zap.Error(err))
			return nil, err
		}
		expenses[categoryID] = amount
	}
	return expenses, nil
}
