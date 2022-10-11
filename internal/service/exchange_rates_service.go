package service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"time"
)

type ExchangeRatesService struct {
	client     CurrencyExtractor
	ratesCache Cache
}

type CurrencyExtractor interface {
	GetLiveCurrency() (map[string]decimal.Decimal, error)
	GetHistoricalCurrency(day time.Time) (map[string]decimal.Decimal, error)
}

type CurrencyExchanger interface {
	GetMultiplier(currency string, date time.Time) (float64, error)
}

const cacheTimeFormat = "2006-01-02"
const defaultExpires = time.Hour * 24 * 30

func (s ExchangeRatesService) GetMultiplier(currency string, date time.Time) (decimal.Decimal, error) {
	if currency == constants.ServerCurrency {
		return decimal.NewFromInt(1), nil
	}
	key := getCurrencyCacheKey(date) // не зависит от валюты, чтобы экономить запросы в сервис запроса курсов валют

	// hit cache
	if v, ok := s.ratesCache.Get(key); ok {
		temp := v.(map[string]decimal.Decimal)
		if v2, ok2 := temp[currency]; ok2 {
			return v2, nil
		}
	}

	// miss cache
	var currencies map[string]decimal.Decimal
	var err error

	dateStr := date.Format(cacheTimeFormat)
	now := time.Now().Format(cacheTimeFormat)
	if dateStr == now {
		currencies, err = s.client.GetLiveCurrency()
	} else {
		currencies, err = s.client.GetHistoricalCurrency(date)
	}
	if err != nil {
		return decimal.Decimal{}, err
	}

	// save to cache
	err = s.ratesCache.Add(getCurrencyCacheKey(date), currencies, defaultExpires)
	if err != nil {
		fmt.Println("cannot save to cache")
	}

	if v, ok := currencies[currency]; ok {
		return v, nil
	} else {
		return decimal.Decimal{}, errors.New(constants.UndefinedCurrencyMsg)
	}
}

func getCurrencyCacheKey(date time.Time) string {
	return "CURRENCY_" + date.Format(cacheTimeFormat)
}

func loadCurrencies(ratesCache Cache, client CurrencyExtractor) {
	key := getCurrencyCacheKey(time.Now())
	if _, ok := ratesCache.Get(key); !ok { // first initialization of currencies
		currencies, err := client.GetLiveCurrency()
		if err != nil {
			fmt.Println("cannot set connection to exchange currency api")
			return
		}
		// save to cache
		err = ratesCache.Add(key, currencies, defaultExpires)
		if err != nil {
			panic("cannot interact with cache")
		}
	}
}

type Cache interface {
	Get(k string) (interface{}, bool)
	Add(k string, x interface{}, d time.Duration) error
}

func NewExchangeRatesService(ctx context.Context, client CurrencyExtractor, ratesCache Cache) (*ExchangeRatesService, error) {
	loadCurrencies(ratesCache, client) // for first run
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("graceful shutdown")
				break
			case <-ticker.C:
				loadCurrencies(ratesCache, client)
			}
		}
	}()

	return &ExchangeRatesService{
		client:     client,
		ratesCache: ratesCache,
	}, nil
}
