package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	log.SetOutput(ioutil.Discard)
	app := &cli.App{
		Name:   "bucketgrowth",
		Usage:  "Display size and growth statistics of an S3 bucket.",
		Action: perform,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "profile",
				Usage:       "AWS profile to use",
				Destination: &flagProfile,
				EnvVars:     []string{"AWS_PROFILE"},
			},
			&cli.StringFlag{
				Name:        "region",
				Usage:       "AWS region to use",
				Destination: &flagRegion,
				EnvVars:     []string{"AWS_DEFAULT_REGION"},
			},
			&cli.BoolFlag{
				Name:        "verbose",
				Usage:       "Enable verbose logging",
				Aliases:     []string{"v"},
				Destination: &flagVerbose,
			},
			&cli.StringFlag{
				Name:        "output",
				Usage:       fmt.Sprintf("Changes the output `TYPE` to %s or %s", outputText, outputJson),
				Value:       outputText,
				Destination: &flagOutputType,
			},
			&cli.BoolFlag{
				Name:        "skip-banner",
				Usage:       "Suppresses output of the banner",
				Destination: &flagSkipBanner,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
