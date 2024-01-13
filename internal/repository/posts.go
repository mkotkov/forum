package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type Posts struct {
    Id       uint16
    Author   string
    Title    string
    PostDate string
    FullText string
}

const SQLSelectAllPosts = "SELECT * FROM posts"

func (r *Repository) GetUserByName(ctx context.Context, login string) (u User, err error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, login, name, surname FROM users WHERE login = $1`, login)

	err = row.Scan(&u.Id, &u.Login, &u.Name, &u.Surname)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, fmt.Errorf("user not found")
		}
		return u, fmt.Errorf("failed to query data: %w", err)
	}

	return u, nil
}


