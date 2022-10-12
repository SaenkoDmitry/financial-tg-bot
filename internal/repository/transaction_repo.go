package repository

import (
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	Amount decimal.Decimal
	Date   time.Time
}

type TransactionOperator interface {
	AddOperation(userID int64, categoryName string, amount decimal.Decimal) error
	GetWallet(userID int64) map[string][]Transaction
}

type TransactionRepository struct {
	m map[int64]map[string][]Transaction // map[userID]map[category][]transaction
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

func (c *TransactionRepository) GetWallet(userID int64) map[string][]Transaction {
	if v, ok := c.m[userID]; ok {
		return v
	}
	return make(map[string][]Transaction)
}
