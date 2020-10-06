package main

import (
	"errors"

	"github.com/urfave/cli/v2"
)

const outputText string = "text"

const outputJson string = "json"

var flagProfile string

var flagRegion string

var flagVerbose bool

var flagOutputType string

var flagSkipBanner bool

var ErrUnsupportedOutputType error = errors.New("Unsupported output type")

var validOutputTypes []string = []string{outputText, outputJson}

func guardOutputType() error {
	for _, val := range validOutputTypes {
		if val == flagOutputType {
			return nil // flag has been found, return no error
		}
	}

	return ErrUnsupportedOutputType
}

func guardBucketArg(c *cli.Context, bucket string) error {
	if bucket != "" {
		return nil
	}

	if err := cli.ShowAppHelp(c); err != nil {
		return err
	}

	return errors.New("") // blank error to show there's an issue
}
