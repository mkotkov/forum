package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func InitDBConn(ctx context.Context) (db *sql.DB, err error) {
	// Используйте путь к файлу SQLite вместо URL
	db, err = sql.Open("sqlite3", "./db/database.db")
	if err != nil {
		err = fmt.Errorf("failed to open SQLite DB connection: %w", err)
		return
	}

	// Установка максимального количества открытых и простаивающих соединений, а также таймаутов
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(24 * time.Hour)
	db.SetConnMaxIdleTime(30 * time.Minute)

	// Устанавливаем таймауты для подключения
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
