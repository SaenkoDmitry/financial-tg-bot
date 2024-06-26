package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/lib/pq"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/logger"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type RateRepository struct {
	pool *pgxpool.Pool
}

func NewRateRepository(pool *pgxpool.Pool) *RateRepository {
	return &RateRepository{
		pool: pool,
	}
}

const onDateTimeFormat = "2006-01-02"

func (r RateRepository) GetBatch(ctx context.Context, dates []time.Time,
	currencies []string) (rates map[string]map[string]decimal.Decimal, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "db:GetBatch")
	defer span.Finish()

	// language=SQL
	sql := `SELECT currency_id, multiplier, on_date 
			FROM financial_bot.rate 
			WHERE currency_id = ANY($1) AND on_date = ANY($2)`
	span.SetTag("sql", sql)
	rows, err := r.pool.Query(ctx, sql, pq.Array(currencies), pq.Array(dates))
	if err != nil {
		span.SetTag("error", err.Error())
		return nil, fmt.Errorf("cannot extract rates by currencies=%v", currencies)
	}
	defer rows.Close()
	rates = make(map[string]map[string]decimal.Decimal)
	for rows.Next() {
		var currencyID string
		var onDate time.Time
		var multiplier decimal.Decimal
		err = rows.Scan(&currencyID, &multiplier, &onDate)
		if err != nil {
			span.SetTag("error", err.Error())
			return nil, fmt.Errorf("cannot extract rates batch: %s", err.Error())
		}
		onDateStr := onDate.Format(onDateTimeFormat)
		if _, ok := rates[onDateStr]; !ok {
			rates[onDateStr] = make(map[string]decimal.Decimal)
		}
		rates[onDateStr][currencyID] = multiplier
	}
	return rates, nil
}

func (r RateRepository) SaveAll(ctx context.Context, rates map[string]decimal.Decimal, inputDate time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "db:SaveAll")
	defer span.Finish()

	if len(rates) == 0 {
		return nil
	}
	currencies := make([]string, 0, len(rates))
	multipliers := make([]decimal.Decimal, 0, len(rates))
	dates := make([]time.Time, 0, len(rates))
	for k, v := range rates {
		currencies = append(currencies, k)
		multipliers = append(multipliers, v)
		dates = append(dates, inputDate)
	}
	// language=SQL
	sql := `INSERT INTO financial_bot.rate (currency_id, multiplier, on_date) 
			(SELECT 
				unnest($1::TEXT[]) AS currency_id,
				unnest($2::NUMERIC[]) AS multiplier,
				unnest($3::DATE[]) AS on_date
 			)
			ON CONFLICT (currency_id, on_date) DO UPDATE SET multiplier = EXCLUDED.multiplier`
	span.SetTag("sql", sql)
	_, err := r.pool.Exec(ctx, sql, pq.Array(currencies), pq.Array(multipliers), pq.Array(dates))
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot execute batch save rates",
			zap.Time("inputDate", inputDate),
			zap.Error(err))
		return err
	}
	return nil
}

func (r RateRepository) GetDatesWithoutRate(ctx context.Context, userID int64, startedFrom time.Time) ([]time.Time, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "db:GetDatesWithoutRate")
	defer span.Finish()

	// language=SQL
	sql := `SELECT t.created_at::DATE
    	FROM financial_bot.transaction t
        	JOIN financial_bot.user u ON t.user_id = u.id
        	LEFT JOIN financial_bot.rate r ON t.created_at::DATE = r.on_date AND r.currency_id = u.currency_id
    	WHERE t.user_id = $1 AND t.created_at > $2 AND r.multiplier IS NULL`
	span.SetTag("sql", sql)
	rows, err := r.pool.Query(ctx, sql, userID, startedFrom)
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot extract dates without rates",
			zap.Int64("userID", userID),
			zap.Time("startedFrom", startedFrom),
			zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	dates := make([]time.Time, 0)
	for rows.Next() {
		var onDate time.Time
		err = rows.Scan(&onDate)
		if err != nil {
			span.SetTag("error", err.Error())
			logger.Error("cannot extract rates without rates",
				zap.Int64("userID", userID),
				zap.Time("startedFrom", startedFrom),
				zap.Error(err))
			return nil, err
		}
		dates = append(dates, onDate)
	}
	return dates, nil
}
