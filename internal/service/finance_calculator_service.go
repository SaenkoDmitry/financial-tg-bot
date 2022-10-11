package service

import (
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/repository"
	"time"
)

type FinanceCalculatorService struct {
	transactionRepo      repository.TransactionOperator
	exchangeRatesService *ExchangeRatesService
}

func NewFinanceCalculatorService(transactionRepo repository.TransactionOperator,
	exchangeRatesService *ExchangeRatesService) *FinanceCalculatorService {
	return &FinanceCalculatorService{
		transactionRepo:      transactionRepo,
		exchangeRatesService: exchangeRatesService,
	}
}

func (c *FinanceCalculatorService) CalcByCurrentWeek(userID int64, currency string) (map[string]decimal.Decimal, error) {
	return c.calcBy(userID, 24*7, currency)
}

func (c *FinanceCalculatorService) CalcByCurrentMonth(userID int64, currency string) (map[string]decimal.Decimal, error) {
	return c.calcBy(userID, 24*30, currency)
}

func (c *FinanceCalculatorService) CalcByCurrentYear(userID int64, currency string) (map[string]decimal.Decimal, error) {
	return c.calcBy(userID, 24*365, currency)
}

func (c *FinanceCalculatorService) calcBy(userID, days int64, currency string) (map[string]decimal.Decimal, error) {
	expenses := make(map[string]decimal.Decimal)
	now := time.Now()
	wallet := c.transactionRepo.GetWallet(userID)
	var resErr error
	for category := range wallet {
		for i := range wallet[category] {
			if now.Sub(wallet[category][i].Date).Hours() >= float64(days) {
				continue
			}
			if _, exists := expenses[category]; !exists {
				expenses[category] = decimal.Zero
			}
			multiplier := decimal.NewFromInt(1)
			if currency != constants.ServerCurrency {
				if temp, err := c.exchangeRatesService.GetMultiplier(currency, wallet[category][i].Date); err == nil {
					multiplier = temp
				} else {
					resErr = constants.MissingCurrencyErr
				}
			}
			newAmount := wallet[category][i].Amount.Mul(multiplier)
			expenses[category] = expenses[category].Add(newAmount)
		}
	}
	return expenses, resErr
}
