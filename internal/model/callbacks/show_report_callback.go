package callbacks

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/utils/expenses"
)

func (s *Model) handleShowReport(ctx context.Context, query *tgbotapi.CallbackQuery, params ...string) (err error) {
	if len(params) == 0 {
		return emptyCallbackErr
	}

	userID := query.From.ID
	selectedCurrency, _ := s.userRepo.GetUserCurrency(ctx, userID)
	var res map[string]decimal.Decimal
	var period string
	switch params[0] {
	case constants.WeekPeriod:
		period = constants.WeekPeriod
		res, err = s.calcService.CalcByCurrentWeek(ctx, userID, selectedCurrency)
	case constants.MonthPeriod:
		period = constants.MonthPeriod
		res, err = s.calcService.CalcByCurrentMonth(ctx, userID, selectedCurrency)
	case constants.YearPeriod:
		period = constants.YearPeriod
		res, err = s.calcService.CalcByCurrentYear(ctx, userID, selectedCurrency)
	}
	if err != nil {
		return s.tgClient.SendMessage(constants.InternalServerErrorMsg, userID)
	}
	categoryIDs := make([]string, 0)
	for k := range res {
		categoryIDs = append(categoryIDs, k)
	}
	categories, err := s.categoryRepo.ResolveCategories(ctx, categoryIDs)
	if err != nil {
		return s.tgClient.SendMessage(constants.InternalServerErrorMsg, userID)
	}
	if errors.Is(err, constants.MissingCurrencyErr) {
		return s.tgClient.SendMessage(constants.ServerProblemMsg+expenses.Format(nil, categories, period, constants.ServerCurrency), userID)
	}
	return s.tgClient.SendMessage(expenses.Format(res, categories, period, selectedCurrency), userID)
}
