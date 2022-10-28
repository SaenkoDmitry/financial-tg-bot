package expenses

import (
	"bytes"
	"fmt"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model"
	"sort"
)

func Format(result map[string]decimal.Decimal, categoriesMap map[string]model.CategoryData, period, currency string) string {
	var formatted bytes.Buffer
	formatted.WriteString(fmt.Sprintf("Расходы за период '%s':\n\n", period))
	if len(result) == 0 {
		formatted.WriteString("Нет трат")
		return formatted.String()
	}
	categoriesList := make([]string, 0)
	for k := range result {
		categoriesList = append(categoriesList, k)
	}
	sort.Slice(categoriesList, func(i, j int) bool {
		return categoriesMap[categoriesList[i]].Name < categoriesMap[categoriesList[j]].Name
	})

	for _, categoryID := range categoriesList {
		formatted.WriteString(categoriesMap[categoryID].Name)
		formatted.WriteString(": ")
		formatted.WriteString(result[categoryID].Round(2).String())
		formatted.WriteString(" " + currency)
		formatted.WriteRune('\n')
		formatted.WriteRune('\n')
	}
	return formatted.String()
}
