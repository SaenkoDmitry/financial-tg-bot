package messages

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/logger"
	"go.uber.org/zap"
)

func (s *Model) start(ctx context.Context, msg Message) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, constants.Start)
	defer span.Finish()

	err := s.userRepo.SetUserCurrency(ctx, msg.UserID, constants.ServerCurrency)
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot set user currency", zap.Error(err))
		return s.tgClient.SendMessage(constants.InternalServerErrorMsg, msg.UserID)
	}
	return s.tgClient.SendMessage(constants.HelloMsg, msg.UserID)
}
