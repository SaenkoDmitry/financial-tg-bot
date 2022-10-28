package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/samber/lo"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model"
)

type CategoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{
		pool: pool,
	}
}

func (c CategoryRepository) GetCategories(ctx context.Context) (category []model.CategoryData, err error) {
	// language=SQL
	rows, err := c.pool.Query(ctx, `SELECT id, name_ru FROM financial_bot.category`)
	if err != nil {
		return nil, fmt.Errorf("cannot extract categories from db: %s", err.Error())
	}
	defer rows.Close()
	categories := make([]model.CategoryData, 0)
	for rows.Next() {
		var temp1, temp2 string
		err := rows.Scan(&temp1, &temp2)
		if err != nil {
			return nil, fmt.Errorf("cannot scan categories from db: %s", err.Error())
		}
		categories = append(categories, model.CategoryData{
			ID:   temp1,
			Name: temp2,
		})
	}
	return categories, nil
}

func (c CategoryRepository) ResolveCategories(ctx context.Context, IDs []string) (category map[string]model.CategoryData, err error) {
	categories, err := c.GetCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve categories from db: %s", err.Error())
	}
	filtered := lo.Filter(categories, func(v model.CategoryData, i int) bool {
		return lo.Contains(IDs, v.ID)
	})
	return lo.SliceToMap(filtered, func(t model.CategoryData) (string, model.CategoryData) {
		return t.ID, t
	}), nil
}
