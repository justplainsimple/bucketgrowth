package bucketgrowth

import (
	"math"
)

// Returns the decimal value that equates to unit/month
// Assumes array in already sorted by time (oldest->newest)
func monthlyGrowth(metrics []DailyMetric) float64 {
	presentDate := metrics[len(metrics)-1].Date
	originalDate := presentDate.AddDate(0, -1, 0)

	idx := 0
	for i, metric := range metrics {
		if metric.Date == originalDate {
			idx = i
			break
		}
	}

	presentValue := metrics[len(metrics)-1].Total
	originalValue := metrics[idx].Total

	growthRate := (presentValue - originalValue) / originalValue
	return float64(growthRate * 100)
}

// Returns the decimal value that equates to unit/year
// Assumes array in already sorted by time (oldest->newest)
func yearlyGrowth(metrics []DailyMetric) float64 {
	// Growth Rate = (Present / Past) * 1/N - 1
	presentValue := metrics[len(metrics)-1].Total
	originalValue := metrics[0].Total

	diff := len(metrics) - 1
	growthRate := (math.Pow(float64(presentValue-originalValue), float64(1/diff)) - 1)

	return float64(growthRate * 100)
}
