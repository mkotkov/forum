package application

import (
	"net/http"
	"strconv"
	"text/template"
	"time"

	"forum/internal/repository"
)

func (a *App) UnregPage(w http.ResponseWriter, r *http.Request, message string, isAuthorized bool) {
	tmpl, err := template.ParseFiles(
		"public/html/unreg-header.html",
		"public/html/footer.html",
		"public/html/forum-card.html",
		"public/html/login-form.html",
		"public/html/signup.html",
		"public/html/rec-post.html",
		"public/html/index.html",
		"public/html/signup.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filter := r.FormValue("filters")
	topicFilter := r.FormValue("topic")
	selectedTopic, err := strconv.Atoi(topicFilter)
	if err != nil {
		selectedTopic = -1
	}

	var posts []repository.Posts

	if selectedTopic == 0 {
		posts, err = QueryData(a.repo.GetDB(), filter, -1)
	} else {
		posts, err = QueryData(a.repo.GetDB(), filter, selectedTopic)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mostRecentPost, err := QueryMostRecentPost(a.repo.GetDB())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	topics, err := a.repo.GetAllTopics(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Posts          []repository.Posts
		MostRecentPost *repository.Posts
		Topics         []repository.Topic
		TimeZone       *time.Location
		Filter         string
		SelectedTopic  int
		Message        string
		IsAuthorized   bool
	}{
		Posts:          posts,
		MostRecentPost: mostRecentPost,
		Topics:         topics,
		TimeZone:       time.UTC,
		Filter:         filter,
		SelectedTopic:  selectedTopic,
		Message:        message,
		IsAuthorized:   isAuthorized,
	}

	err = tmpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
