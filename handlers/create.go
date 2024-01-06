package handlers

import (
	"forum/middleware"
	"net/http"
	"text/template"
)


func CreateHandler(repo *middleware.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Create(w, r, repo)
	}
}

func Create (w http.ResponseWriter, r *http.Request, repo *middleware.Repository) {
	tmpl, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "create", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
}