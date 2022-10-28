package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"time"
)

type LimitationRepository struct {
	pool *pgxpool.Pool
}

const (
	// language=SQL
	checkLimitationSQL = `SELECT upper_border FROM financial_bot.limitation
	                    WHERE user_id = $1 AND category_id = $2
	                      AND until_date > now()`
	// language=SQL
	addLimitationSQL = `INSERT INTO financial_bot.limitation (user_id, category_id, upper_border, until_date) 
			VALUES ($1, $2, $3, $4) ON CONFLICT (user_id, category_id)
			DO UPDATE SET upper_border = EXCLUDED.upper_border, until_date = EXCLUDED.until_date RETURNING (id)`
)

func NewLimitationRepository(pool *pgxpool.Pool) *LimitationRepository {
	return &LimitationRepository{
		pool: pool,
	}
}

func (l LimitationRepository) CheckLimit(ctx context.Context, userID int64, categoryID string, amount decimal.Decimal) (decimal.Decimal, bool, error) {
	row := l.pool.QueryRow(ctx, checkLimitationSQL, userID, categoryID)
	var diff decimal.Decimal
	if err := row.Scan(&diff); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return decimal.Decimal{}, false, nil
		}
		return decimal.Decimal{}, false, fmt.Errorf("cannot do check limit: %s", err.Error())
	}
	return amount.Sub(diff), diff.Cmp(amount) < 0, nil
}

func (l LimitationRepository) AddLimit(ctx context.Context, userID int64, categoryID string, upperBorder decimal.Decimal, untilDate time.Time) error {
	row := l.pool.QueryRow(ctx, addLimitationSQL, userID, categoryID, upperBorder, untilDate)
	var transactionID int64
	if err := row.Scan(&transactionID); err != nil {
		return fmt.Errorf("cannot add limitation: %s", err.Error())
	}
	return nil
}
