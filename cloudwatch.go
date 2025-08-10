package bucketgrowth

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

var now = time.Now()

type Metrics struct {
	TotalSizeBytes   int64 `json:"total_size_bytes"`
	TotalObjectCount int64 `json:"total_object_count"`

	SizeGrowthMonthly float64 `json:"size_growth_monthly"`
	SizeGrowthYearly  float64 `json:"size_growth_yearly"`

	ObjectGrowthMonthly float64 `json:"object_growth_monthly"`
	ObjectGrowthYearly  float64 `json:"object_growth_yearly"`

	Projections ProjectionMetrics `json:"projections"`
}

type ProjectionMetrics struct {
	Size1Year int64 `json:"size_bytes_1_year"`
	Size5Year int64 `json:"size_bytes_5_year"`

	Object1Year int64 `json:"object_count_1_year"`
	Object5Year int64 `json:"object_count_5_year"`
}

type DailyMetric struct {
	Date  time.Time
	Total int64
}

type Request struct {
	BucketName string

	CloudWatchClient *cloudwatch.Client
}

func (self Request) Measure() (Metrics, error) {
	now := time.Now()
	oneYearAgo := now.AddDate(-1, 0, 0) // gives enough data for stats

	sizeMetrics, err := self.sizeData(oneYearAgo)
	if err != nil {
		return Metrics{}, err
	}
	objectMetrics, err := self.objectData(oneYearAgo)
	if err != nil {
		return Metrics{}, err
	}

	if len(sizeMetrics) == 0 {
		return Metrics{}, fmt.Errorf("Bucket is empty")
	}

	sort.Slice(sizeMetrics, func(i, j int) bool { return sizeMetrics[i].Date.Before(sizeMetrics[j].Date) })
	sort.Slice(objectMetrics, func(i, j int) bool { return objectMetrics[i].Date.Before(objectMetrics[j].Date) })

	var metrics Metrics

	// populate the metrics struct
	metrics.TotalSizeBytes = sizeMetrics[len(sizeMetrics)-1].Total
	metrics.TotalObjectCount = objectMetrics[len(objectMetrics)-1].Total

	// calculate growth
	metrics.SizeGrowthMonthly = monthlyGrowthPct(sizeMetrics)
	metrics.ObjectGrowthMonthly = monthlyGrowthPct(objectMetrics)
	metrics.SizeGrowthYearly = yearlyGrowthPct(sizeMetrics)
	metrics.ObjectGrowthYearly = yearlyGrowthPct(objectMetrics)

	// calculation projection
	// using monthly growth rate because it is more accurate than yearly
	proj := ProjectionMetrics{}

	proj.Size1Year = projection(metrics.TotalSizeBytes, 1, metrics.SizeGrowthMonthly/100.0)
	proj.Object1Year = projection(metrics.TotalObjectCount, 1, metrics.ObjectGrowthMonthly/100.0)

	proj.Size5Year = projection(metrics.TotalSizeBytes, 5, metrics.SizeGrowthMonthly/100.0)
	proj.Object5Year = projection(metrics.TotalObjectCount, 5, metrics.ObjectGrowthMonthly/100.0)

	metrics.Projections = proj

	return metrics, nil
}

func endOfResults(token string) bool {
	return token == ""
}

func (self Request) sizeData(start time.Time) ([]DailyMetric, error) {
	input := cloudwatch.GetMetricStatisticsInput{
		StartTime: &start,
		EndTime:   &now,
		Dimensions: []types.Dimension{
			{
				Name:  strPtr("BucketName"),
				Value: strPtr(self.BucketName),
			},
			{
				Name:  strPtr("StorageType"),
				Value: strPtr("StandardStorage"),
			},
		},
		MetricName: strPtr("BucketSizeBytes"),
		Namespace:  strPtr("AWS/S3"),
		Period:     int32Ptr(86400),
		Statistics: []types.Statistic{
			types.StatisticMaximum,
		},
		Unit: types.StandardUnitBytes,
	}

	log.Printf("Retrieving total size from: %s\n", self.BucketName)
	resp, err := self.CloudWatchClient.GetMetricStatistics(context.Background(), &input)
	if err != nil {
		return []DailyMetric{}, err
	}

	results := marshalDailyMetric(resp.Datapoints)

	return results, nil
}

func (self Request) objectData(start time.Time) ([]DailyMetric, error) {
	input := cloudwatch.GetMetricStatisticsInput{
		StartTime: &start,
		EndTime:   &now,
		Dimensions: []types.Dimension{
			{
				Name:  strPtr("BucketName"),
				Value: strPtr(self.BucketName),
			},
			{
				Name:  strPtr("StorageType"),
				Value: strPtr("AllStorageTypes"),
			},
		},
		MetricName: strPtr("NumberOfObjects"),
		Namespace:  strPtr("AWS/S3"),
		Period:     int32Ptr(86400),
		Statistics: []types.Statistic{
			types.StatisticMaximum,
		},
		Unit: types.StandardUnitCount,
	}

	log.Printf("Retrieving total object count from: %s\n", self.BucketName)
	resp, err := self.CloudWatchClient.GetMetricStatistics(context.Background(), &input)
	if err != nil {
		return []DailyMetric{}, err
	}

	results := marshalDailyMetric(resp.Datapoints)

	return results, nil
}

func marshalDailyMetric(results []types.Datapoint) []DailyMetric {
	metrics := make([]DailyMetric, len(results))

	for i, val := range results {
		metrics[i] = DailyMetric{
			Date: *val.Timestamp,
			// loss of precision during cast doesn't matter since we are working
			// with whole numbers here
			Total: int64(*val.Maximum),
		}
	}

	return metrics
}

func strPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}
