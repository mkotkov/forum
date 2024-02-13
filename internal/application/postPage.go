package application

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/repository"
)

func (a *App) PostPage(w http.ResponseWriter, r *http.Request, slug string, isAuthorized bool, message string) {
	userID, err := repository.GetUserIDFromToken(r, a.repo)
	if err != nil {
	}

	fmt.Println("userID:", userID)

	userLogin, err := a.repo.GetUserLoginByID(r.Context(), userID)
	if err != nil {
	}

	fmt.Println("userLogin:", userLogin)

	post, err := a.repo.GetPostBySlug(a.ctx, slug)
	if err != nil {
		http.Error(w, "Error getting post", http.StatusInternalServerError)
		return
	}

	// fmt.Printf("Post: %+v\n", post)

	if post.TopicID != 0 {
		post.Topic = GetTopicName(a.repo.GetDB(), post.TopicID)
	} else {
		post.Topic = "No Topic"
	}

	likeCount, err := a.repo.GetPostLikes(a.ctx, int(post.Id))
	if err != nil {
		http.Error(w, "Error getting like count", http.StatusInternalServerError)
		return
	}

	dislikeCount, err := a.repo.GetPostDislikes(a.ctx, int(post.Id))
	if err != nil {
		http.Error(w, "Error getting dislike count", http.StatusInternalServerError)
		return
	}

	post.LikeCount = likeCount
	post.DislikeCount = dislikeCount

	comments, err := a.repo.GetCommentsByPostID(a.ctx, post.Id)
	if err != nil {
		http.Error(w, "Error getting comments", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(
		"public/html/header.html",
		"public/html/unreg-header.html",
		"public/html/login-form.html",
		"public/html/signup.html",
		"public/html/footer.html",
		"public/html/comment.html",
		"public/html/post.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	topics, err := a.repo.GetAllTopics(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Println("Post:", post)
	// fmt.Println("Comments:", comments)

	data := struct {
		Post         repository.Posts
		Comments     []repository.Comments
		Topics       []repository.Topic
		IsAuthorized bool
		Message      string
		NameUser     string
	}{
		Post:         post,
		Comments:     comments,
		Topics:       topics,
		IsAuthorized: isAuthorized,
		Message:      message,
		NameUser:     userLogin,
	}

	err = tmpl.ExecuteTemplate(w, "post", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) SaveComment(w http.ResponseWriter, r *http.Request) {
    // Извлекаем slug из URL-пути
    slug := strings.TrimPrefix(r.URL.Path, "/save_comment/")

    // Парсим форму из запроса
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Error parsing form data", http.StatusBadRequest)
        return
    }

    // Получаем текст комментария из формы
    commentText := strings.TrimSpace(r.Form.Get("add-text-comment"))

    // Проверяем, что комментарий не пустой
    if commentText == "" {
        errorMes := "Comment cannot be empty!"
        a.PostPage(w, r, slug, true, errorMes)
        return
    }

    // Получаем аутентифицированного пользователя
    user, err := a.getAuthenticatedUser(r)
    if err != nil {
        http.Error(w, "User not authenticated", http.StatusUnauthorized)
        return
    }

    // Получаем пост по его slug
    post, err := a.repo.GetPostBySlug(a.ctx, slug)
    if err != nil {
        log.Printf("Error getting post by slug: %v", err)
        http.Error(w, "Error getting post", http.StatusInternalServerError)
        return
    }

    // Сохраняем комментарий в базе данных
    err = a.repo.SaveComment(a.ctx, int(post.Id), int(user.Id), user.Login, commentText)
    if err != nil {
        log.Printf("Error saving comment: %v", err)
        http.Error(w, "Error saving comment", http.StatusInternalServerError)
        return
    }

    // Перенаправляем пользователя на страницу поста
    http.Redirect(w, r, "/post/"+slug, http.StatusSeeOther)
}



func (a *App) getAuthenticatedUser(r *http.Request) (repository.User, error) {
	// Получение идентификатора пользователя из сессии
	userID, err := repository.GetUserIDFromSessionID(r, a.repo)
	if err != nil {
		fmt.Println("Error getting user ID from session ID:", err)
		return repository.User{}, fmt.Errorf("user not authenticated")
	}

	userIDStr := strconv.Itoa(userID)

	// Получение пользователя из репозитория по идентификатору
	user, err := a.repo.GetUserBySessionID(r.Context(), userIDStr)
	if err != nil {
		fmt.Println("Error getting user by session ID:", err)
		return repository.User{}, fmt.Errorf("user not found")
	}

	return user, nil
}







