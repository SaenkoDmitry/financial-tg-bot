package callbacks

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"time"
)

func (s *Model) handleSetLimitation(ctx context.Context, query *tgbotapi.CallbackQuery, params ...string) error {
	input, err := s.parseCategoryWithAmountInputData(ctx, params, query)
	if err != nil {
		return err
	}

	if err, done := s.makeProcessOfEnteringAmount(params, input, query); err != nil || done {
		return err
	}

	multiplier, err := s.rateService.GetMultiplier(ctx, input.Currency, time.Now())
	if err != nil {
		return s.tgClient.SendEditMessage(fmt.Sprintf(constants.CannotGetRateForYouMsg, constants.ServerCurrency),
			input.UserID, input.MessageID)
	}

	untilDate := EndOfMonth(time.Now())
	err = s.limitationRepo.AddLimit(ctx, input.UserID, input.CategoryID, input.Amount.Div(multiplier), untilDate) // just overwrite if exists
	if err != nil {
		return err
	}

	categories, err := s.categoryRepo.ResolveCategories(ctx, []string{input.CategoryID})
	if err != nil {
		return s.tgClient.SendMessage(constants.InternalServerErrorMsg, input.UserID)
	}

	msg := fmt.Sprintf(constants.SetLimitUntilDateMsg,
		categories[input.CategoryID].Name,
		input.Amount.Round(2).String(),
		input.Currency,
		untilDate.Format(untilDateFormat),
	)
	time.Now()
	return s.tgClient.SendEditMessage(msg, input.UserID, input.MessageID)
}

var untilDateFormat = "Mon, 02 Jan 2006"

func EndOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 1, -date.Day())
}
