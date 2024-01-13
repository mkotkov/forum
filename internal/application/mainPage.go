package application

import (
	"forum/internal/repository"
	"database/sql"
	"html/template"
	"net/http"
)

func (a *App) MainPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"public/html/header.html", 
		"public/html/footer.html", 
		"public/html/forum-card.html", 
		"public/html/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	posts, err := queryData(a.repo.GetDB())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// There are posts, load the "index" page
	data := struct {
		Posts []repository.Posts
	}{
		Posts: posts,
	}
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
        err := rows.Scan(&post.Id, &post.Author, &post.Title, &post.PostDate, &post.FullText)
        if err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }

    return posts, nil
}

