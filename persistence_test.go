package main

import (
	"fmt"
	"testing"
	"time"
)

var (
	t1 time.Time
	t2 time.Time
	t3 time.Time
)

func init() {
	t1 = time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC)
	t2 = time.Date(2017, 1, 3, 4, 5, 6, 0, time.UTC)
	t3 = time.Date(2017, 6, 5, 4, 3, 2, 0, time.UTC)
}

func setup(t *testing.T) {
	err := createDBFromEnv()
	if err != nil {
		t.Error(err)
	}

	_, err = db.Exec("CREATE TABLE occurrences(channelid BIGINT UNSIGNED, ts TIMESTAMP)")
	if err != nil {
		t.Error("Could not create test table", err)
	}

	_, err = db.Exec("INSERT INTO occurrences VALUES (1234, ?), (1234, ?), (5678, ?)", t1, t2, t3)
	if err != nil {
		t.Error("Could not populate test table", err)
	}

	initDB()
}

func teardown(t *testing.T) {
	_, err := db.Exec("DROP TABLE occurrences")
	if err != nil {
		t.Error("Could not drop test table", err)
	}
	db.Close()
}

func TestGetMostRecentOccurrence(t *testing.T) {
	setup(t)
	defer teardown(t)

	result, err := getMostRecentOccurrence("1234")
	if err != nil {
		t.Error(err)
	}
	if result != t2 {
		t.Errorf("Retrieved timestamp did not equal expected timestamp (expected %s, got %s)", t2, result)
	}

	result, err = getMostRecentOccurrence("5678")
	if err != nil {
		t.Error(err)
	}
	if result != t3 {
		t.Errorf("Retrieved timestamp did not equal expected timestamp (expected %s, got %s)", t3, result)
	}
}

func TestInsertOccurrence(t *testing.T) {
	setup(t)
	defer teardown(t)

	channel := uint64(1234)
	ts := time.Now().Round(0)

	err := insertOccurence(fmt.Sprintf("%d", channel), ts)
	if err != nil {
		t.Error(err)
	}

	var retrievedTs time.Time
	err = db.QueryRow("SELECT MAX(ts) FROM occurrences WHERE channelid = ?", channel).Scan(&retrievedTs)
	if err != nil {
		t.Error(err)
	}

	if retrievedTs.Equal(ts) {
		t.Errorf("Retrieved timestamp did not equal inserted timestamp (expected %s, got %s)", ts, retrievedTs)
	}
}
