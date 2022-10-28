package repository

import (
	"context"
	"fmt"
	"github.com/lib/pq"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
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
	// language=SQL
	rows, err := r.pool.Query(ctx, `SELECT currency_id, multiplier, on_date 
			FROM financial_bot.rate 
			WHERE currency_id = ANY($1) AND on_date = ANY($2)`,
		pq.Array(currencies), pq.Array(dates))
	if err != nil {
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
	_, err := r.pool.Exec(ctx, `INSERT INTO financial_bot.rate (currency_id, multiplier, on_date) 
			(SELECT 
				unnest($1::TEXT[]) AS currency_id,
				unnest($2::NUMERIC[]) AS multiplier,
				unnest($3::DATE[]) AS on_date
 			)
			ON CONFLICT (currency_id, on_date) DO UPDATE SET multiplier = EXCLUDED.multiplier`,
		pq.Array(currencies), pq.Array(multipliers), pq.Array(dates))
	if err != nil {
		return fmt.Errorf("cannot execute batch save rates: %s", err.Error())
	}
	return nil
}

func (r RateRepository) GetDatesWithoutRate(ctx context.Context, userID int64, startedFrom time.Time) ([]time.Time, error) {
	// language=SQL
	rows, err := r.pool.Query(ctx, `SELECT t.created_at::DATE
    	FROM financial_bot.transaction t
        	JOIN financial_bot.user u ON t.user_id = u.id
        	LEFT JOIN financial_bot.rate r ON t.created_at::DATE = r.on_date AND r.currency_id = u.currency_id
    	WHERE t.user_id = $1 AND t.created_at > $2 AND r.multiplier IS NULL`,
		userID, startedFrom)
	if err != nil {
		return nil, fmt.Errorf("cannot extract dates without rates: %s", err.Error())
	}
	defer rows.Close()
	dates := make([]time.Time, 0)
	for rows.Next() {
		var onDate time.Time
		err = rows.Scan(&onDate)
		if err != nil {
			return nil, fmt.Errorf("cannot extract rates without rates: %s", err.Error())
		}
		dates = append(dates, onDate)
	}
	return dates, nil
}
