package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	repoMocks "gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/mocks/repository"
	serviceMocks "gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/mocks/service"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/repository"
	"testing"
	"time"
)

var (
	currentWeekDate  = time.Now().Add(time.Hour * -24 * 2)      // -2 days until now
	currentMonthDate = time.Now().Add(time.Hour * -24 * 7 * 2)  // -14 days until now
	currentYearDate  = time.Now().Add(time.Hour * -24 * 7 * 10) // -10 weeks until now
)

var userID = int64(12345)

var mockTestData = map[int64]map[string][]repository.Transaction{
	userID: {
		constants.Education: []repository.Transaction{
			{
				Amount: decimal.NewFromInt(2000),
				Date:   currentWeekDate,
			},
			{
				Amount: decimal.NewFromInt(5000),
				Date:   currentMonthDate,
			},
		},
		constants.Beauty: []repository.Transaction{
			{
				Amount: decimal.NewFromInt(13000),
				Date:   currentYearDate,
			},
		},
		constants.Clothes: []repository.Transaction{
			{
				Amount: decimal.NewFromInt(2132134),
				Date:   currentMonthDate,
			},
		},
	},
}

func TestFinanceCalculatorService_CalcByCurrentWeek(t *testing.T) {
	ctrl := gomock.NewController(t)

	defaultExpiration := time.Hour * 24 * 30
	cleanupInterval := time.Hour
	simpleCache := NewSimpleCache(context.TODO(), defaultExpiration, cleanupInterval)
	currencyExchangeClient := serviceMocks.NewMockCurrencyExtractor(ctrl)
	currencyExchangeClient.EXPECT().GetLiveCurrency().Times(1)
	transactionRepo := repoMocks.NewMockTransactionOperator(ctrl)
	transactionRepo.EXPECT().GetWallet(userID).Return(mockTestData[userID])
	exchangeRatesService, _ := NewExchangeRatesService(context.Background(), currencyExchangeClient, simpleCache)

	type fields struct {
		m map[int64]map[string][]repository.Transaction
	}
	type args struct {
		userID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]decimal.Decimal
		wantErr bool
	}{
		{
			name: "calc by current week",
			fields: fields{
				m: mockTestData,
			},
			args: args{
				userID: 12345,
			},
			want: map[string]decimal.Decimal{
				constants.Education: decimal.NewFromInt(2000),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FinanceCalculatorService{
				transactionRepo:      transactionRepo,
				exchangeRatesService: exchangeRatesService,
			}
			got, err := f.CalcByCurrentWeek(tt.args.userID, constants.ServerCurrency)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got, "CalcByCurrentWeek: got = %v, want %v", got, tt.want)
		})
	}
}

func TestFinanceCalculatorService_CalcByCurrentMonth(t *testing.T) {
	ctrl := gomock.NewController(t)

	defaultExpiration := time.Hour * 24 * 30
	cleanupInterval := time.Hour
	simpleCache := NewSimpleCache(context.TODO(), defaultExpiration, cleanupInterval)
	currencyExchangeClient := serviceMocks.NewMockCurrencyExtractor(ctrl)
	currencyExchangeClient.EXPECT().GetLiveCurrency().Times(1)
	transactionRepo := repoMocks.NewMockTransactionOperator(ctrl)
	transactionRepo.EXPECT().GetWallet(userID).Return(mockTestData[userID])
	exchangeRatesService, _ := NewExchangeRatesService(context.Background(), currencyExchangeClient, simpleCache)

	type fields struct {
		m map[int64]map[string][]repository.Transaction
	}
	type args struct {
		userID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]decimal.Decimal
		wantErr bool
	}{
		{
			name: "calc by current month",
			fields: fields{
				m: mockTestData,
			},
			args: args{
				userID: 12345,
			},
			want: map[string]decimal.Decimal{
				constants.Education: decimal.NewFromInt(7000),
				constants.Clothes:   decimal.NewFromInt(2132134),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FinanceCalculatorService{
				transactionRepo:      transactionRepo,
				exchangeRatesService: exchangeRatesService,
			}
			got, err := f.CalcByCurrentMonth(tt.args.userID, constants.ServerCurrency)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got, "CalcByCurrentMonth: got = %v, want %v", got, tt.want)
		})
	}
}

func TestFinanceCalculatorService_CalcByCurrentYear(t *testing.T) {
	ctrl := gomock.NewController(t)

	defaultExpiration := time.Hour * 24 * 30
	cleanupInterval := time.Hour
	simpleCache := NewSimpleCache(context.TODO(), defaultExpiration, cleanupInterval)
	currencyExchangeClient := serviceMocks.NewMockCurrencyExtractor(ctrl)
	currencyExchangeClient.EXPECT().GetLiveCurrency().Times(1)
	transactionRepo := repoMocks.NewMockTransactionOperator(ctrl)
	transactionRepo.EXPECT().GetWallet(userID).Return(mockTestData[userID])
	exchangeRatesService, _ := NewExchangeRatesService(context.Background(), currencyExchangeClient, simpleCache)

	type fields struct {
		m map[int64]map[string][]repository.Transaction
	}
	type args struct {
		userID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]decimal.Decimal
		wantErr bool
	}{
		{
			name: "calc by current year",
			fields: fields{
				m: mockTestData,
			},
			args: args{
				userID: 12345,
			},
			want: map[string]decimal.Decimal{
				constants.Education: decimal.NewFromInt(7000),
				constants.Clothes:   decimal.NewFromInt(2132134),
				constants.Beauty:    decimal.NewFromInt(13000),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FinanceCalculatorService{
				transactionRepo:      transactionRepo,
				exchangeRatesService: exchangeRatesService,
			}
			got, err := f.CalcByCurrentYear(tt.args.userID, constants.ServerCurrency)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got, "CalcByCurrentYear: got = %v, want %v", got, tt.want)
		})
	}
}
