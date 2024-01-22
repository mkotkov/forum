package application

import (
	"context"
	"database/sql"

	"forum/internal/repository"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	ctx   context.Context
	db    *sql.DB
	repo  *repository.Repository
	cache map[string]repository.User
}

func NewApp(ctx context.Context, db *sql.DB) *App {
	return &App{ctx, db, repository.NewRepository(db), make(map[string]repository.User)}
}
