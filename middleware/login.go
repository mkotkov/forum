package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"forum/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Login(ctx context.Context, login, hashedPassword string) (u models.User, err error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, login, name, surname FROM users WHERE login = ? AND hashed_password = ?", login, hashedPassword)
	err = row.Scan(&u.Id, &u.Login, &u.Name, &u.Surname)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
	}
	return
}
