package callbacks

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/repository"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/service"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/utils/expenses"
	"strings"
	"time"
)

type Model struct {
	tgClient                 CallbackSender
	transactionRepo          *repository.TransactionRepository
	userCurrencyRepo         *repository.CurrencyRepository
	exchangeRatesService     *service.ExchangeRatesService
	financeCalculatorService *service.FinanceCalculatorService
}

func New(tgClient CallbackSender, transactionRepo *repository.TransactionRepository,
	userCurrencyRepo *repository.CurrencyRepository,
	exchangeRatesService *service.ExchangeRatesService,
	financeCalculatorService *service.FinanceCalculatorService) *Model {
	return &Model{
		tgClient:                 tgClient,
		transactionRepo:          transactionRepo,
		userCurrencyRepo:         userCurrencyRepo,
		exchangeRatesService:     exchangeRatesService,
		financeCalculatorService: financeCalculatorService,
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
	case constants.ChangeCurrency:
		err = s.handleChangeCurrency(query, split[1:]...)
	}
	return err
}

func (s *Model) handleAddOperationWithSelectedCategory(query *tgbotapi.CallbackQuery, params ...string) error {
	if len(params) == 0 {
		return emptyCallbackErr
	}
	userCurrency := constants.ServerCurrency
	if v, err1 := s.userCurrencyRepo.GetUserCurrency(query.From.ID); err1 == nil {
		userCurrency = v
	}
	userMsg := fmt.Sprintf(constants.SpecifyAmountMsg, userCurrency) + strings.Join(params[1:], "")
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
	userID := query.From.ID
	messageID := query.Message.MessageID
	if len(params) == 0 {
		return emptyCallbackErr
	}
	categoryName := params[0]
	userCurrency := constants.ServerCurrency // parse user currency
	if v, err := s.userCurrencyRepo.GetUserCurrency(userID); err == nil && v != constants.ServerCurrency {
		userCurrency = v
	}

	inputAmount, err := decimal.NewFromString(params[1]) // amount in user currency (not equals to server)
	if err != nil {
		return s.tgClient.SendMessage(constants.IncorrectAmountClientMsg, userID)
	}

	amountInServerCurrency := inputAmount // convert to server currency if it is needed
	if multiplier, err := s.exchangeRatesService.GetMultiplier(userCurrency, time.Now()); err == nil {
		amountInServerCurrency = inputAmount.Div(multiplier)
	}

	// sending result message and save data
	transactionAddedText := fmt.Sprintf(constants.TransactionAddedMsg, categoryName, inputAmount.String(), userCurrency)
	_ = s.tgClient.SendEditMessage(transactionAddedText, userID, messageID)
	return s.transactionRepo.AddOperation(userID, categoryName, amountInServerCurrency)
}

func (s *Model) handleChangeCurrency(query *tgbotapi.CallbackQuery, params ...string) error {
	if len(params) == 0 {
		return emptyCallbackErr
	}
	userID := query.From.ID
	messageID := query.Message.MessageID
	err := s.userCurrencyRepo.SetUserCurrency(userID, params[0])
	if err != nil {
		return s.tgClient.SendEditMessage(constants.CannotChangeCurrencyMsg, userID, messageID)
	}
	return s.tgClient.SendEditMessage(fmt.Sprintf(constants.CurrencyChangedSuccessfullyMsg, params[0]), userID, messageID)
}

func (s *Model) handleShowReport(query *tgbotapi.CallbackQuery, params ...string) error {
	if len(params) == 0 {
		return emptyCallbackErr
	}

	userID := query.From.ID
	selectedCurrency, err := s.userCurrencyRepo.GetUserCurrency(userID)
	if err != nil {
		selectedCurrency = constants.ServerCurrency
	}
	var res map[string]decimal.Decimal
	var period string
	switch params[0] {
	case constants.WeekPeriod:
		period = constants.WeekPeriod
		res, err = s.financeCalculatorService.CalcByCurrentWeek(userID, selectedCurrency)
	case constants.MonthPeriod:
		period = constants.MonthPeriod
		res, err = s.financeCalculatorService.CalcByCurrentMonth(userID, selectedCurrency)
	case constants.YearPeriod:
		period = constants.YearPeriod
		res, err = s.financeCalculatorService.CalcByCurrentYear(userID, selectedCurrency)
	}
	if err != nil {
		return err
	}
	return s.tgClient.SendMessage(expenses.Format(res, period, selectedCurrency), userID)
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
