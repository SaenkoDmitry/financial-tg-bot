package callbacks

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
)

func (s *Model) handleChangeCurrency(ctx context.Context, query *tgbotapi.CallbackQuery, params ...string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, constants.ChangeCurrency)
	defer span.Finish()

	if len(params) == 0 {
		span.SetTag("error", emptyCallbackErr.Error())
		return emptyCallbackErr
	}
	userID := query.From.ID
	messageID := query.Message.MessageID
	err := s.userRepo.SetUserCurrency(ctx, userID, params[0])
	if err != nil {
		span.SetTag("error", err.Error())
		return s.tgClient.SendEditMessage(constants.CannotChangeCurrencyMsg, userID, messageID)
	}
	return s.tgClient.SendEditMessage(fmt.Sprintf(constants.CurrencyChangedSuccessfullyMsg, params[0]), userID, messageID)
}
