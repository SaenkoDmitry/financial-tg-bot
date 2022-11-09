package service

import (
	"context"
	"testing"
	"time"

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

func TestFinanceCalculatorService_CalcByCurrentWeek(t *testing.T) {
	ctrl := gomock.NewController(t)

	ctx := context.Background()
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
			f := NewCalculatorService(transactionRepoMock, rateRepoMock, exchangeRatesService)
			got, err := f.CalcByCurrentWeek(ctx, tt.args.userID, tt.args.currency)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got, "CalcByCurrentWeek: got = %v, want %v", got, tt.want)
		})
	}
}

func TestFinanceCalculatorService_CalcByCurrentMonth(t *testing.T) {
	ctrl := gomock.NewController(t)

	ctx := context.Background()
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
			f := NewCalculatorService(transactionRepoMock, rateRepoMock, exchangeRatesService)
			got, err := f.CalcByCurrentMonth(ctx, tt.args.userID, constants.ServerCurrency)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got, "CalcByCurrentMonth: got = %v, want %v", got, tt.want)
		})
	}
}

func TestFinanceCalculatorService_CalcByCurrentYear(t *testing.T) {
	ctrl := gomock.NewController(t)

	ctx := context.Background()
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
			f := NewCalculatorService(transactionRepoMock, rateRepoMock, exchangeRatesService)
			got, err := f.CalcByCurrentYear(ctx, tt.args.userID, constants.ServerCurrency)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got, "CalcByCurrentYear: got = %v, want %v", got, tt.want)
		})
	}
}
