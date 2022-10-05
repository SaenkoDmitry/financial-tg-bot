package repository

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"testing"
	"time"
)

var (
	currentWeekDate  = time.Now().Add(time.Hour * -24 * 2)      // -2 days until now
	currentMonthDate = time.Now().Add(time.Hour * -24 * 7 * 2)  // -14 days until now
	currentYearDate  = time.Now().Add(time.Hour * -24 * 7 * 10) // -10 weeks until now
)

var mockTestData = map[int64]map[string][]Transaction{
	12345: {
		constants.Education: []Transaction{
			{
				Amount: decimal.NewFromInt(2000),
				Date:   currentWeekDate,
			},
			{
				Amount: decimal.NewFromInt(5000),
				Date:   currentMonthDate,
			},
		},
		constants.Beauty: []Transaction{
			{
				Amount: decimal.NewFromInt(13000),
				Date:   currentYearDate,
			},
		},
		constants.Clothes: []Transaction{
			{
				Amount: decimal.NewFromInt(2132134),
				Date:   currentMonthDate,
			},
		},
	},
}

func TestTransactionRepository_CalcByCurrentWeek(t *testing.T) {
	type fields struct {
		m map[int64]map[string][]Transaction
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
			c := &TransactionRepository{
				m: tt.fields.m,
			}
			got, err := c.CalcByCurrentWeek(tt.args.userID)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got, "CalcByCurrentWeek: got = %v, want %v", got, tt.want)
		})
	}
}

func TestTransactionRepository_CalcByCurrentMonth(t *testing.T) {
	type fields struct {
		m map[int64]map[string][]Transaction
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
			c := &TransactionRepository{
				m: tt.fields.m,
			}
			got, err := c.CalcByCurrentMonth(tt.args.userID)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got, "CalcByCurrentMonth: got = %v, want %v", got, tt.want)
		})
	}
}

func TestTransactionRepository_CalcByCurrentYear(t *testing.T) {
	type fields struct {
		m map[int64]map[string][]Transaction
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
				constants.Beauty:    decimal.NewFromInt(13000),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &TransactionRepository{
				m: tt.fields.m,
			}
			got, err := c.CalcByCurrentYear(tt.args.userID)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got, "CalcByCurrentYear: got = %v, want %v", got, tt.want)
		})
	}
}
