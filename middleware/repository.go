package middleware

import (
    "context"
    "forum/models"
)

type RepositoryInterface interface {
    Login(ctx context.Context, login, hashedPassword string) (u models.User, err error)
    InsertData(title, fullText, authorName string) error
}
