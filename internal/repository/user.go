package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type User struct {
	Id             int    `json:"id" db:"id"`
	Login          string `json:"login" db:"login"`
	HashedPassword string `json:"hashed_password" db:"hashed_password"`
	Name           string `json:"name" db:"name"`
	Surname        string `json:"surname" db:"surname"`
}

func (r *Repository) Login(ctx context.Context, login, hashedPassword string) (u User, err error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, login, name, surname FROM users WHERE login = $1 AND hashed_password = $2`, login, hashedPassword)
	if err := row.Scan(&u.Id, &u.Login, &u.Name, &u.Surname); err != nil {
		if err == sql.ErrNoRows {
			return u, fmt.Errorf("user not found")
		}
		return u, fmt.Errorf("failed to query data: %w", err)
	}

	return u, nil
}

func (r *Repository) AddNewUser(ctx context.Context, name, surname, login, hashedPassword string) (err error) {
	_, err = r.db.ExecContext(ctx, `INSERT INTO users (name, surname, login, hashed_password) VALUES ($1, $2, $3, $4)`, name, surname, login, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to exec data: %w", err)
	}

	return nil
}

