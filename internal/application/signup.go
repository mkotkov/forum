package application

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

func (a *App) SignupPage(w http.ResponseWriter, message string) {
	tmpl, err := template.ParseFiles(
		"public/html/header.html", 
		"public/html/footer.html", 
		"public/html/forum-card.html", 
		"public/html/start-page.html", 
		"public/html/login-form.html",
		"public/html/signup.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type answer struct {
		Message string
	}
	data := answer{message}

	err = tmpl.ExecuteTemplate(w, "start-page", data)
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
