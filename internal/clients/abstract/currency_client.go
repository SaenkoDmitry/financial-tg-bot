package abstract

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/config"
	"io"
	"net/http"
	"time"
)

type CurrencyClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

type CurrencyExtractor interface {
	GetLiveCurrency() (map[string]float64, error)
	GetHistoricalCurrency(day string) (map[string]float64, error)
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

var (
	CannotCreateReqMsg         = "cannot create request to abstract currency api in method '%s'"
	ErrorWhileMakingReqMsg     = "error while making request to abstract currency api in method '%s'"
	CannotReadResponseMsg      = "cannot read response from abstract currency api in method '%s'"
	CannotUnmarshalResponseMsg = "cannot unmarshal response from abstract currency api in method '%s'"
)

func (s *CurrencyClient) GetLiveCurrency() (map[string]float64, error) {
	method := "v1/live"
	body, err := s.generalCurrencyRequestMaker(method, nil)
	if err != nil {
		return nil, err
	}
	var result *CurrencyLiveResponse
	if err1 := json.Unmarshal(body, &result); err1 != nil || result == nil {
		return nil, errors.Wrap(err, fmt.Sprintf(CannotUnmarshalResponseMsg, method))
	}
	return result.ExchangeRates, nil
}

func (s *CurrencyClient) GetHistoricalCurrency(day string) (map[string]float64, error) {
	method := "v1/historical"
	dateConstraint := fmt.Sprintf("&date=%s", day)
	body, err := s.generalCurrencyRequestMaker(method, &dateConstraint)
	if err != nil {
		return nil, err
	}
	var result *CurrencyHistoricalResponse
	if err1 := json.Unmarshal(body, &result); err1 != nil || result == nil {
		return nil, errors.Wrap(err, fmt.Sprintf(CannotUnmarshalResponseMsg, method))
	}
	return result.ExchangeRates, nil
}

func (s *CurrencyClient) generalCurrencyRequestMaker(method string, constraint *string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s?api_key=%s&base=RUB&target=USD,EUR,CNY", s.baseURL, method, s.apiKey)
	if constraint != nil {
		url += *constraint
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(CannotCreateReqMsg, method))
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrorWhileMakingReqMsg, method))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(CannotReadResponseMsg, method))
	}
	return body, nil
}

type CurrencyLiveResponse struct {
	Base          string             `json:"base"`
	LastUpdated   int64              `json:"last_updated"`
	ExchangeRates map[string]float64 `json:"exchange_rates"`
}

type CurrencyHistoricalResponse struct {
	Base          string             `json:"base"`
	Date          string             `json:"date"`
	ExchangeRates map[string]float64 `json:"exchange_rates"`
}
