package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go-auth-app/config"
	"go-auth-app/utils"

	"github.com/jackc/pgconn" 
	"golang.org/x/crypto/bcrypt"

)

type RegisterInput struct {
	Username string `json:"username"` 
	Email    string `json:"email"`  
	Password string `json:"password"` 
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
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

	// 1. Проверяем, что email имеет правильный формат
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		http.Error(w, "Ошибка хэширования", http.StatusInternalServerError)
		return
	}

    // 2. Проверяем, что email и username уникальны
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`

	_, err = config.DB.Exec(r.Context(), query, input.Username, input.Email, hashedPassword)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			http.Error(w, "Пользователь с таким email или username уже существует", http.StatusConflict)
			return
		}
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Регистрация прошла успешно")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input LoginInput

	// 1. Читаем тело запроса
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	input.Email = strings.TrimSpace(input.Email)
	input.Password = strings.TrimSpace(input.Password)

	if input.Email == "" || input.Password == "" {
		http.Error(w, "Email и пароль обязательны", http.StatusBadRequest)
		return
	}

	// 2. Ищем пользователя в базе
	var userID int
	var hashedPassword string

	query := `SELECT id, password FROM users WHERE email = $1`

	err = config.DB.QueryRow(r.Context(), query, input.Email).Scan(&userID, &hashedPassword)
	if err != nil {
		http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}

	// 3. Сравниваем пароли
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password))
	if err != nil {
		http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}

	// 4. Генерируем токен
	token, err := utils.GenerateJWT(userID)
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	// 5. Отправляем токен клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
