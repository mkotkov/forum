package application

import (
	"errors"
	"net/http"
)

func readCookie(name string, r *http.Request) (value string, err error) {
	if name == "" {
        return value, errors.New("you are trying to read empty cookie")
    }
    cookie, err := r.Cookie(name)
    if err != nil {
        return value, err // Возвращаем ошибку, если куки не найдено
    }
    value = cookie.Value
    return value, err
}

