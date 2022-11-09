package callbacks

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/logger"
	"go.uber.org/zap"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/utils/expenses"
)

func (s *Model) handleShowReport(ctx context.Context, query *tgbotapi.CallbackQuery, params ...string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, constants.ShowReport)
	defer span.Finish()

	if len(params) == 0 {
		span.SetTag("error", emptyCallbackErr.Error())
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
		span.SetTag("error", err.Error())
		logger.Error("cannot make report",
			zap.Int64("userID", userID),
			zap.String("currency", selectedCurrency),
			zap.String("period", period),
			zap.Error(err))
		return s.tgClient.SendMessage(constants.InternalServerErrorMsg, userID)
	}
	categoryIDs := make([]string, 0)
	for k := range res {
		categoryIDs = append(categoryIDs, k)
	}
	categories, err := s.categoryRepo.ResolveCategories(ctx, categoryIDs)
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot make report because of resolving categories problem",
			zap.Int64("userID", userID),
			zap.String("currency", selectedCurrency),
			zap.String("period", period),
			zap.Error(err))
		return s.tgClient.SendMessage(constants.InternalServerErrorMsg, userID)
	}
	if errors.Is(err, constants.MissingCurrencyErr) {
		span.SetTag("error", err.Error())
		logger.Error("cannot make report and fallback with fallback on default currency",
			zap.Int64("userID", userID),
			zap.String("currency", selectedCurrency),
			zap.String("period", period),
			zap.Error(err))
		return s.tgClient.SendMessage(constants.ServerProblemMsg+expenses.Format(nil, categories, period, constants.ServerCurrency), userID)
	}
	return s.tgClient.SendMessage(expenses.Format(res, categories, period, selectedCurrency), userID)
}
