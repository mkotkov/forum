package application

import (
	"net/http"
	"text/template"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"public/html/error.html",
		"public/html/footer.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := struct {
		Status int
	}{
		Status: http.StatusNotFound,
	}

	w.WriteHeader(http.StatusNotFound)
	err = tmpl.ExecuteTemplate(w, "error", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
