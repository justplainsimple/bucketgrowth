package bucketgrowth

import (
	"testing"
)

func TestProjection(t *testing.T) {
	cases := []struct {
		name              string
		val               int64
		numYears          int
		monthlyGrowthRate float64
		expected          int64
	}{
		{
			"happy path",
			1000,
			1,
			0.1,
			3138,
		},
		{
			"invalid val",
			0,
			10,
			0.1,
			0,
		},
		{
			"invalid year",
			1000,
			0,
			0.1,
			0,
		},
		{
			"invalid monthlyGrowthRate",
			1000,
			10,
			0.0,
			1000,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := projection(c.val, c.numYears, c.monthlyGrowthRate)

			if c.expected != actual {
				t.Fatalf("Expected %v but got %v", c.expected, actual)
			}
		})
	}
}
