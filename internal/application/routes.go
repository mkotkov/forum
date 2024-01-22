package application

import (
	"log"
	"net/http"
	"strings"
)

func (a *App) Routes(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/":
		a.authorized(a.MainPage).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/post/"):
		// Extract the slug from the URL path
		slug := strings.TrimPrefix(r.URL.Path, "/post/")
		log.Println("Slug:", slug)
		// Handle the dynamic post page
		a.authorized(func(w http.ResponseWriter, r *http.Request) {
			a.PostPage(w, r, slug)
		}).ServeHTTP(w, r)
	case r.URL.Path == "/create/":
		a.authorized(a.Create).ServeHTTP(w, r)
	case r.URL.Path == "/save_post/":
		if r.Method == http.MethodPost {
			// Handle the save post logic here
			a.SavePost(w, r)
		} else {
			// Respond with an error or handle as needed for other HTTP methods
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	case strings.HasPrefix(r.URL.Path, "/save_comment/"):
		// Handle the save comment logic here
		a.SaveComment(w, r)
	case r.URL.Path == "/login":
		if r.Method == http.MethodGet {
			a.LoginPage(w, "")
		} else if r.Method == http.MethodPost {
			a.Login(w, r)
		}
	case r.URL.Path == "/logout":
		a.Logout(w, r)
	case r.URL.Path == "/signup":
		if r.Method == http.MethodGet {
			a.SignupPage(w, "")
		} else if r.Method == http.MethodPost {
			a.Signup(w, r)
		}
		// в switch в функции Routes
	case strings.HasPrefix(r.URL.Path, "/like_post/"):
		// Извлекаем slug из URL
		slug := strings.TrimPrefix(r.URL.Path, "/like_post/")
		a.authorized(func(w http.ResponseWriter, r *http.Request) {
			a.ReactPost(w, r, slug, "like")
		}).ServeHTTP(w, r)

	case strings.HasPrefix(r.URL.Path, "/dislike_post/"):
		// Извлекаем slug из URL
		slug := strings.TrimPrefix(r.URL.Path, "/dislike_post/")
		a.authorized(func(w http.ResponseWriter, r *http.Request) {
			a.ReactPost(w, r, slug, "dislike")
		}).ServeHTTP(w, r)

	default:
		// Перенаправление запросов для статических файлов в http.FileServer
		http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))).ServeHTTP(w, r)
	}
}
