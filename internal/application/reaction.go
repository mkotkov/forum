package application

import (
	"log"
	"net/http"
	"strconv"
)

func (a *App) ReactPost(w http.ResponseWriter, r *http.Request, slug string, reactionType string, isAuthorized bool) {
    // Проверка аутентификации пользователя
    user, err := a.getAuthenticatedUser(r)
    if err != nil {
        // Обработка ошибки аутентификации
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    post, err := a.repo.GetPostBySlug(a.ctx, slug)
    if err != nil {
        http.Error(w, "Error getting post", http.StatusInternalServerError)
        return
    }

    // Удаление соответствующего состояния пользователя
    err = a.repo.DeleteReaction(a.ctx, int(post.Id), int(user.Id))
    if err != nil {
        http.Error(w, "Error deleting previous reaction", http.StatusInternalServerError)
        return
    }

    // Сохранение новой реакции пользователя
    err = a.repo.ReactPost(a.ctx, int(post.Id), int(user.Id), reactionType)
    if err != nil {
        http.Error(w, "Error reacting to post", http.StatusInternalServerError)
        return
    }

    // Обновление счетчиков лайков и дизлайков в посте
    err = a.repo.UpdatePostReactionsCount(a.ctx, int(post.Id))
    if err != nil {
        http.Error(w, "Error updating post reactions count", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}


func (a *App) LikeComment(w http.ResponseWriter, r *http.Request, commentID string, isAuthorized bool) {
  // Проверка аутентификации пользователя
  sessionID, err := readCookie("session_id", r)
  if err != nil {
	  // Обработка ошибки аутентификации
	  http.Error(w, "Unauthorized", http.StatusUnauthorized)
	  return
  }

  // Получение пользователя из кэша
  user, ok := a.cache[sessionID]
  if !ok {
	  // Обработка ошибки аутентификации
	  http.Error(w, "Unauthorized", http.StatusUnauthorized)
	  return
  }

	// Преобразование commentID в целочисленный формат
	commentIDInt, err := strconv.Atoi(commentID)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	err = a.repo.DeleteReactionComment(a.ctx, commentIDInt, int(user.Id))
	if err != nil {
		log.Printf("Error deleting previous reaction for LikeComment: %v", err)
		http.Error(w, "Error deleting previous reaction", http.StatusInternalServerError)
		return
	}

	// Вызов метода репозитория для постановки лайка к комментарию
	err = a.repo.LikeComment(a.ctx, commentIDInt, int(user.Id))
	if err != nil {
		log.Printf("Error liking comment: %v", err)
		http.Error(w, "Error liking comment", http.StatusInternalServerError)
		return
	}

	// Перенаправление пользователя обратно на страницу комментариев
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func (a *App) DislikeComment(w http.ResponseWriter, r *http.Request, commentID string, isAuthorized bool) {
  // Проверка аутентификации пользователя
  sessionID, err := readCookie("session_id", r)
  if err != nil {
	  // Обработка ошибки аутентификации
	  http.Error(w, "Unauthorized", http.StatusUnauthorized)
	  return
  }

  // Получение пользователя из кэша
  user, ok := a.cache[sessionID]
  if !ok {
	  // Обработка ошибки аутентификации
	  http.Error(w, "Unauthorized", http.StatusUnauthorized)
	  return
  }

	// Преобразование commentID в целочисленный формат
	commentIDInt, err := strconv.Atoi(commentID)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	err = a.repo.DeleteReactionComment(a.ctx, commentIDInt, int(user.Id))
	if err != nil {
		log.Printf("Error deleting previous reaction for DislikeComment: %v", err)
		http.Error(w, "Error deleting previous reaction", http.StatusInternalServerError)
		return
	}

	// Вызов метода репозитория для постановки дизлайка к комментарию
	err = a.repo.DislikeComment(a.ctx, commentIDInt, int(user.Id))
	if err != nil {
		log.Printf("Error disliking comment: %v", err)
		http.Error(w, "Error disliking comment", http.StatusInternalServerError)
		return
	}

	// Перенаправление пользователя обратно на страницу комментариев
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}
