package bucketgrowth

import (
	"testing"

	_ "github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
)

type Mock struct {
	cloudwatchiface.CloudWatchAPI
}

func TestMeasure(t *testing.T) {
}

func TestEndofResults(t *testing.T) {
	cases := []struct {
		name     string
		token    string
		expected bool
	}{
		{
			"no more results",
			"",
			true,
		},
		{
			"more results",
			"somerandomvalue",
			false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := endOfResults(c.token)

			if c.expected != actual {
				t.Fatalf("Expected %v but got %v", c.expected, actual)
			}
		})
	}
}
