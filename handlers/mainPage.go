package handlers

import (
	"database/sql"
	"forum/middleware"
	"forum/models"
	"net/http"
	"text/template"
)

func MainPageHandler(repo *middleware.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html", "templates/forum-card.html", "templates/start-page.html", "templates/login-form.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if the user is authenticated
		if middleware.IsAuthenticated(repo.GetDB(), r) {
			// User is authenticated, load the "index" page
			posts, err := queryData(repo.GetDB())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if len(posts) > 0 {
				// There are posts, load the "index" page
				data := struct {
					Posts []models.Posts
				}{
					Posts: posts,
				}
				err = tmpl.ExecuteTemplate(w, "index", data)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				// There are no posts, load the "start-page"
				err = tmpl.ExecuteTemplate(w, "start-page", nil)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		} else {
			// User is not authenticated, load the "start-page"
			err = tmpl.ExecuteTemplate(w, "start-page", nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
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
