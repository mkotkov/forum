package handlers

import (
	"database/sql"
	"fmt"
	"forum/middleware"
	"forum/models"
	"net/http"
	"text/template"
)

// MainPageHandler обрабатывает главную страницу
func MainPageHandler(repo *middleware.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html", "templates/forum-card.html", "templates/start-page.html", "templates/login-form.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if the user is authenticated
		isAuth := middleware.IsAuthenticated(repo.GetDB(), r)

		// Check if there are any posts in the database
		posts, err := queryData(repo.GetDB())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(posts) > 0 || isAuth {
			// There are posts or the user is authenticated, load the "index" page
			fmt.Println("index ok")
			data := struct {
				Posts []models.Posts
			}{
				Posts: posts,
			}
			err = tmpl.ExecuteTemplate(w, "index", data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			fmt.Println("start-page")
			// User is not authenticated, load the "start-page"
			err = tmpl.ExecuteTemplate(w, "start-page", nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
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


