package bucketgrowth

import (
	"io/ioutil"
	"log"
	"math"
	"testing"
	"time"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestAverage(t *testing.T) {
	cases := []struct {
		name     string
		metrics  []float64
		expected float64
	}{
		{
			"valid result",
			[]float64{float64(1), float64(2)},
			float64(1.5),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := average(c.metrics)

			if c.expected != actual {
				t.Fatalf("Expected %v but got %v", c.expected, actual)
			}
		})
	}
}

func TestMonthlyGrowthPct(t *testing.T) {
	now := time.Now()
	metrics := []DailyMetric{
		{
			Date:  now.AddDate(0, -2, 0),
			Total: 100,
		},
		{
			Date:  now.AddDate(0, -1, 0),
			Total: 200,
		},
		{
			Date:  now,
			Total: 1000,
		},
	}

	expected := float64(16.67)
	actual := monthlyGrowthPct(metrics)
	// due to floating-point precision, we need to round in order to pass
	// this test.
	actual = math.Round(actual*100) / 100

	if expected != actual {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestYearlyGrowthPct(t *testing.T) {
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

	expected := float64(196.04)
	actual := yearlyGrowthPct(metrics)
	// due to floating-point precision, we need to round in order to pass
	// this test.
	actual = math.Round(actual*100) / 100

	if expected != actual {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestMidpointGrowthRate(t *testing.T) {
	cases := []struct {
		name     string
		start    DailyMetric
		end      DailyMetric
		expected float64
	}{
		{
			"valid result",
			DailyMetric{
				Total: 100,
			},
			DailyMetric{
				Total: 1000,
			},
			float64(1.63636),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := midpointGrowthRate(c.start, c.end)
			actual = math.Round(actual*100000) / 100000

			if c.expected != actual {
				t.Fatalf("Expected %v but got %v", c.expected, actual)
			}
		})
	}
}

func TestFindMetricByDate(t *testing.T) {
	now := time.Now()

	cases := []struct {
		name        string
		metrics     []DailyMetric
		date        time.Time
		expected    DailyMetric
		expectedErr error
	}{
		{
			"valid result",
			[]DailyMetric{
				DailyMetric{
					Date: now.AddDate(0, -8, 0),
				},
				DailyMetric{
					Date: now.AddDate(0, -6, 0),
				},
				DailyMetric{
					Date: now,
				},
			},
			now.AddDate(0, -6, 0),
			DailyMetric{
				Date: now.AddDate(0, -6, 0),
			},
			nil,
		},
		{
			"date not found",
			[]DailyMetric{
				DailyMetric{
					Date: now.AddDate(0, -8, 0),
				},
				DailyMetric{
					Date: now.AddDate(0, -6, 0),
				},
				DailyMetric{
					Date: now,
				},
			},
			now.AddDate(0, -7, 0),
			DailyMetric{},
			ErrDailyMetricNotFound,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := findMetricByDate(c.metrics, c.date)

			if c.expectedErr != err {
				t.Fatalf("Expected %v but got %v", c.expectedErr, err)
			}

			if c.expected != actual {
				t.Fatalf("Expected %v but got %v", c.expected, actual)
			}
		})
	}
}
