package middleware

import (
	"fmt"
	"net/http"
)

func SavePost(w http.ResponseWriter, r *http.Request, repo RepositoryInterface) {
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Error parsing form", http.StatusBadRequest)
        return
    }

    title := r.FormValue("title")
    fullText := r.FormValue("full-text")
	authorName:= "Max"

    fmt.Printf("Title: %s\n", title)
    fmt.Printf("Full Text: %s\n", fullText)
	fmt.Printf("Name: %s\n", authorName)

    if title == "" || fullText == ""{
        http.Error(w, "Title and Full Text cannot be empty", http.StatusBadRequest)
        return
    }

    err = repo.InsertData(title, fullText, authorName)
	if err != nil {
		fmt.Printf("Error inserting data: %v\n", err)
		http.Error(w, "Error inserting data", http.StatusInternalServerError)
		return
	}

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (r *Repository) InsertData(title, fullText, authorName string) error {
    stmt, err := r.db.Prepare("INSERT INTO posts (title, full_text, author_name) VALUES (?, ?, ?)")
    if err != nil {
        return fmt.Errorf("failed to prepare statement: %w", err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(title, fullText, authorName)
    if err != nil {
        return fmt.Errorf("failed to execute statement: %w", err)
    }

    fmt.Println("Insertion successful")
    return nil
}