package application

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"forum/internal/repository"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	ctx   context.Context
	db    *sql.DB
	repo  *repository.Repository
	cache map[string]repository.User
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


func (a *App) SignupPage(w http.ResponseWriter, message string) {
	sp := filepath.Join("public", "html", "signup.html")

	tmpl, err := template.ParseFiles(sp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type answer struct {
		Message string
	}
	data := answer{message}

	err = tmpl.ExecuteTemplate(w, "signup", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a *App) Signup(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.FormValue("name"))
	surname := strings.TrimSpace(r.FormValue("surname"))
	login := strings.TrimSpace(r.FormValue("login"))
	password := strings.TrimSpace(r.FormValue("password"))
	password2 := strings.TrimSpace(r.FormValue("password2"))

	if name == "" || surname == "" || login == "" || password == "" {
		a.SignupPage(w, "Все поля должны быть заполнены!")
		return
	}

	if password != password2 {
		a.SignupPage(w, "Пароли не совпадают! Попробуйте еще")
		return
	}

	hash := md5.Sum([]byte(password))
	hashedPass := hex.EncodeToString(hash[:])

	err := a.repo.AddNewUser(a.ctx, name, surname, login, hashedPass)
	if err != nil {
		a.SignupPage(w, fmt.Sprintf("Ошибка создания пользователя: %v", err))
		return
	}

	a.LoginPage(w, fmt.Sprintf("%s, вы успешно зарегистрированы! Теперь вам доступен вход через страницу авторизации", name))
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

func NewApp(ctx context.Context, db *sql.DB) *App {
	return &App{ctx, db, repository.NewRepository(db), make(map[string]repository.User)}
}
