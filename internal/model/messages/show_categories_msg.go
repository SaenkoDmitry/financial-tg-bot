package messages

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/logger"
	"go.uber.org/zap"
)

func (s *Model) showCategories(ctx context.Context, msg Message) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, constants.ShowCategoryList)
	defer span.Finish()

	categories, err := s.categoryRepo.GetAllCategories(ctx)
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot get categories", zap.Error(err))
		return s.tgClient.SendMessage(constants.InternalServerErrorMsg, msg.UserID)
	}
	return s.tgClient.SendMessage(formatCategoryList(categories), msg.UserID)
}
