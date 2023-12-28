package handlers

import (
	"database/sql"
	"forum/models"
	"net/http"
	"text/template"
)

func MainPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html", "templates/forum-card.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db, err := sql.Open("sqlite3", "./db/data.db")
    if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
    defer db.Close()

	posts, err := queryData(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Posts []models.Posts
	}{
		Posts: posts,
	}

	err = tmpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func queryData(db *sql.DB) ([]models.Posts, error) {
    rows, err := db.Query("SELECT * FROM `posts`")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []models.Posts

    for rows.Next() {
        var post models.Posts
        err := rows.Scan(&post.Id, &post.Author, &post.PostDate, &post.Title, &post.FullText)
        if err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }

    return posts, nil
}