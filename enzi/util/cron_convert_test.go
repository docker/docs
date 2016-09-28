package util

import (
	"fmt"
	"testing"
)

type minuteRange struct {
	// From min up to and not including max.
	min int
	max int
	// I is the offset from the min.
	result func(i int) string
}

func TestMinutesToCron(t *testing.T) {
	// These ranges cover all results from 0 minutes to 31 days.
	testRanges := []*minuteRange{
		{
			min:    0,
			max:    60 * 2,
			result: func(i int) string { return "@hourly" },
		},
		{
			min: 60 * 2,
			max: 60 * 24,
			result: func(i int) string {
				hours := (120 + i) / 60
				return fmt.Sprintf("0 0 */%d * *", hours)
			},
		},
		{
			min:    60 * 24,
			max:    60 * 24 * 2,
			result: func(i int) string { return "@daily" },
		},
		{
			min: 60 * 24 * 2,
			max: 60 * 24 * 28,
			result: func(i int) string {
				days := (60*24*2 + i) / (60 * 24)
				return fmt.Sprintf("0 0 0 */%d *", days)
			},
		},
		{
			min:    60 * 24 * 28,
			max:    60 * 24 * 31,
			result: func(i int) string { return "@monthly" },
		},
	}

	for _, testRange := range testRanges {
		for minutes := testRange.min; minutes < testRange.max; minutes++ {
			cronExpression := MinutesToCron(minutes)
			expectedResult := testRange.result(minutes - testRange.min)

			if cronExpression != expectedResult {
				t.Fatalf("%d minutes: expected %s, got %s", minutes, expectedResult, cronExpression)
			}
		}
	}
}
