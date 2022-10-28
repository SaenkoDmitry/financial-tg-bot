package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	Amount     decimal.Decimal
	CategoryID string
	Date       time.Time
}
