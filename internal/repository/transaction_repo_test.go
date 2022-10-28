package repository

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTransactionRepo(t *testing.T) {
	ctx := context.Background()
	dbContainer, connPool := SetupTestDatabase()
	defer dbContainer.Terminate(ctx) // nolint

	repository := NewTransactionRepository(connPool)
	userID := int64(12345678)

	t.Run("calculation amount by categories", func(t *testing.T) {
		err := repository.AddOperation(ctx, userID, "RESTAURANTS", decimal.NewFromInt(1000),
			time.Date(2022, 10, 27, 0, 0, 0, 0, time.UTC))
		assert.NoError(t, err)

		err = repository.AddOperation(ctx, userID, "RESTAURANTS", decimal.NewFromInt(1580),
			time.Date(2022, 10, 27, 0, 0, 0, 0, time.UTC))
		assert.NoError(t, err)

		err = repository.AddOperation(ctx, userID, "CLOTHES", decimal.NewFromInt(1053),
			time.Date(2022, 10, 27, 0, 0, 0, 0, time.UTC))
		assert.NoError(t, err)

		err = repository.AddOperation(ctx, userID, "MEDICINE", decimal.NewFromInt(15807),
			time.Date(2022, 10, 27, 0, 0, 0, 0, time.UTC))
		assert.NoError(t, err)

		err = repository.AddOperation(ctx, userID, "CLOTHES", decimal.NewFromInt(2107),
			time.Date(2022, 10, 27, 0, 0, 0, 0, time.UTC))
		assert.NoError(t, err)

		expenses, err := repository.CalcAmountByPeriod(ctx, userID, time.Now().Add(-time.Hour*24), "RUB")
		assert.NoError(t, err)
		assert.Equal(t, 3, len(expenses))
		assert.Equal(t, "2580", expenses["RESTAURANTS"].String())
		assert.Equal(t, "3160", expenses["CLOTHES"].String())
		assert.Equal(t, "15807", expenses["MEDICINE"].String())
	})
}
