package repository

import (
	"context"
	"fmt"
	"log"
)

func (r *Repository) DeleteReaction(ctx context.Context, postID, userID int) error {
	_, err := r.db.ExecContext(ctx, `
        DELETE FROM reaction
        WHERE user_id = $1 AND post_id = $2
    `, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to delete reaction: %w", err)
	}
	return nil
}

// Метод для реакции пользователя на пост (лайк/дизлайк)
func (r *Repository) ReactPost(ctx context.Context, postID, userID int, reactionType string) error {
	// Удаление предыдущей реакции пользователя
	err := r.DeleteReaction(ctx, postID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete previous reaction: %w", err)
	}

	// Вставка новой реакции пользователя
	_, err = r.db.ExecContext(ctx, `
        INSERT INTO reaction (user_id, post_id, reaction_type)
        VALUES ($1, $2, $3)
    `, userID, postID, reactionType)
	if err != nil {
		return fmt.Errorf("failed to react to post: %w", err)
	}

	// Обновление счетчика лайков/дизлайков для поста
	var updateColumn string
	if reactionType == "like" {
		updateColumn = "like_count"
	} else if reactionType == "dislike" {
		updateColumn = "dislike_count"
	}

	_, err = r.db.ExecContext(ctx, fmt.Sprintf("UPDATE posts SET %s = %s + 1 WHERE id = $1", updateColumn, updateColumn), postID)
	if err != nil {
		return fmt.Errorf("failed to update %s count: %w", reactionType, err)
	}

	return nil
}

// UpdatePostReactionsCount обновляет счетчик лайков и дизлайков для поста
func (r *Repository) UpdatePostReactionsCount(ctx context.Context, postID int) error {
	// Обновление счетчика лайков для поста
	_, err := r.db.ExecContext(ctx, `
        UPDATE posts
        SET like_count = (SELECT COUNT(*) FROM reaction WHERE post_id = $1 AND reaction_type = 'like'),
            dislike_count = (SELECT COUNT(*) FROM reaction WHERE post_id = $1 AND reaction_type = 'dislike')
        WHERE id = $1
    `, postID)
	if err != nil {
		return fmt.Errorf("failed to update post reactions count: %w", err)
	}

	return nil
}

func (r *Repository) LikePost(ctx context.Context, postID, userID int) error {
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO reaction (user_id, post_id, reaction_type)
        VALUES ($1, $2, 'like')
        ON CONFLICT (user_id, post_id, reaction_type) DO NOTHING
    `, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to like post: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE posts SET like_count = like_count + 1 WHERE id = $1", postID)
	if err != nil {
		log.Printf("Failed to update like count: %v", err)
		return fmt.Errorf("failed to update like count: %w", err)
	}

	return nil
}

// DislikePost увеличивает счетчик дизлайков для поста с заданным ID
func (r *Repository) DislikePost(ctx context.Context, postID, userID int) error {
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO reaction (user_id, post_id, reaction_type)
        VALUES ($1, $2, 'dislike')
        ON CONFLICT (user_id, post_id, reaction_type) DO NOTHING
    `, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to dislike post: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE posts SET dislike_count = dislike_count + 1 WHERE id = $1", postID)
	if err != nil {
		return fmt.Errorf("failed to update dislike count: %w", err)
	}

	return nil
}

// GetPostLikes получает количество лайков для поста
func (r *Repository) GetPostLikes(ctx context.Context, postID int) (int, error) {
	var likeCount int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM reaction WHERE post_id = $1 AND reaction_type = 'like'", postID).Scan(&likeCount)
	if err != nil {
		return 0, fmt.Errorf("failed to get post likes: %w", err)
	}
	return likeCount, nil
}

// GetPostDislikes получает количество дизлайков для поста
func (r *Repository) GetPostDislikes(ctx context.Context, postID int) (int, error) {
	var dislikeCount int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM reaction WHERE post_id = $1 AND reaction_type = 'dislike'", postID).Scan(&dislikeCount)
	if err != nil {
		return 0, fmt.Errorf("failed to get post dislikes: %w", err)
	}
	return dislikeCount, nil
}



func (r *Repository) LikeComment(ctx context.Context, commentID, userID int) error {
    _, err := r.db.ExecContext(ctx, `
        INSERT INTO reaction (user_id, comment_id, reaction_type)
        VALUES ($1, $2, 'like')
        ON CONFLICT (user_id, comment_id, reaction_type) DO UPDATE SET reaction_type = 'like', created_at = CURRENT_TIMESTAMP
    `, userID, commentID)
    if err != nil {
        return fmt.Errorf("failed to dislike comment: %w", err)
    }
    
    // Обновление счетчика дизлайков для комментария
    err = r.UpdateCommentReactionsCount(ctx, commentID)
    if err != nil {
        return fmt.Errorf("failed to update comment reactions count: %w", err)
    }
    
    return nil
}


func (r *Repository) DislikeComment(ctx context.Context, commentID, userID int) error {
    _, err := r.db.ExecContext(ctx, `
        INSERT INTO reaction (user_id, comment_id, reaction_type)
        VALUES ($1, $2, 'dislike')
        ON CONFLICT (user_id, comment_id, reaction_type) DO UPDATE SET reaction_type = 'dislike', created_at = CURRENT_TIMESTAMP
    `, userID, commentID)
    if err != nil {
        return fmt.Errorf("failed to dislike comment: %w", err)
    }
    
    // Обновление счетчика дизлайков для комментария
    err = r.UpdateCommentReactionsCount(ctx, commentID)
    if err != nil {
        return fmt.Errorf("failed to update comment reactions count: %w", err)
    }
    
    return nil
}

func (r *Repository) UpdateCommentReactionsCount(ctx context.Context, commentID int) error {
    _, err := r.db.ExecContext(ctx, `
        UPDATE comments
        SET like_count = (SELECT COUNT(*) FROM reaction WHERE comment_id = $1 AND reaction_type = 'like'),
            dislike_count = (SELECT COUNT(*) FROM reaction WHERE comment_id = $1 AND reaction_type = 'dislike')
        WHERE id = $1
    `, commentID)
    if err != nil {
        return fmt.Errorf("failed to update comment reactions count: %w", err)
    }

    return nil
}

func (r *Repository) DeleteReactionComment(ctx context.Context, commentID, userID int) error {
	_, err := r.db.ExecContext(ctx, `
        DELETE FROM reaction
        WHERE user_id = $1 AND comment_id = $2
    `, userID, commentID)
	if err != nil {
		return fmt.Errorf("failed to delete reaction: %w", err)
	}
	return nil
}