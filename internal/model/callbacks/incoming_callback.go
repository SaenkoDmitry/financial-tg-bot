package callbacks

import (
	"context"
	"strings"
	"time"

	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/metrics"

	"github.com/opentracing/opentracing-go"

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

type Cache interface {
	Get(k string) (string, bool)
	Add(k string, x string, d time.Duration) error
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
	reportCache     Cache
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "HandleIncomingCallback")
	defer span.Finish()

	span.SetTag("userID", query.From.ID)
	span.SetTag("messageID", query.Message.MessageID)

	modelType := "callback"
	operation := "unrecognized"
	status := "ok"
	start := time.Now()
	defer func() {
		tookTime := time.Since(start).Seconds()
		metrics.IncomingRequestsTotalCounter.WithLabelValues(modelType, operation, status).Inc()
		metrics.IncomingRequestsHistogramResponseTime.WithLabelValues(modelType, operation, status).Observe(tookTime)
	}()

	split := strings.Split(query.Data, ":")
	if len(split) == 0 {
		status = "error"
		return emptyCallbackErr
	}
	operation = split[0]
	var err error
	switch operation {
	case constants.AddOperation:
		err = s.handleAddOperation(ctx, query, split[1:]...)
	case constants.SetLimitation:
		err = s.handleSetLimitation(ctx, query, split[1:]...)
	case constants.ShowReport:
		err = s.handleShowReport(ctx, query, split[1:]...)
	case constants.ChangeCurrency:
		err = s.handleChangeCurrency(ctx, query, split[1:]...)
	default:
		operation = "unrecognized"
	}
	if err != nil {
		status = "error"
	}
	return err
}
