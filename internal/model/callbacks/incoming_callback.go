package callbacks

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/repository"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/utils/expenses"
	"strings"
)

type Model struct {
	tgClient        CallbackSender
	transactionRepo *repository.TransactionRepository
}

func New(tgClient CallbackSender, transactionRepo *repository.TransactionRepository) *Model {
	return &Model{
		tgClient:        tgClient,
		transactionRepo: transactionRepo,
	}
}

var (
	emptyCallbackErr = errors.New("empty callback data")
)

func (s *Model) HandleIncomingCallback(query *tgbotapi.CallbackQuery) error {
	split := strings.Split(query.Data, ":")
	if len(split) == 0 {
		return emptyCallbackErr
	}
	var err error
	switch split[0] {
	case constants.AddOperation:
		if split[len(split)-1] == "done" {
			err = s.handleAddOperationWithSelectedCategoryAndAmount(query, split[1:]...)
		} else {
			err = s.handleAddOperationWithSelectedCategory(query, split[1:]...)
		}
	case constants.ShowReport:
		err = s.handleShowReport(query, split[1:]...)
	}
	return err
}

func (s *Model) handleAddOperationWithSelectedCategory(query *tgbotapi.CallbackQuery, params ...string) error {
	if len(params) == 0 {
		return emptyCallbackErr
	}
	userMsg := constants.SpecifyAmountMsg + strings.Join(params[1:], "")
	msgID := query.Message.MessageID
	userID := query.From.ID
	markupData := numericKeyboardAccumulator(query.Data)
	if len(params) > 1 {
		return s.tgClient.SendEditMessageWithMarkupAndText(userMsg, markupData, userID, msgID)
	} else {
		return s.tgClient.SendMessageWithMarkup(userMsg, markupData, userID)
	}
}

func (s *Model) handleAddOperationWithSelectedCategoryAndAmount(query *tgbotapi.CallbackQuery, params ...string) error {
	if len(params) == 0 {
		return emptyCallbackErr
	}
	categoryName := params[0]
	amount, err := decimal.NewFromString(params[1])
	if err != nil {
		return s.tgClient.SendMessage(constants.IncorrectAmountClientMsg, query.From.ID)
	}
	transactionAddedText := fmt.Sprintf(constants.TransactionAddedMsg, categoryName, amount.String())
	_ = s.tgClient.SendEditMessage(transactionAddedText, query.From.ID, query.Message.MessageID)
	return s.transactionRepo.AddOperation(query.From.ID, categoryName, amount)
}

func (s *Model) handleShowReport(query *tgbotapi.CallbackQuery, params ...string) error {
	if len(params) == 0 {
		return emptyCallbackErr
	}
	var err error
	var res map[string]decimal.Decimal
	var period string
	switch params[0] {
	case constants.WeekPeriod:
		period = constants.WeekPeriod
		res, err = s.transactionRepo.CalcByCurrentWeek(query.From.ID)
	case constants.MonthPeriod:
		period = constants.MonthPeriod
		res, err = s.transactionRepo.CalcByCurrentMonth(query.From.ID)
	case constants.YearPeriod:
		period = constants.YearPeriod
		res, err = s.transactionRepo.CalcByCurrentYear(query.From.ID)
	}
	if err != nil {
		return err
	}
	return s.tgClient.SendMessage(expenses.Format(res, period), query.From.ID)
}

func numericKeyboardAccumulator(callback string) [][]model.MarkupData {
	return [][]model.MarkupData{
		{
			numericButton("1", callback),
			numericButton("2", callback),
			numericButton("3", callback),
		},
		{
			numericButton("4", callback),
			numericButton("5", callback),
			numericButton("6", callback),
		},
		{
			numericButton("7", callback),
			numericButton("8", callback),
			numericButton("9", callback),
		},
		{
			numericButton(".", callback),
			numericButton("0", callback),
			model.MarkupData{
				Data: callback + ":done",
				Text: "done",
			},
		},
	}
}

func numericButton(text, callback string) model.MarkupData {
	return model.MarkupData{
		Text: text,
		Data: callback + text,
	}
}
