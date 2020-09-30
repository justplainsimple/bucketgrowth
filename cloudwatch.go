package bucketgrowth

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
)

var now = time.Now()

type Metrics struct {
	TotalSizeBytes   int64 `json:"total_size_bytes"`
	TotalObjectCount int64 `json:"total_object_count"`

	SizeGrowthMonthly float64 `json:"size_growth_monthly"`
	SizeGrowthYearly  float64 `json:"size_growth_yearly"`

	ObjectGrowthMonthly float64 `json:"object_growth_monthly"`
	ObjectGrowthYearly  float64 `json:"object_growth_yearly"`
}

type DailyMetric struct {
	Date  time.Time
	Total int64
}

type Request struct {
	BucketName string

	CloudWatchClient cloudwatchiface.CloudWatchAPI
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

	return metrics, nil
}

func endOfResults(token string) bool {
	return token == ""
}

func (self Request) sizeData(start time.Time) ([]DailyMetric, error) {
	input := cloudwatch.GetMetricStatisticsInput{
		StartTime: &start,
		EndTime:   &now,
		Dimensions: []*cloudwatch.Dimension{
			&cloudwatch.Dimension{
				Name:  aws.String("BucketName"),
				Value: aws.String(self.BucketName),
			},
			&cloudwatch.Dimension{
				Name:  aws.String("StorageType"),
				Value: aws.String("StandardStorage"),
			},
		},
		MetricName: aws.String("BucketSizeBytes"),
		Namespace:  aws.String("AWS/S3"),
		Period:     aws.Int64(86400),
		Statistics: []*string{
			aws.String("Maximum"),
		},
		Unit: aws.String("Bytes"),
	}

	log.Printf("Retrieving total size from: %s\n", self.BucketName)
	resp, err := self.CloudWatchClient.GetMetricStatistics(&input)
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
		Dimensions: []*cloudwatch.Dimension{
			&cloudwatch.Dimension{
				Name:  aws.String("BucketName"),
				Value: aws.String(self.BucketName),
			},
			&cloudwatch.Dimension{
				Name:  aws.String("StorageType"),
				Value: aws.String("AllStorageTypes"),
			},
		},
		MetricName: aws.String("NumberOfObjects"),
		Namespace:  aws.String("AWS/S3"),
		Period:     aws.Int64(86400),
		Statistics: []*string{
			aws.String("Maximum"),
		},
		Unit: aws.String("Count"),
	}

	log.Printf("Retrieving total object count from: %s\n", self.BucketName)
	resp, err := self.CloudWatchClient.GetMetricStatistics(&input)
	if err != nil {
		return []DailyMetric{}, err
	}

	results := marshalDailyMetric(resp.Datapoints)

	return results, nil
}

func marshalDailyMetric(results []*cloudwatch.Datapoint) []DailyMetric {
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
