package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var DatabaseAddress string
var DatabaseUsername string
var DatabasePassword string
var db *sql.DB
var insertStmt *sql.Stmt
var mostRecentStmt *sql.Stmt
var longestStmnt *sql.Stmt

var ErrNoOccurrence = errors.New("No occurences found")

func createDBFromEnv() error {
	DatabaseAddress = MustEnv("DATABASE_ADDRESS")
	DatabaseUsername = MustEnv("DATABASE_USERNAME")
	DatabasePassword = MustEnv("DATABASE_PASSWORD")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/c9bot?parseTime=true", DatabaseUsername, DatabasePassword, DatabaseAddress)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	err = db.Ping()
	return err
}

func initDB() error {
	var err error
	insertStmt, err = db.Prepare("INSERT INTO occurrences VALUES (?, ?)")
	if err != nil {
		return err
	}

	mostRecentStmt, err = db.Prepare("SELECT MAX(ts) FROM occurrences WHERE channelid = ?")
	if err != nil {
		return err
	}

	longestStmnt, err = db.Prepare(`
		SELECT MAX(TIMESTAMPDIFF(DAY, prevts, ts))
        FROM (
            SELECT 
                ts,
                @prev AS prevts,
                @prev := ts
            FROM occurrences, (select @prev:=NULL) vars
            WHERE channelid = ?
            ORDER BY ts
        ) lagged`)
	if err != nil {
		return err
	}

	return nil
}

func insertOccurence(channel string, ts time.Time) error {
	cid, err := parseChannelId(channel)
	if err != nil {
		return err
	}
	_, err = insertStmt.Exec(cid, ts)
	return err
}

func getMostRecentOccurrence(channel string) (time.Time, error) {
	var mostRecent time.Time

	cid, err := parseChannelId(channel)
	if err != nil {
		return mostRecent, err
	}

	err = mostRecentStmt.QueryRow(cid).Scan(&mostRecent)
	switch {
	case err == sql.ErrNoRows:
		return mostRecent, ErrNoOccurrence
	case err != nil:
		return mostRecent, err
	}
	return mostRecent, nil
}

func getLongestGap(channel string) (int, error) {
	var longest sql.NullInt64

	cid, err := parseChannelId(channel)
	if err != nil {
		return 0, err
	}

	err = longestStmnt.QueryRow(cid).Scan(&longest)
	switch {
	case err == sql.ErrNoRows:
		return 0, ErrNoOccurrence
	case err != nil:
		return 0, err
	}
	if longest.Valid {
		return int(longest.Int64), nil
	}
	return 0, ErrNoOccurrence
}
