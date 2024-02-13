package application

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"forum/internal/repository"

	"github.com/google/uuid"
)

func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login-inp")
	password := r.FormValue("password-inp")

	fmt.Println("Login and password provided:", login, password)

	if login == "" || password == "" {
		fmt.Println("Empty login or password provided")
		a.UnregPage(w, r, "<div class="+"error"+"><p>You must provide a username and password!</p></div>", false)
		return
	}

	hash := md5.Sum([]byte(password))
	hashedPass := hex.EncodeToString(hash[:])

	user, err := a.repo.Login(a.ctx, login, hashedPass)
	if err != nil {
		fmt.Println("Error during login:", err)
		a.UnregPage(w, r, "<div class="+"error"+"><p>Failed to login!</p></div>", false)
		return
	}

	// Create a session ID
	sessionID := uuid.New().String()

	fmt.Println("User session ID:", sessionID)

	// Save the session ID in the database
	err = a.repo.SaveSessionID(a.ctx, user.Id, sessionID)
	if err != nil {
		fmt.Println("Error saving session ID:", err)
		a.UnregPage(w, r, "<div class="+"error"+"><p>Failed to create a session!</p></div>", false)
		return
	}

	// Set session ID cookie
	cookie := http.Cookie{Name: "session_id", Value: sessionID, Expires: time.Now().Add(24 * time.Hour)}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}


func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	// Delete all cookies
	for _, v := range r.Cookies() {
		c := http.Cookie{
			Name:   v.Name,
			MaxAge: -1,
		}
		http.SetCookie(w, &c)
	}

	fmt.Println("All cookies deleted")

	http.Redirect(w, r, "/", http.StatusSeeOther)

	// Get user ID from session ID
	userID, err := repository.GetUserIDFromSessionID(r, a.repo)
	if err != nil {
		fmt.Println("Error getting user ID from session ID:", err)
		a.UnregPage(w, r, "<div class="+"error"+"><p>Failed to get user ID from session ID!</p></div>", false)
		return
	}

	fmt.Println("User ID obtained from session ID:", userID)

	// Delete session ID from the database
	err = a.repo.DeleteSessionID(a.ctx, userID)
	if err != nil {
		fmt.Println("Error deleting session ID:", err)
		a.UnregPage(w, r, "<div class="+"error"+"><p>Failed to delete session ID!</p></div>", false)
		return
	}

	fmt.Println("Session ID deleted from the database")
}

