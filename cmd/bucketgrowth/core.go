package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/dustin/go-humanize"
	"github.com/urfave/cli/v2"

	"bucketgrowth"
)

func perform(c *cli.Context) error {
	if flagVerbose {
		log.SetOutput(os.Stdout)
	}

	bucket := c.Args().Get(0)

	if err := guardBucketArg(c, bucket); err != nil {
		return err
	}

	if err := guardOutputType(); err != nil {
		return err
	}

	if flagProfile != "" {
		log.Printf("Using AWS profile: %s\n", flagProfile)
	}

	if flagRegion != "" {
		log.Printf("Using AWS region: %s\n", flagRegion)
	}

	opts := session.Options{
		Profile: flagProfile,
		Config: aws.Config{
			Region: aws.String(flagRegion),
		},
	}

	req := bucketgrowth.Request{
		BucketName:       bucket,
		CloudWatchClient: cloudwatch.New(session.Must(session.NewSessionWithOptions(opts))),
	}
	metrics, err := req.Measure()
	if err != nil {
		return err
	}

	if err := displayMetrics(metrics); err != nil {
		return err
	}

	return nil
}

func displayMetrics(metrics bucketgrowth.Metrics) error {
	if flagOutputType == outputJson {
		output, err := json.MarshalIndent(metrics, "", "  ")
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", output)
	}

	if flagOutputType == outputText {
		if !flagSkipBanner {
			fmt.Println("Bucket Growth")
			fmt.Println("=============")
		}

    var bytesAsUint uint64
    if metrics.TotalSizeBytes >= 0 {
      bytesAsUint = uint64(metrics.TotalSizeBytes)
    } else {
      bytesAsUint = 0 // Handle negative value
    }
		fmt.Printf("\nTotal Size: %s\n", humanize.Bytes(bytesAsUint))
		fmt.Printf("Total Objects: %s\n", humanize.Comma(metrics.TotalObjectCount))

		fmt.Printf("\nSize Growth: %s%%/mo, %s%%/yr\n", humanize.CommafWithDigits(metrics.SizeGrowthMonthly, 2), humanize.CommafWithDigits(metrics.SizeGrowthYearly, 2))
		fmt.Printf("Object Growth: %s%%/mo, %s%%/yr\n", humanize.CommafWithDigits(metrics.ObjectGrowthMonthly, 2), humanize.CommafWithDigits(metrics.ObjectGrowthYearly, 2))

    var size1YearAsUint uint64
    if metrics.Projections.Size1Year >= 0 {
      // #nosec G115 - We ensure it is non-negative in the Measure function
      size1YearAsUint = uint64(metrics.Projections.Size1Year)
    } else {
      size1YearAsUint = 0 // Handle negative value
    }
    var size5YearAsUint uint64
    if metrics.Projections.Size5Year >= 0 {
      // #nosec G115 - We ensure it is non-negative in the Measure function
      size5YearAsUint = uint64(metrics.Projections.Size5Year)
    } else {
      size5YearAsUint = 0 // Handle negative value
    }
    fmt.Printf("\nSize Projection: %s (1 yr), %s (5 yr)\n", humanize.Bytes(size1YearAsUint), humanize.Bytes(size5YearAsUint))

    fmt.Printf("Object Count Projection: %s (1 yr), %s (5 yr)\n", humanize.Comma(metrics.Projections.Object1Year), humanize.Comma(metrics.Projections.Object5Year))
	}

	return nil
}
