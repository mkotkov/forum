package middleware

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"forum/models"
	"log"
	"net/http"
)

type RepositoryInterface interface {
	GetDB() *sql.DB
	Login(login, hashedPassword string) (models.User, error)
	InsertData(title, fullText, authorName string) error
	HandleLogin(w http.ResponseWriter, req *http.Request)
}

type Repository struct {
	db       *sql.DB
	ctx      context.Context
	repo     *Repository
	LoggedIn bool
}

// InsertData implements RepositoryInterface.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) SetRepo(repo *Repository) {
	r.repo = repo
}

func (r *Repository) Login(login, hashedPassword string) (models.User, error) {
	if r == nil || r.repo == nil {
		return models.User{}, fmt.Errorf("repository is nil")
	}

	row := r.db.QueryRowContext(r.ctx, "SELECT id, login, name, surname FROM user WHERE login = ? AND hashed_password = ?", login, hashedPassword)
	var u models.User
	err := row.Scan(&u.Id, &u.Login, &u.Name, &u.Surname)
	if err == sql.ErrNoRows {
		// No user found, indicating login failure
		return models.User{}, fmt.Errorf("user not found")
	} else if err != nil {
		// Other error occurred
		return models.User{}, fmt.Errorf("failed to query data: %w", err)
	}

	// User found, indicating successful login
	r.LoggedIn = true
	return u, nil
}

func (r *Repository) Logout() {
	r.LoggedIn = false
}

func (r *Repository) HandleLogin(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	login := req.FormValue("login-inp")
	password := req.FormValue("password-inp")

	fmt.Println("Login:", login)
	fmt.Println("Password:", password)

	if login == "" || password == "" {
		http.Error(w, "Необходимо указать логин и пароль!", http.StatusBadRequest)
		return
	}

	hash := md5.Sum([]byte(password))
	hashedPass := hex.EncodeToString(hash[:])
	fmt.Println("Password hash:", hashedPass)

	user, err := r.repo.Login(login, hashedPass)
	if err != nil {
		log.Printf("Error during login for user %s: %v", login, err)
		http.Error(w, "Неверные учетные данные. Пожалуйста, повторите попытку позже.", http.StatusUnauthorized)
		return
	}

	if user.Id == 0 {
		r.LoggedIn = false
		http.Error(w, "Вы ввели неверный логин или пароль!", http.StatusUnauthorized)
		return
	}

	r.LoggedIn = true
	fmt.Println("login OK!")
	AuthenticateUser(w, req, &user)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (r *Repository) GetDB() *sql.DB {
	return r.db
}
