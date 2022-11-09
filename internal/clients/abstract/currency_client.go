package abstract

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/config"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/logger"
	"go.uber.org/zap"
)

type CurrencyClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewCurrencyClient(config *config.Service) *CurrencyClient {
	transport := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 1 * time.Second,
	}
	client := &http.Client{Transport: transport}
	return &CurrencyClient{
		baseURL: "https://exchange-rates.abstractapi.com",
		apiKey:  config.AbstractAPIKey(),
		client:  client,
	}
}

func (s *CurrencyClient) GetLiveCurrency(ctx context.Context) (map[string]decimal.Decimal, error) {
	body, err := s.generalCurrencyRequestMaker(ctx, "v1/live")
	if err != nil {
		logger.Error("error while request api in method GetLiveCurrency", zap.Error(err))
		return nil, err
	}
	var result *CurrencyLiveResponse
	if err1 := json.Unmarshal(body, &result); err1 != nil || result == nil {
		logger.Error("cannot unmarshal response in method GetLiveCurrency", zap.Error(err))
		return nil, err
	}
	return result.ExchangeRates, nil
}

const historicalDateFormat = "2006-01-02"

func (s *CurrencyClient) GetHistoricalCurrency(ctx context.Context, day time.Time) (map[string]decimal.Decimal, error) {
	dateConstraint := fmt.Sprintf("&date=%s", day.Format(historicalDateFormat))
	body, err := s.generalCurrencyRequestMaker(ctx, "v1/historical", dateConstraint)
	if err != nil {
		logger.Error("error while request api in method GetHistoricalCurrency", zap.Error(err))
		return nil, err
	}
	var result *CurrencyHistoricalResponse
	if err1 := json.Unmarshal(body, &result); err1 != nil || result == nil {
		logger.Error("cannot unmarshal response in method GetHistoricalCurrency", zap.Error(err))
		return nil, err
	}
	return result.ExchangeRates, nil
}

func (s *CurrencyClient) generalCurrencyRequestMaker(ctx context.Context, method string, constraints ...string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s?api_key=%s&base=RUB&target=USD,EUR,CNY", s.baseURL, method, s.apiKey)
	if len(constraints) > 0 {
		url = fmt.Sprintf("%s&%s", url, strings.Join(constraints, "&"))
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("cannot create request to abstract currency api in method '%s'", method), zap.Error(err))
		return nil, err
	}

	logger.Debug("outgoing http request", zap.String("url", url), zap.String("method", http.MethodGet))
	resp, err := s.client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("error while making request to abstract currency api in method '%s'", method), zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("cannot read response from abstract currency api in method '%s'", method), zap.Error(err))
		return nil, err
	}
	return body, nil
}

type CurrencyLiveResponse struct {
	Base          string                     `json:"base"`
	LastUpdated   int64                      `json:"last_updated"`
	ExchangeRates map[string]decimal.Decimal `json:"exchange_rates"`
}

type CurrencyHistoricalResponse struct {
	Base          string                     `json:"base"`
	Date          string                     `json:"date"`
	ExchangeRates map[string]decimal.Decimal `json:"exchange_rates"`
}
