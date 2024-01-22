package application

import (
	"net/http"

	"forum/internal/repository"
)

func (a *App) authorized(handler func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := readCookie("token", r)
		if err != nil || a.cache[token] == (repository.User{}) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		handler(w, r)
	})
}
