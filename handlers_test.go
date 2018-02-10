package main

import (
	"testing"
	"time"
)

func TestDurationString(t *testing.T) {
	testCase := func(d time.Duration, exp string) func(*testing.T) {
		return func(t *testing.T) {
			result := DurationString(d)
			if result != exp {
				t.Errorf("DurationString was incorrect. Expected %s, got %s with input %s", exp, result, d)
			}
		}
	}

	t.Run("2 days", testCase(48*time.Hour, "2 days"))
	t.Run("1.5 days", testCase(36*time.Hour, "1 day"))
	t.Run("1 day", testCase(24*time.Hour, "1 day"))
	t.Run("2 hours", testCase(2*time.Hour, "2 hours"))
	t.Run("1.5 hours", testCase(90*time.Minute, "1 hour"))
	t.Run("1 hour", testCase(time.Hour, "1 hour"))
	t.Run("> hour", testCase(59*time.Minute, "less than 1 hour"))
}
