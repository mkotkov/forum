package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Comments struct {
	ID           int       `db:"id"`
	UserID       int       `db:"user_id"`
	PostID       int       `db:"post_id"`
	UserName     string    `db:"user_name"`
	Comment      string    `db:"comment"`
	CommentDate  time.Time `db:"comment_date"`
	LikeCount    int       `db:"like_count"`
	DislikeCount int       `db:"dislike_count"`
}

func (r *Repository) GetCommentsByPostID(ctx context.Context, postID uint16) ([]Comments, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, post_id, author_id, user_name, text, comment_date, like_count, dislike_count FROM comments WHERE post_id = $1 ORDER BY comment_date DESC
	`, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	var comments []Comments
	for rows.Next() {
		var comment Comments
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.UserName, &comment.Comment, &comment.CommentDate, &comment.LikeCount, &comment.DislikeCount); err != nil {
			return nil, fmt.Errorf("failed to scan comments: %w", err)
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *Repository) SaveComment(ctx context.Context, postID, userID int, userName, commentText string) error {
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO comments (post_id, author_id, user_name, text)
        VALUES ($1, $2, $3, $4)
    `, postID, userID, userName, commentText)
	if err != nil {
		return fmt.Errorf("failed to save comment: %w", err)
	}
	return nil
}

func (r *Repository) GetPostByPostID(ctx context.Context, postID uint16) (Posts, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, author, title, post_date, full_text, slug FROM posts WHERE id = $1`, postID)
	var p Posts
	if err := row.Scan(&p.Id, &p.Author, &p.Title, &p.PostDate, &p.FullText, &p.Slug); err != nil {
		if err == sql.ErrNoRows {
			return p, fmt.Errorf("post not found")
		}
		return p, fmt.Errorf("failed to query data: %w", err)
	}

	return p, nil
}
