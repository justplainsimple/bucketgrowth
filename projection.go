package bucketgrowth

import (
	"math"
)

func projection(val int64, numYears int, monthlyGrowthRate float64) int64 {
	if numYears == 0 {
		return 0.0
	}

	if val == 0 {
		return 0.0
	}

	if monthlyGrowthRate == 0.0 {
		return val
	}

	periods := numYears * 12

	// value * (1 + rate)^periods
	compoundGrowthRate := math.Pow((1 + monthlyGrowthRate), float64(periods))

	result := float64(val) * compoundGrowthRate
	return int64(result)
}
