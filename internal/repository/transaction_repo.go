package repository

import (
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	Amount decimal.Decimal
	Date   time.Time
}

type TransactionRepository struct {
	m map[int64]map[string][]Transaction
}

func NewTransactionRepository() (*TransactionRepository, error) {
	return &TransactionRepository{
		m: make(map[int64]map[string][]Transaction),
	}, nil
}

func (c *TransactionRepository) AddOperation(userID int64, categoryName string, amount decimal.Decimal) error {
	if _, ok := c.m[userID]; !ok {
		c.m[userID] = make(map[string][]Transaction)
	}
	if len(c.m[userID][categoryName]) == 0 {
		c.m[userID][categoryName] = make([]Transaction, 0)
	}
	c.m[userID][categoryName] = append(c.m[userID][categoryName], Transaction{
		Amount: amount,
		Date:   time.Now(),
	})
	return nil
}

func (c *TransactionRepository) CalcByCurrentWeek(userID int64) (map[string]decimal.Decimal, error) {
	return c.calcBy(userID, 24*7)
}

func (c *TransactionRepository) CalcByCurrentMonth(userID int64) (map[string]decimal.Decimal, error) {
	return c.calcBy(userID, 24*30)
}

func (c *TransactionRepository) CalcByCurrentYear(userID int64) (map[string]decimal.Decimal, error) {
	return c.calcBy(userID, 24*365)
}

func (c *TransactionRepository) calcBy(userID, days int64) (map[string]decimal.Decimal, error) {
	expenses := make(map[string]decimal.Decimal)
	now := time.Now()
	if wallet, ok := c.m[userID]; ok {
		for category := range wallet {
			for i := range wallet[category] {
				if now.Sub(wallet[category][i].Date).Hours() >= float64(days) {
					continue
				}
				if _, ok2 := expenses[category]; !ok2 {
					expenses[category] = decimal.Zero
				}
				expenses[category] = expenses[category].Add(wallet[category][i].Amount)
			}
		}
	}
	return expenses, nil
}
