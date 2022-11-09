package messages

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/metrics"

	"github.com/opentracing/opentracing-go"

	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/logger"
	"go.uber.org/zap"

	"github.com/samber/lo"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model"
)

type UserStore interface {
	SetUserCurrency(ctx context.Context, userID int64, newCurrency string) error
	GetCurrenciesFilteredByUser(ctx context.Context, userID int64) ([]string, error)
}

type CategoryStore interface {
	GetAllCategories(ctx context.Context) (category []model.CategoryData, err error)
}

type MessageSender interface {
	SendMessage(text string, userID int64) error
	SendMessageWithMarkup(text string, markup [][]model.MarkupData, userID int64) error
}

type Model struct {
	tgClient     MessageSender
	userRepo     UserStore
	categoryRepo CategoryStore
}

func New(tgClient MessageSender,
	userRepo UserStore,
	categoryRepo CategoryStore,
) *Model {
	return &Model{
		tgClient:     tgClient,
		userRepo:     userRepo,
		categoryRepo: categoryRepo,
	}
}

type Message struct {
	Text   string
	UserID int64
}

func (s *Model) IncomingMessage(ctx context.Context, msg Message) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IncomingMessage")
	defer span.Finish()

	span.SetTag("command", msg.Text)
	span.SetTag("userID", msg.UserID)

	modelType := "callback"
	operation := "unrecognized"
	status := "ok"
	start := time.Now()
	defer func() {
		tookTime := time.Since(start).Seconds()
		metrics.IncomingRequestsTotalCounter.WithLabelValues(modelType, operation, status).Inc()
		metrics.IncomingRequestsHistogramResponseTime.WithLabelValues(modelType, operation, status).Observe(tookTime)
	}()

	var err error
	switch msg.Text {
	case "/" + constants.Start:
		err = s.start(ctx, msg)
	case "/" + constants.AddOperation:
		err = s.chooseCategory(ctx, msg.UserID, constants.AddOperation)
	case "/" + constants.SetLimitation:
		err = s.chooseCategory(ctx, msg.UserID, constants.SetLimitation)
	case "/" + constants.ShowCategoryList:
		err = s.showCategories(ctx, msg)
	case "/" + constants.ChangeCurrency:
		err = s.changeCurrency(ctx, msg)
	case "/" + constants.ShowReport:
		err = s.showReport(ctx, msg)
	default:
		operation = "unrecognized"
		err = s.tgClient.SendMessage(constants.UnrecognizedCommandMsg, msg.UserID)
	}
	if err != nil {
		status = "error"
	}
	return err
}

func (s *Model) chooseCategory(ctx context.Context, userID int64, operation string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, operation)
	defer span.Finish()

	categories, err := s.categoryRepo.GetAllCategories(ctx)
	if err != nil {
		logger.Error("cannot make choosing category",
			zap.Int64("userID", userID),
			zap.String("operation", operation),
			zap.Error(err))
		return s.tgClient.SendMessage(constants.InternalServerErrorMsg, userID)
	}
	return s.tgClient.SendMessageWithMarkup(constants.SpecifyCategoryMsg, s.collectCategories(categories, operation), userID)
}

func getCurrencies(currencies []string) [][]model.MarkupData {
	result := make([][]model.MarkupData, 0, 1)
	result = append(result, lo.Map(currencies, func(t string, _ int) model.MarkupData {
		return mapToMarkupData(constants.ChangeCurrency, t)
	}))
	return result
}

func getPeriods() [][]model.MarkupData {
	result := make([][]model.MarkupData, 0, 1)
	filtered := []string{constants.WeekPeriod, constants.MonthPeriod, constants.YearPeriod}
	result = append(result, lo.Map(filtered, func(t string, _ int) model.MarkupData {
		return mapToMarkupData(constants.ShowReport, t)
	}))
	return result
}

func mapToMarkupData(callback, input string) model.MarkupData {
	return model.MarkupData{
		Text: input,
		Data: fmt.Sprintf("%s:%s", callback, input),
	}
}

func (s *Model) collectCategories(categories []model.CategoryData, callback string) [][]model.MarkupData {
	buttons := make([][]model.MarkupData, 0, len(categories))
	for i := range categories {
		buttons = append(buttons, []model.MarkupData{
			{
				Text: categories[i].Name,
				Data: fmt.Sprintf("%s:%s:", callback, categories[i].ID),
			},
		})
	}
	return buttons
}

func formatCategoryList(categories []model.CategoryData) string {
	var formatted bytes.Buffer
	for i := range categories {
		formatted.WriteString(categories[i].Name)
		formatted.WriteRune('\n')
		formatted.WriteRune('\n')
	}
	return formatted.String()
}
