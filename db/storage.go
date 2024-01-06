package db

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    _ "github.com/mattn/go-sqlite3"
)

func InitDBConn(ctx context.Context) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", "./db/data.db")
    if err != nil {
        return nil, fmt.Errorf("failed to open SQLite connection: %w", err)
    }

    if err := db.Ping(); err != nil {
        db.Close()
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    db.SetMaxOpenConns(5)
    db.SetMaxIdleConns(1)
    db.SetConnMaxLifetime(24 * time.Hour)
    db.SetConnMaxIdleTime(30 * time.Minute)

    fmt.Println("Connected to SQLite database")

    return db, nil
}
