package messages

import (
	"bytes"
	"fmt"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model"
)

type Model struct {
	tgClient MessageSender
}

func New(tgClient MessageSender) *Model {
	return &Model{
		tgClient: tgClient,
	}
}

type Message struct {
	Text   string
	UserID int64
}

func (s *Model) IncomingMessage(msg Message) error {
	switch msg.Text {
	case "/" + constants.Start:
		return s.tgClient.SendMessage(constants.HelloMsg, msg.UserID)
	case "/" + constants.AddOperation:
		return s.tgClient.SendMessageWithMarkup(constants.SpecifyCategoryMsg, collectCategories(), msg.UserID)
	case "/" + constants.ShowCategoryList:
		return s.tgClient.SendMessage(formatCategoryList(constants.CategoryList), msg.UserID)
	case "/" + constants.ShowReport:
		return s.tgClient.SendMessageWithMarkup(constants.SpecifyPeriodMsg, periods, msg.UserID)
	default:
		return s.tgClient.SendMessage(constants.UnrecognizedCommandMsg, msg.UserID)
	}
}

var periods = [][]model.MarkupData{
	{
		{
			Text: constants.WeekPeriod,
			Data: fmt.Sprintf("%s:%s", constants.ShowReport, constants.WeekPeriod),
		},
		{
			Text: constants.MonthPeriod,
			Data: fmt.Sprintf("%s:%s", constants.ShowReport, constants.MonthPeriod),
		},
		{
			Text: constants.YearPeriod,
			Data: fmt.Sprintf("%s:%s", constants.ShowReport, constants.YearPeriod),
		},
	},
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
