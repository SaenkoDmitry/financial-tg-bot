package service

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"time"

	"github.com/shopspring/decimal"
)

type TransactionStore interface {
	CalcAmountByPeriod(ctx context.Context, userID int64, moment time.Time, currencyID string) (map[string]decimal.Decimal, error)
}

type CurrencyExchanger interface {
	GetMultiplier(ctx context.Context, currency string, date time.Time) (decimal.Decimal, error)
}

type calculatorService struct {
	transactionRepo TransactionStore
	rateRepo        RateStore
	rateService     CurrencyExchanger
}

func NewCalculatorService(
	transactionRepo TransactionStore,
	rateRepo RateStore,
	rateService CurrencyExchanger,
) *calculatorService {
	return &calculatorService{
		transactionRepo: transactionRepo,
		rateRepo:        rateRepo,
		rateService:     rateService,
	}
}

func (c *calculatorService) CalcByCurrentWeek(ctx context.Context, userID int64, currency string) (map[string]decimal.Decimal, error) {
	return c.calcBy(ctx, userID, 7, currency)
}

func (c *calculatorService) CalcByCurrentMonth(ctx context.Context, userID int64, currency string) (map[string]decimal.Decimal, error) {
	return c.calcBy(ctx, userID, 30, currency)
}

func (c *calculatorService) CalcByCurrentYear(ctx context.Context, userID int64, currency string) (map[string]decimal.Decimal, error) {
	return c.calcBy(ctx, userID, 365, currency)
}

func (c *calculatorService) CalcSinceStartOfMonth(ctx context.Context, userID int64, currency string, days int64) (map[string]decimal.Decimal, error) {
	return c.calcBy(ctx, userID, days, currency)
}

func (c *calculatorService) calcBy(ctx context.Context, userID, days int64, currency string) (map[string]decimal.Decimal, error) {
	momentInThePast := time.Now().Add(-time.Hour * 24 * time.Duration(days))
	if currency != constants.ServerCurrency {
		dates, err := c.rateRepo.GetDatesWithoutRate(ctx, userID, momentInThePast)
		if err != nil {
			return nil, fmt.Errorf("cannot extract transactions by user=%d", userID)
		}
		for i := range dates { // try load new rates and persist if needed
			c.rateService.GetMultiplier(ctx, currency, dates[i]) // nolint
			time.Sleep(time.Second)                              // because API one request per second constraint
		}
	}

	expenses, err := c.transactionRepo.CalcAmountByPeriod(ctx, userID, momentInThePast, currency)
	if err != nil {
		return nil, fmt.Errorf("cannot extract transactions by user=%d", userID)
	}

	return expenses, nil
}
