package repository

import (
	"github.com/pkg/errors"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/utils/slices"
)

type UserCurrencyRepository struct {
	m map[int64]string // map[userID]default_currency
}

func NewUserCurrencyRepository() (*UserCurrencyRepository, error) {
	return &UserCurrencyRepository{
		m: make(map[int64]string),
	}, nil
}

func (c *UserCurrencyRepository) GetCurrency(userID int64) (string, error) {
	if v, ok := c.m[userID]; ok {
		return v, nil
	}
	return "", errors.New("no such user")
}

func (c *UserCurrencyRepository) SetCurrency(userID int64, newCurrency string) error {
	if !slices.Contains(constants.AllowedCurrencies, newCurrency) {
		return errors.New(constants.UndefinedCurrencyMsg)
	}
	c.m[userID] = newCurrency
	return nil
}
