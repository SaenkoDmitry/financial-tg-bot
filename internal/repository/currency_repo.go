package repository

import (
	"github.com/pkg/errors"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/utils/slices"
)

type CurrencyRepository struct {
	m              map[int64]string // map[userID]default_currency
	currencyGetter CurrencyGetter
}

type CurrencyGetter interface {
	Currencies() []string
}

func NewUserCurrencyRepository(currencyGetter CurrencyGetter) (*CurrencyRepository, error) {
	return &CurrencyRepository{
		m:              make(map[int64]string),
		currencyGetter: currencyGetter,
	}, nil
}

func (c *CurrencyRepository) GetUserCurrency(userID int64) (string, error) {
	if v, ok := c.m[userID]; ok {
		return v, nil
	}
	return "", errors.New("no such user")
}

func (c *CurrencyRepository) SetUserCurrency(userID int64, newCurrency string) error {
	if !slices.Contains(c.currencyGetter.Currencies(), newCurrency) {
		return errors.New(constants.UndefinedCurrencyMsg)
	}
	c.m[userID] = newCurrency
	return nil
}

func (c *CurrencyRepository) GetFilteredByUserCurrencies(userID int64) []string {
	userCurrency, err := c.GetUserCurrency(userID)
	if err != nil || userCurrency == "" {
		userCurrency = constants.ServerCurrency
	}
	return slices.Filter(c.currencyGetter.Currencies(), userCurrency)
}
