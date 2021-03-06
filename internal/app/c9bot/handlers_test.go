package c9bot

import (
	"testing"
	"time"
)

func TestDurationString(t *testing.T) {
	testCase := func(d time.Duration, exp string) func(*testing.T) {
		return func(t *testing.T) {
			result := durationString(d)
			if result != exp {
				t.Errorf("DurationString was incorrect. Expected %s, got %s with input %s", exp, result, d)
			}
		}
	}

	t.Run("2 days", testCase(48*time.Hour, "2 days"))
	t.Run("1.5 days", testCase(36*time.Hour, "1 day"))
	t.Run("1 day", testCase(24*time.Hour, "1 day"))
	t.Run("> 1 day", testCase(2*time.Hour, "less than a day"))
}
