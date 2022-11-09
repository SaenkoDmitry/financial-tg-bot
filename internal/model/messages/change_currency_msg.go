package messages

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/logger"
	"go.uber.org/zap"
)

func (s *Model) changeCurrency(ctx context.Context, msg Message) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, constants.ChangeCurrency)
	defer span.Finish()

	userCurrencies, err := s.userRepo.GetCurrenciesFilteredByUser(ctx, msg.UserID)
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot change user currency", zap.Error(err))
		return s.tgClient.SendMessage(constants.CannotShowCurrencyMenuMsg, msg.UserID)
	}
	return s.tgClient.SendMessageWithMarkup(constants.SpecifyCurrencyMsg, getCurrencies(userCurrencies), msg.UserID)
}
