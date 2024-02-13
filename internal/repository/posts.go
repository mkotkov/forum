package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type Posts struct {
	Id           uint16
	Author       string
	PostDate     string
	Title        string
	FullText     string
	Slug         string
	LikeCount    int
	DislikeCount int
	Topic        string
	TopicID      int
}

type Topic struct {
	ID   int
	Name string
}

const (
	SQLSelectAllPosts = "SELECT id, author, post_date, title, full_text, slug, like_count, dislike_count, COALESCE(topic_id, -1) FROM posts"

	SQLSelectMostRecentPost = "SELECT id, author, post_date, title, full_text, slug, like_count, dislike_count, COALESCE(topic_id, -1) FROM posts ORDER BY post_date DESC LIMIT 1"

	SQLSelectAllTopics = "SELECT id, name FROM topics ORDER BY name"
)

func (r *Repository) GetUserByName(ctx context.Context, login string) (u User, err error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, login FROM users WHERE login = $1`, login)

	err = row.Scan(&u.Id, &u.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, fmt.Errorf("user not found")
		}
		return u, fmt.Errorf("failed to query data: %w", err)
	}

	return u, nil
}

func (r *Repository) GetPostBySlug(ctx context.Context, slug string) (p Posts, err error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, author, title, post_date, full_text, slug, COALESCE(topic_id, -1) FROM posts WHERE slug = $1`, slug)
	if err := row.Scan(&p.Id, &p.Author, &p.Title, &p.PostDate, &p.FullText, &p.Slug, &p.TopicID); err != nil {
		if err == sql.ErrNoRows {
			return p, fmt.Errorf("post not found")
		}
		return p, fmt.Errorf("failed to query data: %w", err)
	}

	// Добавим вывод для проверки
	fmt.Printf("Post from GetPostBySlug: %+v\n", p)

	return p, nil
}

func (r *Repository) GetAllTopics(ctx context.Context) ([]Topic, error) {
	rows, err := r.db.QueryContext(ctx, SQLSelectAllTopics)
	if err != nil {
		return nil, fmt.Errorf("failed to query topics: %w", err)
	}
	defer rows.Close()

	var topics []Topic
	for rows.Next() {
		var topic Topic
		if err := rows.Scan(&topic.ID, &topic.Name); err != nil {
			return nil, fmt.Errorf("failed to scan topic: %w", err)
		}
		topics = append(topics, topic)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return topics, nil
}
