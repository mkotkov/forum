package middleware

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"forum/models"
	"log"
	"net/http"
	"golang.org/x/crypto/bcrypt"
)

type RepositoryInterface interface {
	GetDB() *sql.DB
	Login(login, password string) (models.User, error)
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

func (r *Repository) Login(login, password string) (models.User, error) {
	if r == nil || r.repo == nil {
		return models.User{}, fmt.Errorf("repository is nil")
	}

	row := r.db.QueryRowContext(r.ctx, "SELECT id, login, name, surname, hashed_password FROM user WHERE login = ?", login)
	var u models.User
	err := row.Scan(&u.Id, &u.Login, &u.Name, &u.Surname, &u.HashedPassword)
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
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Println("Error hashing password:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }


	fmt.Println("Login:", login)
	fmt.Println("Password:", password)
	fmt.Println("Bcrypt-хеш пароля:", string(hashedPass))

	if login == "" || password == "" {
		http.Error(w, "Необходимо указать логин и пароль!", http.StatusBadRequest)
		return
	}

	user, err := r.repo.getUserByLogin(login)
	if err != nil {
		log.Println("Error getting user by login:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Stored Hashed Password:", user.HashedPassword)

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
        log.Println("Incorrect login or password")
        http.Error(w, "Invalid login or password", http.StatusUnauthorized)
        return
    }

	r.LoggedIn = true
	fmt.Println("login OK!")
	AuthenticateUser(w, req, user)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}



func (r *Repository) GetDB() *sql.DB {
	return r.db
}

func (repo *Repository) getUserByLogin(login string) (*models.User, error) {
	// Используйте ваш SQL-запрос для поиска пользователя по логину
	row := repo.db.QueryRow("SELECT id, login, hashed_password FROM user WHERE login = ?", login)

	var user models.User
	err := row.Scan(&user.Id, &user.Login, &user.HashedPassword)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
