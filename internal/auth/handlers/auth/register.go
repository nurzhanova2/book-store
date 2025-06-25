package handlers

import (
	"book-store/internal/auth/utils"
	"book-store/internal/auth/config"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)


func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на регистрацию")

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	var input RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
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

	// ВАЖНО: проверка, что DB инициализирован
	if config.DB == nil {
		log.Println("Ошибка: DB pool не инициализирован")
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
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
