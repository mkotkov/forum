package application

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func (a *App) Signup(w http.ResponseWriter, r *http.Request) {
	login := strings.TrimSpace(r.FormValue("login"))
	password := strings.TrimSpace(r.FormValue("password"))
	password2 := strings.TrimSpace(r.FormValue("password2"))
	email := strings.TrimSpace(r.FormValue("email"))

	if login == "" || password == "" || email == "" {
		a.UnregPage(w, r, "<div class="+"error"+"><p>All fields must be filled in!</p></div>", false)
		return
	}

	if password != password2 {
		a.UnregPage(w, r, "<div class="+"error"+"><p>Passwords do not match! Please try again.</p></div>", false)
		return
	}

	// Check if the user already exists
	exists, err := a.repo.UserExists(a.ctx, login, email)
	if err != nil {
		a.UnregPage(w, r, fmt.Sprintf("<div class="+"error"+"><p>Error checking user existence: %v</p></div>", err), false)
		return
	}

	if exists {
		a.UnregPage(w, r, "<div class="+"error"+"><p>User with this login or email already exists. Please log in or use different credentials.</p></div>", false)
		return
	}

	if !isValidEmail(email) {
		a.UnregPage(w, r, "<div class="+"error"+"><p>Invalid email format. Please enter a valid email address.</p></div>", false)
		return
	}

	hash := md5.Sum([]byte(password))
	hashedPass := hex.EncodeToString(hash[:])

	err = a.repo.AddNewUser(a.ctx, login, email, hashedPass)
	if err != nil {
		a.UnregPage(w, r, fmt.Sprintf("<div class="+"error"+"><p>Error creating user: %v</p></div>", err), false)
		return
	}

	a.UnregPage(w, r, fmt.Sprintf("<div class="+"error"+"><p>%s, you have successfully registered! Now you can log in.</p></div>", login), false)
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
