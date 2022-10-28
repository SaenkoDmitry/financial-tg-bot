package repository

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRateRepo(t *testing.T) {
	ctx := context.Background()
	dbContainer, connPool := SetupTestDatabase()
	defer dbContainer.Terminate(ctx) // nolint

	repository := NewRateRepository(connPool)

	t.Run("get batch rates", func(t *testing.T) {
		yesterday := time.Now().Add(-time.Hour * 24)
		today := time.Now()
		err := repository.SaveAll(ctx, map[string]decimal.Decimal{
			"USD": decimal.NewFromFloat(0.009524),
			"EUR": decimal.NewFromFloat(0.009502),
			"CNY": decimal.NewFromFloat(0.068365),
		}, yesterday)
		assert.NoError(t, err)

		err = repository.SaveAll(ctx, map[string]decimal.Decimal{
			"USD": decimal.NewFromFloat(0.009624),
			"EUR": decimal.NewFromFloat(0.009402),
			"CNY": decimal.NewFromFloat(0.069365),
		}, today)
		assert.NoError(t, err)

		res, err := repository.GetBatch(ctx, []time.Time{yesterday}, []string{"EUR", "USD"})
		assert.NoError(t, err)
		assert.Equal(t, len(res), 1)
		assert.Equal(t, map[string]map[string]decimal.Decimal{
			yesterday.Format("2006-01-02"): {
				"USD": decimal.NewFromFloat(0.009524),
				"EUR": decimal.NewFromFloat(0.009502),
			},
		}, res)
	})
}
