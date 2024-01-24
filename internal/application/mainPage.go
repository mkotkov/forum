package application

import (
	"database/sql"
	"html/template"
	"net/http"
	"time"

	"forum/internal/repository"
)

func (a *App) MainPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"public/html/header.html",
		"public/html/footer.html",
		"public/html/forum-card.html",
		"public/html/rec-post.html",
		"public/html/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch all posts for the main content
	posts, err := queryData(a.repo.GetDB())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the most recent post for the "Recently created" section
	mostRecentPost, err := queryMostRecentPost(a.repo.GetDB())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Include the most recent post in the data for rendering the template
	data := struct {
		Posts          []repository.Posts
		MostRecentPost *repository.Posts
		TimeZone       *time.Location
	}{
		Posts:          posts,
		MostRecentPost: mostRecentPost,
		TimeZone:       time.UTC,
	}

	// Render the template
	err = tmpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func queryData(db *sql.DB) ([]repository.Posts, error) {
	rows, err := db.Query(repository.SQLSelectAllPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []repository.Posts

	for rows.Next() {
		var post repository.Posts
		err := rows.Scan(&post.Id, &post.Author, &post.PostDate, &post.Title, &post.FullText, &post.Slug, &post.LikeCount, &post.DislikeCount)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func queryMostRecentPost(db *sql.DB) (*repository.Posts, error) {
	row := db.QueryRow(repository.SQLSelectMostRecentPost)

	var post repository.Posts
	err := row.Scan(&post.Id, &post.Author, &post.PostDate, &post.Title, &post.FullText, &post.Slug, &post.LikeCount, &post.DislikeCount)
	if err != nil {
		return nil, err
	}

	return &post, nil
}
