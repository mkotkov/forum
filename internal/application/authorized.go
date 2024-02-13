package application

import (
	"fmt"
	"net/http"

	"forum/internal/repository"
)

func (a *App) authorized(handler func(http.ResponseWriter, *http.Request, bool)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        sessionID, err := readCookie("session_id", r) // Чтение куки сессии
        if err != nil {
            fmt.Println("Error reading session ID cookie:", err)
            handler(w, r, false)
            return
        }

        userID, err := repository.GetUserIDFromSessionID(r, a.repo)
        if err != nil {
            fmt.Println("Error getting user ID from session ID:", err)
            handler(w, r, false)
            return
        }

        // Проверка, что существует пользователь с данным sessionID
        if userID == 0 {
            fmt.Println("User with session ID", sessionID, "not found")
            handler(w, r, false)
            return
        }

        // Используем sessionID, например, для журналирования
        fmt.Println("User session ID:", sessionID)

        handler(w, r, true)
    })
}
