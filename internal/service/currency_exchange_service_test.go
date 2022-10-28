package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	serviceMocks "gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/mocks/service"
	"testing"
	"time"
)

func Test_currencyExchangeService_GetMultiplier_CustomCurrency(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	currencyClientMock := serviceMocks.NewMockCurrencyExtractor(ctrl)
	rateRepoMock := serviceMocks.NewMockRateStore(ctrl)
	defaultExpiration := time.Hour * 24 * 30
	cleanupInterval := time.Hour
	simpleCache := NewSimpleCache(ctx, defaultExpiration, cleanupInterval)

	rates := map[string]decimal.Decimal{
		"USD": decimal.NewFromFloat(0.03),
	}
	currencies := []string{"USD", "EUR", "CNY"}
	firstDateStr := "2022-10-18"
	firstDate, _ := time.Parse("2006-01-02", firstDateStr)
	rateUSDValue := decimal.NewFromFloat(0.03)
	currencyClientMock.EXPECT().GetLiveCurrency(ctx).Return(rates, nil).AnyTimes()
	rateRepoMock.EXPECT().SaveAll(ctx, rates, gomock.Any()).Return(nil).AnyTimes()
	rateRepoMock.EXPECT().GetBatch(ctx, gomock.Any(), currencies).Return(
		map[string]map[string]decimal.Decimal{
			firstDateStr: rates,
		}, nil,
	)
	type args struct {
		currency  string
		inputDate time.Time
	}
	tests := []struct {
		name string
		args args
		want decimal.Decimal
	}{
		{
			name: "get multiplier for currency=USD date=2022-10-18",
			args: args{
				currency:  "USD",
				inputDate: firstDate,
			},
			want: rateUSDValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewCurrencyExchangeService(ctx, currencyClientMock, simpleCache, rateRepoMock)
			got, err := s.GetMultiplier(ctx, tt.args.currency, tt.args.inputDate)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got, "GetMultiplier: got = %v, want %v", got, tt.want)
		})
	}
}

func Test_currencyExchangeService_GetMultiplier_ServerCurrency(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	currencyClientMock := serviceMocks.NewMockCurrencyExtractor(ctrl)
	rateRepoMock := serviceMocks.NewMockRateStore(ctrl)
	defaultExpiration := time.Hour * 24 * 30
	cleanupInterval := time.Hour
	simpleCache := NewSimpleCache(ctx, defaultExpiration, cleanupInterval)

	firstDateStr := "2022-10-18"
	firstDate, _ := time.Parse("2006-01-02", firstDateStr)
	rateRUBValue := decimal.NewFromInt(1)
	currencyClientMock.EXPECT().GetLiveCurrency(ctx).AnyTimes()
	rateRepoMock.EXPECT().SaveAll(ctx, gomock.Any(), gomock.Any()).AnyTimes()
	rateRepoMock.EXPECT().GetBatch(gomock.Any(), gomock.Any(), gomock.Any()).Return(
		map[string]map[string]decimal.Decimal{
			firstDateStr: {
				"CNY": decimal.NewFromFloat(0.3),
			},
		}, nil,
	)
	type args struct {
		currency  string
		inputDate time.Time
	}
	tests := []struct {
		name string
		args args
		want decimal.Decimal
	}{
		{
			name: "get multiplier for currency=RUB date=2022-10-18",
			args: args{
				currency:  "RUB",
				inputDate: firstDate,
			},
			want: rateRUBValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewCurrencyExchangeService(ctx, currencyClientMock, simpleCache, rateRepoMock)
			got, err := s.GetMultiplier(ctx, tt.args.currency, tt.args.inputDate)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got, "GetMultiplier: got = %v, want %v", got, tt.want)
		})
	}
}
