package bucketgrowth

import (
	"testing"
	"time"
)

func TestMonthlyGrowth(t *testing.T) {
	now := time.Now()
	metrics := []DailyMetric{
		{
			Date:  now.AddDate(0, -1, 0),
			Total: 100,
		},
		{
			Date:  now,
			Total: 1000,
		},
	}

	expected := float64(900)
	actual := monthlyGrowth(metrics)

	if expected != actual {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestYearlyGrowth(t *testing.T) {
	now := time.Now()
	metrics := []DailyMetric{
		{
			Date:  now.AddDate(-1, 0, 0),
			Total: 10,
		},
		{
			Date:  now.AddDate(0, -1, 0),
			Total: 100,
		},
		{
			Date:  now,
			Total: 1000,
		},
	}

	expected := float64(9900)
	actual := yearlyGrowth(metrics)

	if expected != actual {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}
