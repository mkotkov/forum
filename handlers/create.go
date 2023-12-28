package handlers

import (
	"net/http"
	"text/template"
)

func Create(w http.ResponseWriter, r *http.Request) {
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