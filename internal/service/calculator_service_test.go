package service

import (
	"context"
	"testing"
	"time"

	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/utils"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	serviceMocks "gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/mocks/service"
)

const (
	EducationCategoryID = "EDUCATION"
	BeautyCategoryID    = "BEAUTY"
	ClothesCategoryID   = "CLOTHES"
)

type MockConfig struct {
}

func (t *MockConfig) CalcCacheDefaultExpiration() time.Duration {
	return defaultExpires
}

func TestFinanceCalculatorService_CalcByCurrentWeek(t *testing.T) {
	ctrl := gomock.NewController(t)

	ctx := context.Background()
	cfg := &MockConfig{}
	userID := int64(12345)
	weekExpensesExpected := map[string]decimal.Decimal{
		EducationCategoryID: decimal.NewFromInt(2000),
	}
	currencyID := "RUB"
	defaultExpiration := time.Hour * 24 * 30
	cleanupInterval := time.Hour
	simpleCache := NewSimpleCache(ctx, defaultExpiration, cleanupInterval)
	currencyExchangeClientMock := serviceMocks.NewMockCurrencyExtractor(ctrl)
	transactionRepoMock := serviceMocks.NewMockTransactionStore(ctrl)
	rateRepoMock := serviceMocks.NewMockRateStore(ctrl)
	transactionRepoMock.EXPECT().CalcAmountByPeriod(gomock.Any(), userID, gomock.Any(), currencyID).Return(weekExpensesExpected, nil)
	currencyExchangeClientMock.EXPECT().GetLiveCurrency(ctx).Times(1)
	rateRepoMock.EXPECT().GetBatch(ctx, gomock.Any(), gomock.Any()).Times(1)
	rateRepoMock.EXPECT().SaveAll(ctx, gomock.Any(), gomock.Any())
	exchangeRatesService := NewCurrencyExchangeService(ctx, currencyExchangeClientMock, simpleCache, rateRepoMock)
	reportCacheMock := serviceMocks.NewMockCache(ctrl)
	reportCacheMock.EXPECT().Get(utils.GetCalcCacheKey(userID, currencyID, 7))
	reportCacheMock.EXPECT().Add(utils.GetCalcCacheKey(userID, currencyID, 7), "{\"EDUCATION\":\"2000\"}", defaultExpires)

	type args struct {
		userID   int64
		currency string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]decimal.Decimal
		wantErr bool
	}{
		{
			name: "calc by current week",
			args: args{
				userID:   userID,
				currency: constants.ServerCurrency,
			},
			want:    weekExpensesExpected,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewCalculatorService(cfg, transactionRepoMock, rateRepoMock, exchangeRatesService, reportCacheMock)
			got, err := f.CalcByCurrentWeek(ctx, tt.args.userID, tt.args.currency)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got, "CalcByCurrentWeek: got = %v, want %v", got, tt.want)
		})
	}
}

