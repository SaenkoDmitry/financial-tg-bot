package service

import (
	"context"
	"encoding/json"
	"time"

	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/utils"

	"github.com/opentracing/opentracing-go"

	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/logger"
	"go.uber.org/zap"

	"github.com/shopspring/decimal"
)

type TransactionStore interface {
	CalcAmountByPeriod(ctx context.Context, userID int64, moment time.Time, currencyID string) (map[string]decimal.Decimal, error)
}

type CurrencyExchanger interface {
	GetMultiplier(ctx context.Context, currency string, date time.Time) (decimal.Decimal, error)
}

type CalculatorConfig interface {
	CalcCacheDefaultExpiration() time.Duration
}

type calculatorService struct {
	transactionRepo TransactionStore
	rateRepo        RateStore
	rateService     CurrencyExchanger
	reportCache     Cache
	config          CalculatorConfig
}

func NewCalculatorService(config CalculatorConfig, transactionRepo TransactionStore, rateRepo RateStore, rateService CurrencyExchanger, reportCache Cache) *calculatorService {
	return &calculatorService{
		config:          config,
		transactionRepo: transactionRepo,
		rateRepo:        rateRepo,
		rateService:     rateService,
		reportCache:     reportCache,
	}
}

func (c *calculatorService) CalcByCurrentWeek(ctx context.Context, userID int64, currency string) (map[string]decimal.Decimal, error) {
	return c.calcBy(ctx, "CalcByCurrentWeek", userID, 7, currency)
}

func (c *calculatorService) CalcByCurrentMonth(ctx context.Context, userID int64, currency string) (map[string]decimal.Decimal, error) {
	return c.calcBy(ctx, "CalcByCurrentMonth", userID, 30, currency)
}

func (c *calculatorService) CalcByCurrentYear(ctx context.Context, userID int64, currency string) (map[string]decimal.Decimal, error) {
	return c.calcBy(ctx, "CalcByCurrentYear", userID, 365, currency)
}

func (c *calculatorService) CalcSinceStartOfMonth(ctx context.Context, userID int64, currency string, days int64) (map[string]decimal.Decimal, error) {
	return c.calcBy(ctx, "CalcSinceStartOfMonth", userID, days, currency)
}

func (c *calculatorService) calcBy(ctx context.Context, operationName string,
	userID, days int64, currency string) (map[string]decimal.Decimal, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, operationName)
	defer span.Finish()

	cacheKey := utils.GetCalcCacheKey(userID, currency, days)
	if res, ok := c.reportCache.Get(cacheKey); ok {
		var temp map[string]decimal.Decimal
		if err := json.Unmarshal([]byte(res), &temp); err == nil {
			return temp, nil
		}
	}

	momentInThePast := time.Now().Add(-time.Hour * 24 * time.Duration(days))
	if currency != constants.ServerCurrency {
		dates, err := c.rateRepo.GetDatesWithoutRate(ctx, userID, momentInThePast)
		if err != nil {
			span.SetTag("error", err.Error())
			logger.Error("cannot extract all dates without rates for calcBy amount",
				zap.Int64("userID", userID),
				zap.Int64("days", days),
				zap.String("currency", currency),
				zap.Error(err))
			return nil, err
		}
		for i := range dates { // try load new rates and persist if needed
			c.rateService.GetMultiplier(ctx, currency, dates[i]) // nolint
			time.Sleep(time.Second)                              // because API one request per second constraint
		}
	}

	expenses, err := c.transactionRepo.CalcAmountByPeriod(ctx, userID, momentInThePast, currency)
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot get amount from database for period",
			zap.Int64("userID", userID),
			zap.String("currency", currency),
			zap.Time("afterDate", momentInThePast),
			zap.Error(err))
		return nil, err
	}

	if b, err := json.Marshal(expenses); err == nil {
		err2 := c.reportCache.Add(cacheKey, string(b), c.config.CalcCacheDefaultExpiration())
		if err2 != nil {
			logger.Warn("cannot save calculated report to cache while requesting report", zap.Error(err2))
		}
	}
	return expenses, nil
}
