package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type Posts struct {
	Id            uint16
	Author        string
	PostDate      string
	Title         string
	FullText      string
	Slug          string
	LikeCount     int
	DislikeCount  int
}

const (
	SQLSelectAllPosts       = "SELECT id, author, post_date, title, full_text, slug, like_count, dislike_count FROM posts ORDER BY post_date DESC"
	SQLSelectMostRecentPost = "SELECT id, author, post_date, title, full_text, slug, like_count, dislike_count FROM posts ORDER BY post_date DESC LIMIT 1"
)

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

func (r *Repository) GetPostBySlug(ctx context.Context, slug string) (p Posts, err error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, author, title, post_date, full_text, slug FROM posts WHERE slug = $1`, slug)
	if err := row.Scan(&p.Id, &p.Author, &p.Title, &p.PostDate, &p.FullText, &p.Slug); err != nil {
		if err == sql.ErrNoRows {
			return p, fmt.Errorf("post not found")
		}
		return p, fmt.Errorf("failed to query data: %w", err)
	}

	return p, nil
}
