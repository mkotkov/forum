package application

import (
	"fmt"
	"forum/internal/repository"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func (a *App) PostPage(w http.ResponseWriter, r *http.Request, slug string) {
    post, err := a.repo.GetPostBySlug(a.ctx, slug)
    if err != nil {
        log.Printf("Error getting post by slug: %v", err)
        if strings.Contains(err.Error(), "post not found") {
            http.Error(w, "Post not found", http.StatusNotFound)
        } else {
            http.Error(w, "Error getting post", http.StatusInternalServerError)
        }
        return
    }

    comments, err := a.repo.GetCommentsByPostID(a.ctx, post.Id)
    if err != nil {
        log.Printf("Error getting comments: %v", err)
        http.Error(w, "Error getting comments", http.StatusInternalServerError)
        return
    }

    tmpl, err := template.ParseFiles(
        "public/html/header.html",
        "public/html/footer.html",
        "public/html/comment.html",
        "public/html/post.html",
    )
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    data := struct {
        Post     repository.Posts
        Comments []repository.Comments
    }{
        Post:     post,
        Comments: comments,
    }

    err = tmpl.ExecuteTemplate(w, "post", data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}


func (a *App) SaveComment(w http.ResponseWriter, r *http.Request) {
    // Извлекаем slug из URL
    slug := strings.TrimPrefix(r.URL.Path, "/save_comment/")

    // Получаем данные из формы
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Error parsing form data", http.StatusBadRequest)
        return
    }

    // Получаем текст комментария и текущего авторизованного пользователя
    commentText := r.FormValue("add-text-comment")
    user, err := a.getAuthenticatedUser(r)
    if err != nil {
        http.Error(w, "User not authenticated", http.StatusUnauthorized)
        return
    }

    // Получаем пост по slug
    post, err := a.repo.GetPostBySlug(a.ctx, slug)
    if err != nil {
        log.Printf("Error getting post by slug: %v", err)
        http.Error(w, "Error getting post", http.StatusInternalServerError)
        return
    }

    // Сохраняем комментарий в базе данных, используя post.ID
    err = a.repo.SaveComment(a.ctx, int(post.Id), int(user.Id), user.Login, commentText)
    if err != nil {
        log.Printf("Error saving comment: %v", err)
        http.Error(w, "Error saving comment", http.StatusInternalServerError)
        return
    }

    // Редиректим обратно на страницу поста с использованием slug
    http.Redirect(w, r, "/post/"+slug, http.StatusSeeOther)
}



func (a *App) getAuthenticatedUser(r *http.Request) (repository.User, error) {
    // Получаем куки из запроса
    cookie, err := r.Cookie("token")
    if err != nil {
        return repository.User{}, fmt.Errorf("no authentication token found")
    }

    // Получаем пользователя из кэша по токену
    user, ok := a.cache[cookie.Value]
    if !ok {
        return repository.User{}, fmt.Errorf("user not authenticated")
    }

    return user, nil
}

