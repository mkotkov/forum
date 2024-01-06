// main.go
package main

import (
	"context"
	"forum/db"
	"forum/handlers"
	"forum/middleware"
	"log"
	"net/http"
	"text/template"
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

	// Middleware initialization
	authMiddleware := middleware.NewRepository(database)
	authMiddleware.SetRepo(authMiddleware)

	// Создание маршрутизатора
	mux := http.NewServeMux()

	// Обработка статических файлов
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// Главная страница
	mux.HandleFunc("/", middleware.AuthenticateHandler(repo, handlers.MainPageHandler(repo)))

	// Стартовая страница (новый обработчик)
	mux.HandleFunc("/start-page", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./templates/header.html", "templates/footer.html", "templates/forum-card.html", "templates/start-page.html", "templates/login-form.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		err = tmpl.ExecuteTemplate(w, "start-page", nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
	})

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
