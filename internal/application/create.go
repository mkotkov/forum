package application

import (
	"forum/internal/repository"
	"fmt"
	"net/http"
	"text/template"
)

func (a *App) Create(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("public/html/create.html", "public/html/header.html", "public/html/footer.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "create", nil)
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

	if title == "" || fullText == "" {
		http.Error(w, "Title and Full Text cannot be empty", http.StatusBadRequest)
		return
	}

	// Insert the post data into the database using the authorized user's information
	err = a.InsertData(title, fullText, authorizedUser.Login)
	if err != nil {
		fmt.Printf("Error inserting data: %v\n", err)
		http.Error(w, "Error inserting data", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *App) InsertData(title, fullText, authorName string) error {
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

	stmt, err := tx.Prepare("INSERT INTO posts (title, full_text, author, post_date) VALUES (?, ?, ?, CURRENT_TIMESTAMP)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(title, fullText, user.Login)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	fmt.Println("Insertion successful")
	return nil
}

