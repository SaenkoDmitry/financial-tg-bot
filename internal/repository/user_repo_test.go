package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository(t *testing.T) {
	ctx := context.Background()
	dbContainer, connPool := SetupTestDatabase()
	defer dbContainer.Terminate(ctx) // nolint

	repository := NewUserRepository(connPool)

	t.Run("getting user information", func(t *testing.T) {
		currency, err := repository.GetUserCurrency(ctx, 123548568)
		assert.NoError(t, err)
		assert.Equal(t, "USD", currency)
	})

	t.Run("changing user currency", func(t *testing.T) {
		err := repository.SetUserCurrency(ctx, 123548568, "RUB")
		assert.NoError(t, err)

		currency, err := repository.GetUserCurrency(ctx, 123548568)
		assert.NoError(t, err)
		assert.Equal(t, "RUB", currency)
	})

	t.Run("getting currencies of user without selected", func(t *testing.T) {
		err := repository.SetUserCurrency(ctx, 123548568, "EUR")
		assert.NoError(t, err)
		currencies, err := repository.GetCurrenciesFilteredByUser(ctx, 123548568)
		assert.NoError(t, err)
		assert.Equal(t, []string{"RUB", "USD", "CNY"}, currencies)
	})
}
