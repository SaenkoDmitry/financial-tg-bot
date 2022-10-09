package service

import (
	"context"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/clients/abstract"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"time"
)

type ExchangeRatesService struct {
	client     abstract.CurrencyExtractor
	ratesCache *cache.Cache
}

type CurrencyExchanger interface {
	GetMultiplier(currency string, date time.Time) (*float64, error)
}

const cacheTimeFormat = "2006-01-02"
const defaultExpires = time.Hour * 24 * 30

func (s ExchangeRatesService) GetMultiplier(currency string, date time.Time) (*float64, error) {
	key := getCurrencyCacheKey(date) // не зависит от валюты, чтобы экономить запросы в сервис запроса курсов валют

	// hit cache
	if v, ok := s.ratesCache.Get(key); ok {
		temp := v.(map[string]float64)
		if v2, ok2 := temp[currency]; ok2 {
			return &v2, nil
		}
	}

	// miss cache
	var currencies map[string]float64
	var err error

	dateStr := date.Format(cacheTimeFormat)
	now := time.Now().Format(cacheTimeFormat)
	if dateStr == now {
		currencies, err = s.client.GetLiveCurrency()
	} else {
		currencies, err = s.client.GetHistoricalCurrency(dateStr)
	}
	if err != nil {
		return nil, err
	}

	// save to cache
	err = s.ratesCache.Add(getCurrencyCacheKey(date), currencies, defaultExpires)
	if err != nil {
		fmt.Println("cannot save to cache")
	}

	if v, ok := currencies[currency]; ok {
		return &v, nil
	} else {
		return nil, errors.New(constants.UndefinedCurrencyMsg)
	}
}

func getCurrencyCacheKey(date time.Time) string {
	return "CURRENCY_" + date.Format(cacheTimeFormat)
}

func loadCurrencies(ratesCache *cache.Cache, client abstract.CurrencyExtractor) {
	key := getCurrencyCacheKey(time.Now())
	if _, ok := ratesCache.Get(key); !ok { // first initialization of currencies
		currencies, err := client.GetLiveCurrency()
		if err != nil {
			panic("cannot set connection to exchange currency api")
		}
		err = ratesCache.Add(key, currencies, defaultExpires)
		if err != nil {
			panic("cannot interact with cache")
		}
	}
}

func NewExchangeRatesService(ctx context.Context, client abstract.CurrencyExtractor) (*ExchangeRatesService, error) {
	ratesCache := cache.New(defaultExpires, time.Hour)

	loadCurrencies(ratesCache, client) // for first run
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("graceful shutdown")
			break
		case <-ticker.C:
			loadCurrencies(ratesCache, client)
		}
	}()

	return &ExchangeRatesService{
		client:     client,
		ratesCache: ratesCache,
	}, nil
}
