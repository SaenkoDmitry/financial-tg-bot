package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
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
	// language=SQL
	row := c.pool.QueryRow(ctx, `SELECT currency_id FROM financial_bot.user WHERE id = $1`, userID)
	if err := row.Scan(&currency); err != nil {
		return constants.ServerCurrency, fmt.Errorf("no such user=%d", userID)
	}
	return
}

func (c *UserRepository) SetUserCurrency(ctx context.Context, userID int64, newCurrency string) error {
	// language=SQL
	_, err := c.pool.Exec(ctx, `INSERT INTO financial_bot.user (id, currency_id) 
			VALUES ($1, $2) ON CONFLICT (id) 
			DO UPDATE SET currency_id = EXCLUDED.currency_id`, userID, newCurrency)
	if err != nil {
		return err
	}
	return nil
}

func (c *UserRepository) GetCurrenciesFilteredByUser(ctx context.Context, userID int64) ([]string, error) {
	// language=SQL
	rows, err := c.pool.Query(ctx, `SELECT id FROM financial_bot.currency
			WHERE id NOT IN (SELECT u.currency_id 
			FROM financial_bot.user u where u.id = $1)`, userID)
	if err != nil {
		return nil, fmt.Errorf("cannot select currencies by user=%d", userID)
	}
	defer rows.Close()
	currencies := make([]string, 0)
	for rows.Next() {
		var temp string
		err := rows.Scan(&temp)
		if err != nil {
			return nil, fmt.Errorf("error while scan currencies by user=%d", userID)
		}
		currencies = append(currencies, temp)
	}
	return currencies, nil
}
