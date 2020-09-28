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
		fmt.Println("Bucket Growth")
		fmt.Println("=============")

		fmt.Printf("Total Size: %s\n", humanize.Bytes(uint64(metrics.TotalSizeBytes)))
		fmt.Printf("Total Objects: %s\n", humanize.Comma(metrics.TotalObjectCount))
		fmt.Println("")
		fmt.Printf("Size Growth: %s%%/mo, %s%%/yr\n", humanize.CommafWithDigits(metrics.SizeGrowthMonthly, 2), humanize.CommafWithDigits(metrics.SizeGrowthYearly, 2))
		fmt.Printf("Object Growth: %s%%/mo, %s%%/yr\n", humanize.CommafWithDigits(metrics.ObjectGrowthMonthly, 2), humanize.CommafWithDigits(metrics.ObjectGrowthYearly, 2))
	}

	return nil
}

