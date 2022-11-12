package service

import (
	"context"
	"encoding/json"
	"time"

	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/metrics"

	"github.com/opentracing/opentracing-go"

	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/logger"
	"go.uber.org/zap"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
)

type RateStore interface {
	GetDatesWithoutRate(ctx context.Context, id int64, startedFrom time.Time) ([]time.Time, error)
	GetBatch(ctx context.Context, dates []time.Time, currencies []string) (rates map[string]map[string]decimal.Decimal, err error)
	SaveAll(ctx context.Context, rates map[string]decimal.Decimal, date time.Time) error
}

type currencyExchangeService struct {
	client    CurrencyExtractor
	rateRepo  RateStore
	rateCache Cache
}

type CurrencyExtractor interface {
	GetLiveCurrency(ctx context.Context) (map[string]decimal.Decimal, error)
	GetHistoricalCurrency(ctx context.Context, day time.Time) (map[string]decimal.Decimal, error)
}

const cacheTimeFormat = "2006-01-02"
const dbTimeFormat = "2006-01-02"
const defaultExpires = time.Hour * 24 * 30

func (s currencyExchangeService) GetMultiplier(ctx context.Context, currency string, inputDate time.Time) (decimal.Decimal, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetMultiplier")
	defer span.Finish()

	if currency == constants.ServerCurrency {
		span.SetTag("result", "returned value=1 for default currency")
		return decimal.NewFromInt(1), nil
	}

	// check cache
	key := getCurrencyCacheKey(inputDate)
	if v, ok := s.rateCache.Get(key); ok {
		metrics.RatesSourceCounter.WithLabelValues(metrics.CacheLabel).Inc()
		metrics.CacheHitCounter.WithLabelValues(metrics.HitLabel).Inc()
		var temp map[string]decimal.Decimal
		err := json.Unmarshal([]byte(v), &temp)
		if err == nil {
			span.SetTag("result", "returned value from cache")
			return temp[currency], nil
		}
		logger.Error("cannot unmarshal extracted value from cache", zap.Error(err))
	}
	metrics.CacheHitCounter.WithLabelValues(metrics.MissLabel).Inc()

	// check db
	res, err := s.rateRepo.GetBatch(ctx, []time.Time{inputDate}, []string{currency})
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot get batch rates from db",
			zap.Time("inputDate", inputDate),
			zap.String("currency", currency),
			zap.Error(err))
		return decimal.Decimal{}, errors.New(constants.InternalServerErrorMsg)
	}
	saveToCache(s.rateCache, res)
	if v, ok := res[inputDate.Format(dbTimeFormat)]; ok {
		metrics.RatesSourceCounter.WithLabelValues(metrics.DBLabel).Inc()
		span.SetTag("result", "extracted value from database")
		return v[currency], nil
	}

	// load new rates
	var rates map[string]decimal.Decimal
	var callType string
	callTypeStatus := "ok"
	if inputDate.Format(cacheTimeFormat) == time.Now().Format(cacheTimeFormat) {
		callType = metrics.LiveCallTypeLabel
		rates, err = s.client.GetLiveCurrency(ctx)
	} else {
		callType = metrics.HistoricalCallTypeLabel
		rates, err = s.client.GetHistoricalCurrency(ctx, inputDate)
	}
	if err != nil {
		callTypeStatus = "error"
		metrics.RatesAPICallCounter.WithLabelValues(callType, callTypeStatus).Inc()
		span.SetTag("error", err.Error())
		logger.Error("cannot get rates from external currency api", zap.Error(err))
		return decimal.Decimal{}, err
	}
	metrics.RatesAPICallCounter.WithLabelValues(callType, callTypeStatus).Inc()

	// persist new rates
	err = s.rateRepo.SaveAll(ctx, rates, inputDate)
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot save loaded rates to database", zap.Error(err))
	}
	metrics.RatesSourceCounter.WithLabelValues(metrics.APILabel).Inc()

	multiplier, ok := rates[currency]
	if !ok {
		span.SetTag("error", err.Error())
		logger.Error("cannot load correct rate for currency", zap.String("currency", currency), zap.Error(err))
		return decimal.Decimal{}, errors.New(constants.UndefinedCurrencyMsg)
	}

	span.SetTag("result", "got value by http request to external rates api")
	return multiplier, nil
}

func getCurrencyCacheKey(date time.Time) string {
	return "CURRENCY_" + date.Format(cacheTimeFormat)
}

func getCurrencyCacheKeyFromStr(date string) string {
	return "CURRENCY_" + date
}

func loadPersistedRates(ctx context.Context, rateCache Cache, rateRepo RateStore) {
	dates := []time.Time{
		time.Now(),
	}
	currencies := []string{"USD", "EUR", "CNY"}
	rates, err := rateRepo.GetBatch(ctx, dates, currencies)
	if err != nil {
		logger.Error("cannot load persisted rates from database", zap.Error(err))
		return
	}

	saveToCache(rateCache, rates)
}

func saveToCache(rateCache Cache, rates map[string]map[string]decimal.Decimal) {
	if len(rates) == 0 {
		return
	}
	for k, ratesForDay := range rates {
		key := getCurrencyCacheKeyFromStr(k)
		err := rateCache.Add(key, toJson(ratesForDay), defaultExpires)
		if err != nil {
			logger.Error("cannot interact with cache", zap.Error(err))
		}
	}
}

func loadNewRates(ctx context.Context, ratesCache Cache, rateRepo RateStore, client CurrencyExtractor) {
	key := getCurrencyCacheKey(time.Now())
	if _, ok := ratesCache.Get(key); !ok { // first initialization of currencies
		rates, err := client.GetLiveCurrency(ctx)
		if err != nil {
			logger.Error("cannot get rates from external currency api", zap.Error(err))
			return
		}

		// save to mem
		err = ratesCache.Add(key, toJson(rates), defaultExpires)
		if err != nil {
			logger.Error("cannot interact with cache", zap.Error(err))
		}

		// save to db
		err = rateRepo.SaveAll(ctx, rates, time.Now())
		if err != nil {
			logger.Error("cannot save loaded rates to database", zap.Error(err))
		}
	}
}

func toJson(input interface{}) string {
	b, err := json.Marshal(input)
	if err != nil {
		logger.Error("cannot marshal to json")
		return ""
	}
	return string(b)
}

type Cache interface {
	Get(k string) (string, bool)
	Add(k string, x string, d time.Duration) error
}

func NewCurrencyExchangeService(ctx context.Context, currencyClient CurrencyExtractor, rateCache Cache,
	rateRepo RateStore) *currencyExchangeService {
	loadPersistedRates(ctx, rateCache, rateRepo)
	loadNewRates(ctx, rateCache, rateRepo, currencyClient) // for first run
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info("graceful shutdown")
				break
			case <-ticker.C:
				loadNewRates(ctx, rateCache, rateRepo, currencyClient)
			}
		}
	}()

	return &currencyExchangeService{
		client:    currencyClient,
		rateRepo:  rateRepo,
		rateCache: rateCache,
	}
}
