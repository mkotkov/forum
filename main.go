// main.go
package main

import (
	"context"
	"fmt"
	"net/http"
	"forum/db"
	"forum/handlers"
	"forum/middleware"
)

func main() {
	// Инициализация базы данных
	ctx := context.Background()
	db, err := db.InitDBConn(ctx)
	if err != nil {
		fmt.Println("Failed to initialize database:", err)
		return
	}
	defer db.Close()

	// Создание репозитория с использованием экземпляра *sql.DB
	repo := middleware.NewRepository(db)

	// Setting up routes and starting the server
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	http.HandleFunc("/", handlers.MainPage)
	http.HandleFunc("/create/", handlers.Create)
	http.HandleFunc("/save_post/", func(w http.ResponseWriter, r *http.Request) {
		middleware.SavePost(w, r, repo)
	})
	http.ListenAndServe(":8080", nil)
}
