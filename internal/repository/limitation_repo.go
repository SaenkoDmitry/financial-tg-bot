package repository

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/logger"
	"go.uber.org/zap"
)

type LimitationRepository struct {
	pool *pgxpool.Pool
}

func NewLimitationRepository(pool *pgxpool.Pool) *LimitationRepository {
	return &LimitationRepository{
		pool: pool,
	}
}

func (l LimitationRepository) CheckLimit(ctx context.Context, userID int64, categoryID string, amount decimal.Decimal) (decimal.Decimal, bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "db:CheckLimit")
	defer span.Finish()

	// language=SQL
	sql := `SELECT upper_border FROM financial_bot.limitation
	        WHERE user_id = $1 AND category_id = $2 AND until_date > now()`
	span.SetTag("sql", sql)

	row := l.pool.QueryRow(ctx, sql, userID, categoryID)
	var diff decimal.Decimal
	if err := row.Scan(&diff); err != nil {
		span.SetTag("error", err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return decimal.Decimal{}, false, nil
		}
		logger.Error("cannot do check limit",
			zap.Int64("userID", userID),
			zap.String("categoryID", categoryID),
			zap.String("amount", amount.String()),
			zap.Error(err))
		return decimal.Decimal{}, false, err
	}
	return amount.Sub(diff), diff.Cmp(amount) < 0, nil
}

func (l LimitationRepository) AddLimit(ctx context.Context, userID int64, categoryID string, upperBorder decimal.Decimal, untilDate time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "db:AddLimit")
	defer span.Finish()

	// language=SQL
	sql := `INSERT INTO financial_bot.limitation (user_id, category_id, upper_border, until_date) 
			VALUES ($1, $2, $3, $4) ON CONFLICT (user_id, category_id)
			DO UPDATE SET upper_border = EXCLUDED.upper_border, until_date = EXCLUDED.until_date RETURNING (id)`
	span.SetTag("sql", sql)

	row := l.pool.QueryRow(ctx, sql, userID, categoryID, upperBorder, untilDate)
	var transactionID int64
	if err := row.Scan(&transactionID); err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot add limitation",
			zap.Int64("userID", userID),
			zap.String("categoryID", categoryID),
			zap.String("upperBorder", upperBorder.String()),
			zap.Error(err))
		return err
	}
	return nil
}
