package middleware

import (
	"fmt"
	"net/http"
)


func SavePost(w http.ResponseWriter, r *http.Request, repo RepositoryInterface) {

    if !IsAuthenticated(repo.GetDB(), r) {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

    if repo == nil {
        http.Error(w, "Repository is nil", http.StatusInternalServerError)
        return
    }

    // Получаем пользователя из сессии
    user, err := getUserBySessionToken(repo.GetDB(), extractSessionToken(r))
    if err != nil {
        // Обработка ошибки аутентификации
        http.Error(w, "User not authenticated", http.StatusUnauthorized)
        return
    }

    err = r.ParseForm()
    if err != nil {
        http.Error(w, "Error parsing form", http.StatusBadRequest)
        return
    }

    title := r.FormValue("title")
    fullText := r.FormValue("full-text")
    authorName := user.Name

    fmt.Printf("Title: %s\n", title)
    fmt.Printf("Full Text: %s\n", fullText)
    fmt.Printf("Name: %s\n", authorName)

    if title == "" || fullText == "" {
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
    tx, err := r.db.Begin()
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }

    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p) // Повторно взлетаем панику после отката транзакции
        } else if err != nil {
            tx.Rollback()
        } else {
            err = tx.Commit()
            if err != nil {
                fmt.Printf("Error committing transaction: %v\n", err)
            }
        }
    }()

    stmt, err := tx.Prepare("INSERT INTO posts (title, full_text, author_name) VALUES (?, ?, ?)")
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

func extractSessionToken(r *http.Request) string {
    cookie, err := r.Cookie(sessionCookieName)
    if err == nil && cookie != nil {
        return cookie.Value
    }
    return ""
}