package messages

import (
	"bytes"
	"fmt"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/repository"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/utils/slices"
)

type Model struct {
	tgClient         MessageSender
	userCurrencyRepo *repository.UserCurrencyRepository
}

func New(tgClient MessageSender, userCurrencyRepo *repository.UserCurrencyRepository) *Model {
	return &Model{
		tgClient:         tgClient,
		userCurrencyRepo: userCurrencyRepo,
	}
}

type Message struct {
	Text   string
	UserID int64
}

func (s *Model) IncomingMessage(msg Message) error {
	userCurrency, err := s.userCurrencyRepo.GetCurrency(msg.UserID)
	if err != nil || userCurrency == "" {
		userCurrency = constants.RUB
	}
	switch msg.Text {
	case "/" + constants.Start:
		return s.tgClient.SendMessage(constants.HelloMsg, msg.UserID)
	case "/" + constants.AddOperation:
		return s.tgClient.SendMessageWithMarkup(constants.SpecifyCategoryMsg, collectCategories(), msg.UserID)
	case "/" + constants.ShowCategoryList:
		return s.tgClient.SendMessage(formatCategoryList(constants.CategoryList), msg.UserID)
	case "/" + constants.ChangeCurrency:
		return s.tgClient.SendMessageWithMarkup(constants.SpecifyCurrencyMsg, getCurrencies(userCurrency), msg.UserID)
	case "/" + constants.ShowReport:
		return s.tgClient.SendMessageWithMarkup(constants.SpecifyPeriodMsg, getPeriods(), msg.UserID)
	default:
		return s.tgClient.SendMessage(constants.UnrecognizedCommandMsg, msg.UserID)
	}
}

func getCurrencies(userCurrency string) [][]model.MarkupData {
	result := make([][]model.MarkupData, 0, 1)
	filtered := slices.Filter(constants.AllowedCurrencies, userCurrency)
	result = append(result, slices.Map(filtered, func(t string) model.MarkupData {
		return mapToMarkupData(constants.ChangeCurrency, t)
	}))
	return result
}

func getPeriods() [][]model.MarkupData {
	result := make([][]model.MarkupData, 0, 1)
	filtered := []string{constants.WeekPeriod, constants.MonthPeriod, constants.YearPeriod}
	result = append(result, slices.Map(filtered, func(t string) model.MarkupData {
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

func collectCategories() [][]model.MarkupData {
	buttons := make([][]model.MarkupData, 0, len(constants.CategoryList))
	for i := range constants.CategoryList {
		categoryName := constants.CategoryList[i]
		buttons = append(buttons, []model.MarkupData{
			{
				Text: categoryName,
				Data: fmt.Sprintf("%s:%s:", constants.AddOperation, categoryName),
			},
		})
	}
	return buttons
}

func formatCategoryList(categories []string) string {
	var formatted bytes.Buffer
	for i := range categories {
		formatted.WriteString(categories[i])
		formatted.WriteRune('\n')
		formatted.WriteRune('\n')
	}
	return formatted.String()
}
