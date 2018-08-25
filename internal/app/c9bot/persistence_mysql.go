package c9bot

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // load the mysql driver
)

type mySQLPersistenceLayer struct {
	db                                      *sql.DB
	insertStmt, mostRecentStmt, longestStmt *sql.Stmt
}

func createmySQLPersistenceLayer(address string, username string, password string) (mySQLPersistenceLayer, error) {

	p := mySQLPersistenceLayer{}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/c9bot?parseTime=true", username, password, address)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return p, err
	}

	err = db.Ping()
	if err != nil {
		return p, err
	}

	insertStmt, err := db.Prepare("INSERT INTO occurrences VALUES (?, ?)")
	if err != nil {
		return p, err
	}

	mostRecentStmt, err := db.Prepare("SELECT MAX(ts) FROM occurrences WHERE channelid = ?")
	if err != nil {
		return p, err
	}

	longestStmt, err := db.Prepare(`
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
		return p, err
	}
	p = mySQLPersistenceLayer{db, insertStmt, mostRecentStmt, longestStmt}
	return p, nil
}

func (p mySQLPersistenceLayer) Close() {
	p.db.Close()
}

func (p mySQLPersistenceLayer) insertOccurence(channelID string, ts time.Time) error {
	cid, err := parseChannelID(channelID)
	if err != nil {
		return err
	}
	_, err = p.insertStmt.Exec(cid, ts)
	return err
}

func (p mySQLPersistenceLayer) getMostRecentOccurrence(channelID string) (time.Time, error) {
	var mostRecent time.Time

	cid, err := parseChannelID(channelID)
	if err != nil {
		return mostRecent, err
	}

	err = p.mostRecentStmt.QueryRow(cid).Scan(&mostRecent)
	switch {
	case err == sql.ErrNoRows:
		return mostRecent, errNoOccurrence
	case err != nil:
		return mostRecent, err
	}
	return mostRecent, nil
}

func (p mySQLPersistenceLayer) getLongestGap(channelID string) (int, error) {
	var longest sql.NullInt64

	cid, err := parseChannelID(channelID)
	if err != nil {
		return 0, err
	}

	err = p.longestStmt.QueryRow(cid).Scan(&longest)
	switch {
	case err == sql.ErrNoRows:
		return 0, errNoOccurrence
	case err != nil:
		return 0, err
	}
	if longest.Valid {
		return int(longest.Int64), nil
	}
	return 0, errNoOccurrence
}
