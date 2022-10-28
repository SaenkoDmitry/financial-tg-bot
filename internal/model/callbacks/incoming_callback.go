package callbacks

import (
	"context"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model"
)

type UserStore interface {
	GetUserCurrency(ctx context.Context, userID int64) (currency string, err error)
	SetUserCurrency(ctx context.Context, userID int64, newCurrency string) error
}

type CategoryStore interface {
	ResolveCategories(ctx context.Context, IDs []string) (category map[string]model.CategoryData, err error)
}

type TransactionStore interface {
	AddOperation(ctx context.Context, userID int64, categoryID string, amount decimal.Decimal, createdAt time.Time) error
}

type LimitationRepo interface {
	CheckLimit(ctx context.Context, userID int64, categoryID string, amount decimal.Decimal) (decimal.Decimal, bool, error)
	AddLimit(ctx context.Context, userID int64, categoryID string, upperBorder decimal.Decimal, untilDate time.Time) error
}

type CurrencyExchanger interface {
	GetMultiplier(ctx context.Context, currency string, date time.Time) (decimal.Decimal, error)
}

type Calculator interface {
	CalcByCurrentWeek(ctx context.Context, userID int64, currency string) (map[string]decimal.Decimal, error)
	CalcByCurrentMonth(ctx context.Context, userID int64, currency string) (map[string]decimal.Decimal, error)
	CalcByCurrentYear(ctx context.Context, userID int64, currency string) (map[string]decimal.Decimal, error)
	CalcSinceStartOfMonth(ctx context.Context, userID int64, currency string, days int64) (map[string]decimal.Decimal, error)
}

type Model struct {
	tgClient        CallbackSender
	transactionRepo TransactionStore
	userRepo        UserStore
	categoryRepo    CategoryStore
	limitationRepo  LimitationRepo
	rateService     CurrencyExchanger
	calcService     Calculator
}

func New(tgClient CallbackSender, transactionRepo TransactionStore, userRepo UserStore, categoryRepo CategoryStore,
	limitationRepo LimitationRepo, rateService CurrencyExchanger, calcService Calculator) *Model {
	return &Model{
		tgClient:        tgClient,
		transactionRepo: transactionRepo,
		categoryRepo:    categoryRepo,
		userRepo:        userRepo,
		limitationRepo:  limitationRepo,
		rateService:     rateService,
		calcService:     calcService,
	}
}

var (
	emptyCallbackErr = errors.New("empty callback data")
)

func (s *Model) HandleIncomingCallback(ctx context.Context, query *tgbotapi.CallbackQuery) error {
	split := strings.Split(query.Data, ":")
	if len(split) == 0 {
		return emptyCallbackErr
	}
	var err error
	switch split[0] {
	case constants.AddOperation:
		err = s.handleAddOperation(ctx, query, split[1:]...)
	case constants.SetLimitation:
		err = s.handleSetLimitation(ctx, query, split[1:]...)
	case constants.ShowReport:
		err = s.handleShowReport(ctx, query, split[1:]...)
	case constants.ChangeCurrency:
		err = s.handleChangeCurrency(ctx, query, split[1:]...)
	}
	return err
}
