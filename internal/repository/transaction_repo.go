package repository

import (
	"context"
	"fmt"
	"time"

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
	// language=SQL
	row := c.pool.QueryRow(ctx, `INSERT INTO financial_bot.transaction 
			(user_id, category_id, amount, created_at) 
			VALUES($1, $2, $3, $4) RETURNING id`, userID, categoryID, amount, time.Now())
	var transactionID int64
	err := row.Scan(&transactionID)
	if err != nil {
		return fmt.Errorf("cannot add new operation: %s", err.Error())
	}
	return nil
}

func (c *TransactionRepository) CalcAmountByPeriod(ctx context.Context, userID int64, moment time.Time, currencyID string) (map[string]decimal.Decimal, error) {
	// language=SQL
	rows, err := c.pool.Query(ctx, `SELECT 
    		t.category_id, SUM(t.amount * COALESCE(r.multiplier, 1)) AS amount
    	FROM financial_bot.transaction t
    		LEFT JOIN financial_bot.rate r on t.created_at::date = r.on_date AND r.currency_id = $3
			WHERE t.user_id = $1 AND t.created_at > $2
    		GROUP BY t.category_id`,
		userID, moment, currencyID)
	if err != nil {
		return nil, fmt.Errorf("cannot extract transactions by user=%d", userID)
	}
	defer rows.Close()
	expenses := make(map[string]decimal.Decimal)
	for rows.Next() {
		var categoryID string
		var amount decimal.Decimal
		err = rows.Scan(&categoryID, &amount)
		if err != nil {
			return nil, fmt.Errorf("cannot scan transactions by user=%d", userID)
		}
		expenses[categoryID] = amount
	}
	return expenses, nil
}
