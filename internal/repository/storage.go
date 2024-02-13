package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func InitDBConn(ctx context.Context) (db *sql.DB, err error) {
	// Use SQLite file path instead of URL
	db, err = sql.Open("sqlite3", "./db/database.db")
	if err != nil {
		err = fmt.Errorf("failed to open SQLite DB connection: %w", err)
		return
	}

	// Setting the maximum number of open and idle connections, as well as timeouts
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(24 * time.Hour)
	db.SetConnMaxIdleTime(30 * time.Minute)

	// Setting timeouts for connections
	conn, err := db.Conn(ctx)
	if err != nil {
		err = fmt.Errorf("failed to establish a new connection: %w", err)
		return
	}
	defer conn.Close()

	err = conn.PingContext(ctx)
	if err != nil {
		err = fmt.Errorf("failed to ping the database: %w", err)
		return
	}

	return
}
