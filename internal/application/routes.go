package application

import (
	"log"
	"net/http"
	"strings"
)

func (a *App) Routes(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/":
		a.authorized(func(w http.ResponseWriter, r *http.Request, isAuthorized bool) {
			a.MainPage(w, r, "", isAuthorized)
		}).ServeHTTP(w, r)

	case strings.HasPrefix(r.URL.Path, "/post/"):
		slug := strings.TrimPrefix(r.URL.Path, "/post/")
		log.Println("Slug:", slug)
		a.authorized(func(w http.ResponseWriter, r *http.Request, isAuthorized bool) {
			a.PostPage(w, r, slug, isAuthorized, "")
		}).ServeHTTP(w, r)

	case r.URL.Path == "/create/":
		a.authorized(func(w http.ResponseWriter, r *http.Request, isAuthorized bool) {
			a.Create(w, r, "")
		}).ServeHTTP(w, r)

	case r.URL.Path == "/save_post/":
		if r.Method == http.MethodPost {
			a.SavePost(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}

	case strings.HasPrefix(r.URL.Path, "/save_comment/"):
		a.SaveComment(w, r)

	case r.URL.Path == "/login":
		if r.Method == http.MethodGet {
			a.UnregPage(w, r, "", false)
		} else if r.Method == http.MethodPost {
			a.Login(w, r)
		}

	case r.URL.Path == "/logout":
		a.Logout(w, r)

	case r.URL.Path == "/signup":
		if r.Method == http.MethodGet {
			a.UnregPage(w, r, "", false)
		} else if r.Method == http.MethodPost {
			a.Signup(w, r)
		}

	case strings.HasPrefix(r.URL.Path, "/like_post/"):
		slug := strings.TrimPrefix(r.URL.Path, "/like_post/")
		a.authorized(func(w http.ResponseWriter, r *http.Request, isAuthorized bool) {
			a.ReactPost(w, r, slug, "like", isAuthorized)
		}).ServeHTTP(w, r)

	case strings.HasPrefix(r.URL.Path, "/dislike_post/"):
		slug := strings.TrimPrefix(r.URL.Path, "/dislike_post/")
		a.authorized(func(w http.ResponseWriter, r *http.Request, isAuthorized bool) {
			a.ReactPost(w, r, slug, "dislike", isAuthorized)
		}).ServeHTTP(w, r)

	case strings.HasPrefix(r.URL.Path, "/like_comment/"):
		commentID := strings.TrimPrefix(r.URL.Path, "/like_comment/")
		a.authorized(func(w http.ResponseWriter, r *http.Request, isAuthorized bool) {
			a.LikeComment(w, r, commentID, isAuthorized)
		}).ServeHTTP(w, r)

	case strings.HasPrefix(r.URL.Path, "/dislike_comment/"):
		commentID := strings.TrimPrefix(r.URL.Path, "/dislike_comment/")
		a.authorized(func(w http.ResponseWriter, r *http.Request, isAuthorized bool) {
			a.DislikeComment(w, r, commentID, isAuthorized)
		}).ServeHTTP(w, r)
	default:
		ErrorHandler(w, r)
	}
}
