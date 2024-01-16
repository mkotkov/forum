package application

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

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


func (a *App) authorized(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := readCookie("token", r)
		if err != nil || a.cache[token] == (repository.User{}) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func readCookie(name string, r *http.Request) (value string, err error) {
	if name == "" {
		return value, errors.New("you are trying to read empty cookie")
	}
	cookie, err := r.Cookie(name)
	if err != nil {
		return value, err
	}
	value = cookie.Value
	return value, err
}

