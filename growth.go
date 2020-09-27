package bucketgrowth

import (
  "fmt"
  "log"
  "time"
)

// Custom error for when a metric isn't found in an array of metrics
var ErrDailyMetricNotFound error = fmt.Errorf("Metric not found")

// Returns the decimal value that equates to unit/month
// Assumes array in already sorted by time (oldest->newest)
func monthlyGrowthPct(metrics []DailyMetric) float64 {
  var monthlyGrowthRates [12]float64

  for i, ending := 0, metrics[len(metrics)-1]; i < 12; i++ {
    var starting DailyMetric
    startingDate := ending.Date.AddDate(0, -1, 0)
    starting, err := findMetricByDate(metrics, startingDate)
    if err != nil {
      starting = metrics[0]
    }
    log.Printf("[monthlyGrowth] Starting Date: %s\n", starting.Date)
    log.Printf("[monthlyGrowth] Ending Date: %s\n", ending.Date)
    growthRate := midpointGrowthRate(starting, ending)

    monthlyGrowthRates[12-1-i] = growthRate
    ending = starting
  }

  log.Printf("Growth rates: %v\n", monthlyGrowthRates)

  avgMonthlyGrowthRate := average(monthlyGrowthRates[:])

  log.Printf("[monthlyGrowth] Growth rate: %f\n", avgMonthlyGrowthRate)

  return avgMonthlyGrowthRate * 100
}

func average(metrics []float64) float64 {
  var sum float64
  for _, val := range metrics {
    sum += val
  }

  return sum / float64(len(metrics))
}

func findMetricByDate(metrics []DailyMetric, date time.Time) (DailyMetric, error) {
  var metric DailyMetric
	for i, val := range metrics {
		if val.Date == date {
      metric = metrics[i]
			break
		}
	}

  if metric == (DailyMetric{}) {
    return DailyMetric{}, ErrDailyMetricNotFound
  }

  return metric, nil
}

// Returns the decimal value that equates to unit/year
// Assumes array in already sorted by time (oldest->newest)
func yearlyGrowthPct(metrics []DailyMetric) float64 {
  ending := metrics[len(metrics)-1]
  starting := metrics[0]

  log.Printf("[yearlyGrowth] Starting value: %d\n", starting.Total)
  log.Printf("[yearlyGrowth] Ending value: %d\n", ending.Total)

  growthRate := midpointGrowthRate(starting, ending)
  log.Printf("[yearlyGrowth] Growth rate: %f\n", growthRate)

	return float64(growthRate * 100)
}

func midpointGrowthRate(starting, ending DailyMetric) float64 {
	return float64(ending.Total - starting.Total)/float64((starting.Total + ending.Total)/2)
}
