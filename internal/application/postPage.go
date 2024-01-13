package application

import (
	"forum/internal/repository"
	"html/template"
	"net/http"
)

func (a *App) PostPage(w http.ResponseWriter, r *http.Request, slug string) {
	// Retrieve the post based on the slug
	post, err := a.repo.GetPostBySlug(a.ctx, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(
		"public/html/header.html",
		"public/html/footer.html",
		"public/html/post.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    data := struct {
        Post repository.Posts
    }{
        Post: post,
    }

    err = tmpl.ExecuteTemplate(w, "post", data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
