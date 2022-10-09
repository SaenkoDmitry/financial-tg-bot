package expenses

import (
	"bytes"
	"fmt"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
)

func Format(result map[string]decimal.Decimal, period, currency string) string {
	var formatted bytes.Buffer
	formatted.WriteString(fmt.Sprintf("Расходы за период '%s':\n\n", period))
	if len(result) == 0 {
		formatted.WriteString("Нет трат")
		return formatted.String()
	}
	for i := range constants.CategoryList {
		categoryName := constants.CategoryList[i]
		if amount, ok := result[categoryName]; ok {
			formatted.WriteString(categoryName)
			formatted.WriteString(": ")
			formatted.WriteString(amount.Round(2).String())
			formatted.WriteString(" " + currency)
			formatted.WriteRune('\n')
			formatted.WriteRune('\n')
		}
	}
	return formatted.String()
}
