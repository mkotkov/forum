package application

import (
	"crypto/md5"
	"encoding/hex"
	"html/template"
	"net/http"
	"time"
)

func (a *App) LoginPage(w http.ResponseWriter, message string) {
	tmpl, err := template.ParseFiles(
		"public/html/header.html", 
		"public/html/footer.html", 
		"public/html/forum-card.html", 
		"public/html/start-page.html", 
		"public/html/login-form.html")
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

func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login-inp")
	password := r.FormValue("password-inp")

	if login == "" || password == "" {
		a.LoginPage(w, "Необходимо указать логин и пароль!")
		return
	}

	hash := md5.Sum([]byte(password))
	hashedPass := hex.EncodeToString(hash[:])

	user, err := a.repo.Login(a.ctx, login, hashedPass)
	if err != nil {
		a.LoginPage(w, "Вы ввели неверный логин или пароль!")
		return
	}

	time64 := time.Now().Unix()
	timeInt := string(rune(time64))
	token := login + password + timeInt

	hashToken := md5.Sum([]byte(token))
	hashedToken := hex.EncodeToString(hashToken[:])

	a.cache[hashedToken] = user

	livingTime := 60 * time.Minute
	expiration := time.Now().Add(livingTime)

	cookie := http.Cookie{Name: "token", Value: hashedToken, Expires: expiration}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	for _, v := range r.Cookies() {
		c := http.Cookie{
			Name:   v.Name,
			MaxAge: -1}
		http.SetCookie(w, &c)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
