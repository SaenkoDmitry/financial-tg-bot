package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCategoryRepository(t *testing.T) {
	ctx := context.Background()
	dbContainer, connPool := SetupTestDatabase()
	defer dbContainer.Terminate(ctx) // nolint

	repository := NewCategoryRepository(connPool)

	t.Run("getting category list", func(t *testing.T) {
		categories, err := repository.GetCategories(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 11, len(categories))
	})

	t.Run("resolving only several categories", func(t *testing.T) {
		categories, err := repository.ResolveCategories(ctx, []string{"RESTAURANTS", "EDUCATION", "MEDICINE"})
		assert.NoError(t, err)
		assert.Equal(t, len(categories), 3)
		assert.Equal(t, "🎓 Образование", categories["EDUCATION"].Name)
	})
}
