package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"go-auth-app/internal/config"
	"go-auth-app/internal/utils"

	"github.com/jackc/pgconn"
)

type RegisterInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на регистрацию")

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	var input RegisterInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	input.Username = strings.TrimSpace(input.Username)
	input.Email = strings.TrimSpace(input.Email)
	input.Password = strings.TrimSpace(input.Password)

	if input.Username == "" || input.Email == "" || input.Password == "" {
		http.Error(w, "Все поля обязательны", http.StatusBadRequest)
		return
	}

	if _, err := mail.ParseAddress(input.Email); err != nil {
		http.Error(w, "Невалидный email", http.StatusBadRequest)
		return
	}

	if err := utils.ValidatePassword(input.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		http.Error(w, "Ошибка хэширования пароля", http.StatusInternalServerError)
		return
	}

	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`
	_, err = config.DB.Exec(r.Context(), query, input.Username, input.Email, hashedPassword)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			http.Error(w, "Пользователь уже существует", http.StatusConflict)
			return
		}
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Регистрация прошла успешно")
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	var storedHash string
	var userID int

	query := `SELECT id, password FROM users WHERE email = $1`
	err := config.DB.QueryRow(r.Context(), query, input.Email).Scan(&userID, &storedHash)
	if err != nil || !utils.CheckPasswordHash(input.Password, storedHash) {
		http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(userID)
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}