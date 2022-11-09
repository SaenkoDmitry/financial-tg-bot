package messages

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
)

func (s *Model) showReport(ctx context.Context, msg Message) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, constants.ShowReport) // nolint
	defer span.Finish()

	return s.tgClient.SendMessageWithMarkup(constants.SpecifyPeriodMsg, getPeriods(), msg.UserID)
}
