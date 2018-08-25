package c9bot

import (
	"errors"
	"time"
)

type persistenceLayer interface {
	insertOccurence(string, time.Time) error
	getMostRecentOccurrence(string) (time.Time, error)
	getLongestGap(string) (int, error)
	Close()
}

var errNoOccurrence = errors.New("No occurences found")
