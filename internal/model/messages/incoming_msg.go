package messages

import (
	"bytes"
	"context"
	"fmt"

	"github.com/samber/lo"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model"
)

type UserStore interface {
	SetUserCurrency(ctx context.Context, userID int64, newCurrency string) error
	GetCurrenciesFilteredByUser(ctx context.Context, userID int64) ([]string, error)
}

type CategoryStore interface {
	GetCategories(ctx context.Context) (category []model.CategoryData, err error)
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
	switch msg.Text {
	case "/" + constants.Start:
		err := s.userRepo.SetUserCurrency(ctx, msg.UserID, constants.ServerCurrency)
		if err != nil {
			return s.tgClient.SendMessage(constants.InternalServerErrorMsg, msg.UserID)
		}
		return s.tgClient.SendMessage(constants.HelloMsg, msg.UserID)
	case "/" + constants.AddOperation:
		return s.chooseCategory(ctx, msg.UserID, constants.AddOperation)
	case "/" + constants.SetLimitation:
		return s.chooseCategory(ctx, msg.UserID, constants.SetLimitation)
	case "/" + constants.ShowCategoryList:
		categories, err := s.categoryRepo.GetCategories(ctx)
		if err != nil {
			return s.tgClient.SendMessage(constants.InternalServerErrorMsg, msg.UserID)
		}
		return s.tgClient.SendMessage(formatCategoryList(categories), msg.UserID)
	case "/" + constants.ChangeCurrency:
		userCurrencies, err := s.userRepo.GetCurrenciesFilteredByUser(ctx, msg.UserID)
		if err != nil {
			return s.tgClient.SendMessage(constants.CannotShowCurrencyMenuMsg, msg.UserID)
		}
		return s.tgClient.SendMessageWithMarkup(constants.SpecifyCurrencyMsg, getCurrencies(userCurrencies), msg.UserID)
	case "/" + constants.ShowReport:
		return s.tgClient.SendMessageWithMarkup(constants.SpecifyPeriodMsg, getPeriods(), msg.UserID)
	default:
		return s.tgClient.SendMessage(constants.UnrecognizedCommandMsg, msg.UserID)
	}
}

func (s *Model) chooseCategory(ctx context.Context, userID int64, operation string) error {
	categories, err := s.categoryRepo.GetCategories(ctx)
	if err != nil {
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
