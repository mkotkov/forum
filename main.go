package main

import (
	"context"
	"forum/db"
	"forum/handlers"
	"forum/middleware"
	"log"
	"net/http"
)

func main() {
	// Инициализация базы данных
	database, err := db.InitDBConn(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Создание репозитория с использованием экземпляра *sql.DB
	repo := middleware.NewRepository(database)
	repo.SetRepo(repo)

	// Создание маршрутизатора
	mux := http.NewServeMux()

	// Обработка статических файлов
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// Главная страница
	mux.HandleFunc("/", handlers.MainPageHandler(repo))

	// Log out
	mux.HandleFunc("/logout/", middleware.LogoutHandler)

	// Создание поста
	mux.HandleFunc("/create/", handlers.CreateHandler(repo))

	// Сохранение поста
	mux.HandleFunc("/save_post/", func(w http.ResponseWriter, r *http.Request) {
		middleware.SavePost(w, r, repo)
	})

	// Логин
	mux.HandleFunc("/login/", func(w http.ResponseWriter, r *http.Request) {
		repo.HandleLogin(w, r)
	})

	// Дополнительные маршруты...

	// Запуск сервера
	http.ListenAndServe(":8080", mux)
}
