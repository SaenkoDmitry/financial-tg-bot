package utils

import (
	"fmt"
	"time"
)

func GetCalcCacheKey(userID int64, currency string, days int64) string {
	date := time.Now().Format(nowDateFormat)
	return fmt.Sprintf("CALC_%s_%d_%s_%d", date, userID, currency, days)
}

var nowDateFormat = "2006-01-02"
