package application

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"text/template"

	"forum/internal/repository"
)

func (a *App) Create(w http.ResponseWriter, r *http.Request) {
    // Fetch all topics for the "create" page
    topics, err := a.repo.GetAllTopics(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Include topics in the data for rendering the "create" template
    data := struct {
        Topics []repository.Topic
    }{
        Topics: topics,
    }

    tmpl, err := template.ParseFiles(
        "public/html/create.html",
        "public/html/header.html",
        "public/html/footer.html",
    )
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    err = tmpl.ExecuteTemplate(w, "create", data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}


func (a *App) SavePost(w http.ResponseWriter, r *http.Request) {
	// Ensure the user is authenticated
	token, err := readCookie("token", r)
	if err != nil || a.cache[token] == (repository.User{}) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Retrieve the authenticated user
	authorizedUser := a.cache[token]

    title := r.FormValue("title")
    fullText := r.FormValue("full-text")
    topicID := r.FormValue("topic")

    if title == "" || fullText == "" || topicID == "" {
        http.Error(w, "Title, Full Text, and Topic cannot be empty", http.StatusBadRequest)
        return
    }

    slug := generateSlug(title) + topicID 

    // Insert the post data into the database using the authorized user's information and the selected topic
    err = a.InsertData(title, fullText, authorizedUser.Login, slug, topicID)
    if err != nil {
        fmt.Printf("Error inserting data: %v\n", err)
        http.Error(w, "Error inserting data", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, fmt.Sprintf("/post/%s", slug), http.StatusSeeOther)
}

func (a *App) InsertData(title, fullText, authorName, slug, topicID string) error {
	// Retrieve the user information based on the author's name
	user, err := a.repo.GetUserByName(a.ctx, authorName)
    if err != nil {
        return fmt.Errorf("failed to retrieve user: %w", err)
    }

	tx, err := a.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-panic after rolling back the transaction
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				fmt.Printf("Error committing transaction: %v\n", err)
			}
		}
	}()

	stmt, err := tx.Prepare("INSERT INTO posts (title, full_text, author, post_date, slug, topic_id) VALUES (?, ?, ?, CURRENT_TIMESTAMP, ?, ?)")
    if err != nil {
        return fmt.Errorf("failed to prepare statement: %w", err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(title, fullText, user.Login, slug, topicID)
    if err != nil {
        return fmt.Errorf("failed to execute statement: %w", err)
    }

    fmt.Println("Insertion successful")
    return nil
}

func generateSlug(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")

	// Remove special characters and symbols
	reg, err := regexp.Compile("[^a-z0-9-]+")
	if err != nil {
		// handle error
		return ""
	}
	s = reg.ReplaceAllString(s, "")

	// Truncate if needed
	if len(s) > 50 {
		s = s[:50]
	}

	return s
}
