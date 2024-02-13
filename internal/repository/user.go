package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	Id             int    `json:"id" db:"id"`
	Login          string `json:"login" db:"login"`
	HashedPassword string `json:"hashed_password" db:"hashed_password"`
	Email          string `json:"email" db:"email"`
	SessionID      string `json:"session_id" db:"session_id"`
}

func (r *Repository) Login(ctx context.Context, login, hashedPassword string) (u User, err error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, login, hashed_password, email FROM users WHERE login = $1 AND hashed_password = $2`, login, hashedPassword)
	if err := row.Scan(&u.Id, &u.Login, &u.HashedPassword, &u.Email); err != nil {
		if err == sql.ErrNoRows {
			return u, fmt.Errorf("user not found")
		}
		return u, fmt.Errorf("failed to query data: %w", err)
	}

	return u, nil
}

func (r *Repository) AddNewUser(ctx context.Context, login, email, hashedPassword string) (err error) {
	_, err = r.db.ExecContext(ctx, `INSERT INTO users (login, email, hashed_password) VALUES ($1, $2, $3)`, login, email, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to exec data: %w", err)
	}

	return nil
}

func (r *Repository) UserExists(ctx context.Context, login, email string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE login = $1 OR email = $2`, login, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return count > 0, nil
}

func GetUserIDFromToken(r *http.Request, repo *Repository) (userID int, err error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Println("Session ID cookie not found")
		return 0, errors.New("session ID cookie not found")
	}

	sessionID := cookie.Value

	user, err := repo.GetUserBySessionID(r.Context(), sessionID)
	if err != nil {
		log.Printf("Failed to get user by session ID: %v\n", err)
		return 0, fmt.Errorf("failed to get user by session ID: %w", err)
	}

	return user.Id, nil
}


func (r *Repository) GetUserLoginByID(ctx context.Context, userID int) (string, error) {
	var login string
	err := r.db.QueryRowContext(ctx, "SELECT login FROM users WHERE id = ?", userID).Scan(&login)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user with ID %d not found", userID)
		}
		return "", fmt.Errorf("failed to get user login by ID: %w", err)
	}
	return login, nil
}

func (r *Repository) GetUserBySessionID(ctx context.Context, sessionID string) (User, error) {
	var user User
	err := r.db.QueryRowContext(ctx, "SELECT id, login, email FROM users WHERE session_id = ?", sessionID).Scan(&user.Id, &user.Login, &user.Email)
	if err != nil {
		return user, fmt.Errorf("failed to get user by session ID: %w", err)
	}
	return user, nil
}

func GetUserIDFromSessionID(r *http.Request, repo *Repository) (userID int, err error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Println("Session ID cookie not found")
		return 0, errors.New("session ID cookie not found")
	}

	sessionID := cookie.Value

	user, err := repo.GetUserBySessionID(r.Context(), sessionID)
	if err != nil {
		log.Printf("Failed to get user by session ID: %v\n", err)
		return 0, fmt.Errorf("failed to get user by session ID: %w", err)
	}

	return user.Id, nil
}

func (r *Repository) SaveSessionID(ctx context.Context, userID int, sessionID string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET session_id = ? WHERE id = ?", sessionID, userID)
	if err != nil {
		return fmt.Errorf("failed to save session ID: %w", err)
	}
	return nil
}

func (r *Repository) DeleteSessionID(ctx context.Context, userID int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET session_id = NULL WHERE id = ?", userID)
	if err != nil {
		return fmt.Errorf("failed to delete session ID: %w", err)
	}
	return nil
}
