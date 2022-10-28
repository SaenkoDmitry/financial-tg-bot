package service

import (
	"context"
	"fmt"
	"log"
	"time"

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
	if currency == constants.ServerCurrency {
		return decimal.NewFromInt(1), nil
	}

	// check cache
	key := getCurrencyCacheKey(inputDate)
	if v, ok := s.rateCache.Get(key); ok {
		temp := v.(map[string]decimal.Decimal)
		return temp[currency], nil
	}

	// check db
	res, err := s.rateRepo.GetBatch(ctx, []time.Time{inputDate}, []string{currency})
	if err != nil {
		return decimal.Decimal{}, errors.New(constants.InternalServerErrorMsg)
	}
	saveToCache(s.rateCache, res)
	if v, ok := res[inputDate.Format(dbTimeFormat)]; ok {
		return v[currency], nil
	}

	// load new rates
	var rates map[string]decimal.Decimal
	if inputDate.Format(cacheTimeFormat) == time.Now().Format(cacheTimeFormat) {
		rates, err = s.client.GetLiveCurrency(ctx)
	} else {
		rates, err = s.client.GetHistoricalCurrency(ctx, inputDate)
	}
	if err != nil {
		return decimal.Decimal{}, err
	}

	// persist new rates
	err = s.rateRepo.SaveAll(ctx, rates, inputDate)
	if err != nil {
		log.Printf("cannot save rate to database: %s", err.Error())
	}

	multiplier, ok := rates[currency]
	if !ok {
		return decimal.Decimal{}, errors.New(constants.UndefinedCurrencyMsg)
	}
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
		err := rateCache.Add(key, ratesForDay, defaultExpires)
		if err != nil {
			log.Printf("cannot interact with cache: %s", err.Error())
		}
	}
}

func loadNewRates(ctx context.Context, ratesCache Cache, rateRepo RateStore, client CurrencyExtractor) {
	key := getCurrencyCacheKey(time.Now())
	if _, ok := ratesCache.Get(key); !ok { // first initialization of currencies
		rates, err := client.GetLiveCurrency(ctx)
		if err != nil {
			fmt.Println("cannot set connection to exchange currency api")
			return
		}

		// save to cache
		err = ratesCache.Add(key, rates, defaultExpires)
		if err != nil {
			log.Printf("cannot interact with cache: %s", err.Error())
		}

		// save to db
		err = rateRepo.SaveAll(ctx, rates, time.Now())
		if err != nil {
			log.Printf("cannot save rate to database: %s", err.Error())
		}
	}
}

type Cache interface {
	Get(k string) (interface{}, bool)
	Add(k string, x interface{}, d time.Duration) error
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
				fmt.Println("graceful shutdown")
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
