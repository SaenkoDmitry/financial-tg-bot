package callbacks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/utils"

	"github.com/opentracing/opentracing-go"

	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/logger"
	"go.uber.org/zap"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model"
)

func (s *Model) handleAddOperation(ctx context.Context, query *tgbotapi.CallbackQuery, params ...string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, constants.AddOperation)
	defer span.Finish()

	input, err := s.parseCategoryWithAmountInputData(ctx, params, query)
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot parse input while adding new operation", zap.Error(err))
		return err
	}
	span.SetTag("parse input category", "success")

	if err, needBreak := s.makeProcessOfEnteringAmount(params, input, query); err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot parse amount while adding new operation", zap.Error(err))
		return err
	} else if needBreak {
		return nil
	}
	span.SetTag("parse input amount", "success")

	multiplier, err := s.rateService.GetMultiplier(ctx, input.Currency, time.Now())
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot get multiplier while adding new operation", zap.Error(err))
		return s.tgClient.SendEditMessage(fmt.Sprintf(constants.CannotGetRateForYouMsg, constants.ServerCurrency),
			input.UserID, input.MessageID)
	}
	span.SetTag("got multiplier", multiplier.String())

	if multiplier.IsZero() {
		multiplier = decimal.NewFromInt(1)
	}
	amount := input.Amount.Div(multiplier)

	// resolve categories to display
	categories, err := s.categoryRepo.ResolveCategories(ctx, []string{input.CategoryID})
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot resolve categories while adding new operation", zap.Error(err))
		return s.tgClient.SendMessage(constants.InternalServerErrorMsg, input.UserID)
	}

	// persist data
	err = s.transactionRepo.AddOperation(ctx, input.UserID, input.CategoryID, amount, time.Now())
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot persist data while adding new operation", zap.Error(err))
		return err
	}

	s.deleteCacheValueForPeriod(input, 7)
	s.deleteCacheValueForPeriod(input, 30)
	s.deleteCacheValueForPeriod(input, 365)

	spend, err := s.getSpendSinceStartOfMonth(ctx, input, multiplier)
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot get operations since start of month while adding new operation", zap.Error(err))
		return err
	}

	diff, exceeded, err := s.limitationRepo.CheckLimit(ctx, input.UserID, input.CategoryID, spend)
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot check limit while adding new operation", zap.Error(err))
		return err
	}

	if exceeded {
		amountExceededText := fmt.Sprintf(constants.LimitExceededMsg,
			categories[input.CategoryID].Name,
			input.Amount.Round(2).String(),
			input.Currency, diff.Mul(multiplier).Round(2).String(), input.Currency)
		return s.tgClient.SendEditMessage(amountExceededText, input.UserID, input.MessageID)
	}

	span.SetTag("adding transaction", "success")
	transactionAddedText := fmt.Sprintf(constants.TransactionAddedMsg, categories[input.CategoryID].Name, input.Amount.Round(2).String(), input.Currency)
	return s.tgClient.SendEditMessage(transactionAddedText, input.UserID, input.MessageID)
}

func (s *Model) deleteCacheValueForPeriod(input *addOperationInputData, days int64) {
	key := utils.GetCalcCacheKey(input.UserID, input.Currency, days)
	if err := s.reportCache.Delete(key); err != nil {
		logger.Warn("cannot delete value from cache for period", zap.Error(err))
	}
}

func (s *Model) getSpendSinceStartOfMonth(ctx context.Context, input *addOperationInputData, multiplier decimal.Decimal) (decimal.Decimal, error) {
	spendByCategories, err := s.calcService.CalcSinceStartOfMonth(ctx, input.UserID, input.Currency, int64(time.Now().Day()))
	if err != nil {
		return decimal.Zero, err
	}
	if v, ok := spendByCategories[input.CategoryID]; ok {
		return v.Div(multiplier), nil
	}
	return decimal.Zero, nil
}

func (s *Model) makeProcessOfEnteringAmount(params []string, input *addOperationInputData, query *tgbotapi.CallbackQuery) (error, bool) {
	// process of entering whole amount (accumulation)
	if params[len(params)-1] != "done" {
		userMsg := fmt.Sprintf(constants.SpecifyAmountMsg, input.Currency) + strings.Join(params[1:], "")
		markupData := numericKeyboardAccumulator(query.Data)
		if len(params) > 1 {
			return s.tgClient.SendEditMessageWithMarkupAndText(userMsg, markupData, input.UserID, input.MessageID), true
		} else {
			return s.tgClient.SendMessageWithMarkup(userMsg, markupData, input.UserID), true
		}
	}
	return nil, false
}

type addOperationInputData struct {
	UserID     int64
	MessageID  int
	CategoryID string
	Currency   string
	Amount     decimal.Decimal
}

func (s *Model) parseCategoryWithAmountInputData(ctx context.Context, params []string, query *tgbotapi.CallbackQuery) (*addOperationInputData, error) {
	if len(params) == 0 {
		return nil, emptyCallbackErr
	}
	userID := query.From.ID
	messageID := query.Message.MessageID
	var amount decimal.Decimal
	if len(params) > 2 {
		var err error
		amount, err = decimal.NewFromString(params[1])
		if err != nil {
			return nil, s.tgClient.SendMessage(constants.IncorrectAmountClientMsg, userID)
		}
	}
	return &addOperationInputData{
		UserID:     userID,
		MessageID:  messageID,
		CategoryID: params[0],
		Currency:   s.getUserCurrency(ctx, userID),
		Amount:     amount,
	}, nil
}

func (s *Model) getUserCurrency(ctx context.Context, userID int64) string {
	if v, err := s.userRepo.GetUserCurrency(ctx, userID); err == nil && v != constants.ServerCurrency {
		return v
	}
	return constants.ServerCurrency
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
