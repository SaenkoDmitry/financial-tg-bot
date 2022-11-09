package repository

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestLimitationRepo(t *testing.T) {
	ctx := context.Background()
	dbContainer, connPool := SetupTestDatabase()
	defer dbContainer.Terminate(ctx) // nolint

	repository := NewLimitationRepository(connPool)
	userID := int64(123548568)

	t.Run("check that added limit is applied", func(t *testing.T) {
		err := repository.AddLimit(
			ctx,
			userID,
			"EDUCATION",
			decimal.NewFromInt(1000),
			time.Now().Add(time.Hour*24*5),
		)
		assert.NoError(t, err)

		diff, exceed, err := repository.CheckLimit(ctx, userID, "EDUCATION", decimal.NewFromInt(1500))
		assert.NoError(t, err)
		assert.Equal(t, true, exceed)
		assert.Equal(t, diff, decimal.NewFromInt(500))
	})

	t.Run("check limit if no one exists", func(t *testing.T) {
		_, under, err := repository.CheckLimit(ctx, userID, "TRANSPORT", decimal.NewFromInt(1500))
		assert.NoError(t, err)
		assert.Equal(t, false, under)
	})
}
