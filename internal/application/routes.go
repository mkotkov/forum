package application

import (
	"net/http"
)

func (a *App) Routes(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
    case "/":
		a.authorized(a.MainPage).ServeHTTP(w, r)
	case "/create/":
		a.authorized(a.Create).ServeHTTP(w, r)
	case "/save_post/":
		if r.Method == http.MethodPost {
			// Handle the save post logic here
			a.SavePost(w, r)
		} else {
			// Respond with an error or handle as needed for other HTTP methods
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
	case "/login":
		if r.Method == http.MethodGet {
			a.LoginPage(w, "")
		} else if r.Method == http.MethodPost {
			a.Login(w, r)
		}
	case "/logout":
		a.Logout(w, r)
	case "/signup":
		if r.Method == http.MethodGet {
			a.SignupPage(w, "")
		} else if r.Method == http.MethodPost {
			a.Signup(w, r)
		}
	default:
		// Перенаправление запросов для статических файлов в http.FileServer
		http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))).ServeHTTP(w, r)
	}
}