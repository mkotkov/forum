package middleware

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/models"
	"net/http"
	"time"
)

const sessionCookieName = "session_token"


func AuthenticateHandler(repo RepositoryInterface, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверка аутентификации пользователя
		if !IsAuthenticated(repo.GetDB(), r) {
			// Пользователь не аутентифицирован, перенаправляем на страницу входа
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Пользователь аутентифицирован, продолжаем выполнение следующего обработчика
		next(w, r)
	}
}


func AuthenticateUser(w http.ResponseWriter, r *http.Request, user *models.User) {
	// Сохраняем токен среди данных пользователя
	user.SessionToken = "some_unique_session_token"
	// Создаем cookie и добавляем туда токен
	http.SetCookie(w, &http.Cookie{
		Name:  sessionCookieName,
		Value: user.SessionToken,
		Path:  "/",
	})
}

func ClearSessionCookie(w http.ResponseWriter) {
    http.SetCookie(w, &http.Cookie{
        Name:    sessionCookieName,
        Value:   "",
        Expires: time.Unix(0, 0), 
        Path:    "/",
    })
}


// IsAuthenticated checks if the user is authenticated based on the session cookie.
func IsAuthenticated(db *sql.DB, r *http.Request) bool {
    cookie, err := r.Cookie(sessionCookieName)
    if err == nil && cookie != nil {
        user, err := getUserBySessionToken(db, cookie.Value)
        fmt.Println("sAuthenticated OK!")
        return err == nil && user != nil
    }
    
    return false
}

func getUserBySessionToken(db *sql.DB, token string) (*models.User, error) {
    // Используйте ваш SQL-запрос для поиска пользователя по токену сессии
    row := db.QueryRow("SELECT id, login, name, surname FROM user WHERE session_token = ?", token)
    
    var user models.User
    err := row.Scan(&user.Id, &user.Login, &user.Name, &user.Surname)
    if err != nil {
        // Если пользователя не найдено, вернуть nil и ошибку.
        return nil, errors.New("user not found")
    }

    return &user, nil
}



func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    ClearSessionCookie(w)
    fmt.Println("logout")

    // Добавление HTTP-заголовка для предотвращения кэширования
    w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

    http.Redirect(w, r, "/", http.StatusSeeOther)
}