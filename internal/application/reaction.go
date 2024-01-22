package application

import (
	"net/http"
)

func (a *App) ReactPost(w http.ResponseWriter, r *http.Request, slug string, reactionType string) {
	post, err := a.repo.GetPostBySlug(a.ctx, slug)
	if err != nil {
		http.Error(w, "Error getting post", http.StatusInternalServerError)
		return
	}

	user, err := a.getAuthenticatedUser(r)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Удаление предыдущей реакции пользователя
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

	http.Redirect(w, r, "/post/"+slug, http.StatusSeeOther)
}
