package callbacks

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/logger"
	"go.uber.org/zap"
)

func (s *Model) handleSetLimitation(ctx context.Context, query *tgbotapi.CallbackQuery, params ...string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, constants.SetLimitation)
	defer span.Finish()

	input, err := s.parseCategoryWithAmountInputData(ctx, params, query)
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot parse input while setting limit", zap.Error(err))
		return err
	}
	span.SetTag("parse input category", "success")

	if err, needBreak := s.makeProcessOfEnteringAmount(params, input, query); err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot parse amount while adding new operation", zap.Error(err))
		return err
	} else if needBreak {
		return nil
	}
	span.SetTag("parse input amount", "success")

	multiplier, err := s.rateService.GetMultiplier(ctx, input.Currency, time.Now())
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot get multiplier while setting limit", zap.Error(err))
		return s.tgClient.SendEditMessage(fmt.Sprintf(constants.CannotGetRateForYouMsg, constants.ServerCurrency),
			input.UserID, input.MessageID)
	}
	span.SetTag("got multiplier", multiplier.String())

	untilDate := EndOfMonth(time.Now())
	err = s.limitationRepo.AddLimit(ctx, input.UserID, input.CategoryID, input.Amount.Div(multiplier), untilDate) // just overwrite if exists
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot persist data while setting limit", zap.Error(err))
		return err
	}
	span.SetTag("adding limit", "success")

	categories, err := s.categoryRepo.ResolveCategories(ctx, []string{input.CategoryID})
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot resolve categories while setting limit", zap.Error(err))
		return s.tgClient.SendMessage(constants.InternalServerErrorMsg, input.UserID)
	}

	msg := fmt.Sprintf(constants.SetLimitUntilDateMsg,
		categories[input.CategoryID].Name,
		input.Amount.Round(2).String(),
		input.Currency,
		untilDate.Format(untilDateFormat),
	)
	return s.tgClient.SendEditMessage(msg, input.UserID, input.MessageID)
}

var untilDateFormat = "Mon, 02 Jan 2006"

func EndOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 1, -date.Day())
}
