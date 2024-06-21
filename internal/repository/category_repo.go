package repository

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/logger"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model"
	"go.uber.org/zap"
)

type CategoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{
		pool: pool,
	}
}

func (c CategoryRepository) GetAllCategories(ctx context.Context) (category []model.CategoryData, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "db:GetAllCategories")
	defer span.Finish()

	// language=SQL
	sql := `SELECT id, name_ru FROM financial_bot.category`
	span.SetTag("sql", sql)
	rows, err := c.pool.Query(ctx, sql)
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot extract categories from db", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	categories := make([]model.CategoryData, 0)
	for rows.Next() {
		var temp1, temp2 string
		err = rows.Scan(&temp1, &temp2)
		if err != nil {
			span.SetTag("error", err.Error())
			logger.Error("cannot scan categories from db", zap.Error(err))
			return nil, err
		}
		categories = append(categories, model.CategoryData{
			ID:   temp1,
			Name: temp2,
		})
	}
	return categories, nil
}

func (c CategoryRepository) ResolveCategories(ctx context.Context, IDs []string) (category map[string]model.CategoryData, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "db:ResolveCategories")
	defer span.Finish()

	categories, err := c.GetAllCategories(ctx)
	if err != nil {
		span.SetTag("error", err.Error())
		logger.Error("cannot resolve categories from db", zap.Error(err))
		return nil, err
	}
	filtered := lo.Filter(categories, func(v model.CategoryData, i int) bool {
		return lo.Contains(IDs, v.ID)
	})
	return lo.SliceToMap(filtered, func(t model.CategoryData) (string, model.CategoryData) {
		return t.ID, t
	}), nil
}