func TestFinanceCalculatorService_CalcByCurrentMonth(t *testing.T) {
	ctrl := gomock.NewController(t)

	ctx := context.Background()
	cfg := &MockConfig{}
	userID := int64(12345)
	currencyID := "RUB"
	monthExpensesExpected := map[string]decimal.Decimal{
		EducationCategoryID: decimal.NewFromInt(7000),
		ClothesCategoryID:   decimal.NewFromInt(2132134),
	}
	defaultExpiration := time.Hour * 24 * 30
	cleanupInterval := time.Hour
	simpleCache := NewSimpleCache(ctx, defaultExpiration, cleanupInterval)
	currencyExchangeClientMock := serviceMocks.NewMockCurrencyExtractor(ctrl)
	transactionRepoMock := serviceMocks.NewMockTransactionStore(ctrl)
	rateRepoMock := serviceMocks.NewMockRateStore(ctrl)
	transactionRepoMock.EXPECT().CalcAmountByPeriod(gomock.Any(), userID, gomock.Any(), currencyID).Return(monthExpensesExpected, nil)
	currencyExchangeClientMock.EXPECT().GetLiveCurrency(ctx).Times(1)
	rateRepoMock.EXPECT().GetBatch(ctx, gomock.Any(), gomock.Any()).Times(1)
	rateRepoMock.EXPECT().SaveAll(ctx, gomock.Any(), gomock.Any())
	exchangeRatesService := NewCurrencyExchangeService(ctx, currencyExchangeClientMock, simpleCache, rateRepoMock)
	reportCacheMock := serviceMocks.NewMockCache(ctrl)
	reportCacheMock.EXPECT().Get(utils.GetCalcCacheKey(userID, currencyID, 30))
	reportCacheMock.EXPECT().Add(utils.GetCalcCacheKey(userID, currencyID, 30),
		"{\"CLOTHES\":\"2132134\",\"EDUCATION\":\"7000\"}", defaultExpires)

	type args struct {
		userID   int64
		currency string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]decimal.Decimal
		wantErr bool
	}{
		{
			name: "calc by current month",
			args: args{
				userID:   userID,
				currency: constants.ServerCurrency,
			},
			want:    monthExpensesExpected,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewCalculatorService(cfg, transactionRepoMock, rateRepoMock, exchangeRatesService, reportCacheMock)
			got, err := f.CalcByCurrentMonth(ctx, tt.args.userID, constants.ServerCurrency)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got, "CalcByCurrentMonth: got = %v, want %v", got, tt.want)
		})
	}
}

func TestFinanceCalculatorService_CalcByCurrentYear(t *testing.T) {
	ctrl := gomock.NewController(t)

	ctx := context.Background()
	cfg := &MockConfig{}
	userID := int64(12345)
	currencyID := "RUB"
	yearExpensesExpected := map[string]decimal.Decimal{
		EducationCategoryID: decimal.NewFromInt(7000),
		ClothesCategoryID:   decimal.NewFromInt(2132134),
		BeautyCategoryID:    decimal.NewFromInt(13000),
	}
	defaultExpiration := time.Hour * 24 * 30
	cleanupInterval := time.Hour
	simpleCache := NewSimpleCache(ctx, defaultExpiration, cleanupInterval)
	currencyExchangeClientMock := serviceMocks.NewMockCurrencyExtractor(ctrl)
	transactionRepoMock := serviceMocks.NewMockTransactionStore(ctrl)
	rateRepoMock := serviceMocks.NewMockRateStore(ctrl)
	transactionRepoMock.EXPECT().CalcAmountByPeriod(gomock.Any(), userID, gomock.Any(), currencyID).Return(yearExpensesExpected, nil)
	currencyExchangeClientMock.EXPECT().GetLiveCurrency(ctx).Times(1)
	rateRepoMock.EXPECT().GetBatch(ctx, gomock.Any(), gomock.Any()).Times(1)
	rateRepoMock.EXPECT().SaveAll(ctx, gomock.Any(), gomock.Any())
	exchangeRatesService := NewCurrencyExchangeService(ctx, currencyExchangeClientMock, simpleCache, rateRepoMock)
	reportCacheMock := serviceMocks.NewMockCache(ctrl)
	reportCacheMock.EXPECT().Get(utils.GetCalcCacheKey(userID, currencyID, 365))
	reportCacheMock.EXPECT().Add(utils.GetCalcCacheKey(userID, currencyID, 365),
		"{\"BEAUTY\":\"13000\",\"CLOTHES\":\"2132134\",\"EDUCATION\":\"7000\"}", defaultExpires)

	type args struct {
		userID   int64
		currency string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]decimal.Decimal
		wantErr bool
	}{
		{
			name: "calc by current year",
			args: args{
				userID:   userID,
				currency: constants.ServerCurrency,
			},
			want:    yearExpensesExpected,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewCalculatorService(cfg, transactionRepoMock, rateRepoMock, exchangeRatesService, reportCacheMock)
			got, err := f.CalcByCurrentYear(ctx, tt.args.userID, constants.ServerCurrency)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got, "CalcByCurrentYear: got = %v, want %v", got, tt.want)
		})
	}
}
