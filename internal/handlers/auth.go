package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"go-auth-app/internal/config"
	"go-auth-app/internal/models"
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

	accessToken, err := utils.GenerateJWT(userID)
	if err != nil {
		http.Error(w, "Ошибка генерации access токена", http.StatusInternalServerError)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		http.Error(w, "Ошибка генерации refresh токена", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(time.Duration(config.RefreshTokenDays) * 24 * time.Hour)

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  expiresAt,
	})

	err = models.SaveRefreshToken(r.Context(), userID, refreshToken, expiresAt)
	if err != nil {
		http.Error(w, "Ошибка сохранения refresh токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Refresh token отсутствует", http.StatusUnauthorized)
		return
	}

	storedToken, err := models.GetRefreshToken(r.Context(), cookie.Value)
	if err != nil {
		http.Error(w, "Недействительный refresh токен", http.StatusUnauthorized)
		return
	}

	if storedToken.Revoked || storedToken.ExpiresAt.Before(time.Now()) {
		http.Error(w, "Refresh токен истёк или отозван", http.StatusUnauthorized)
		return
	}

	accessToken, err := utils.GenerateJWT(storedToken.UserID)
	if err != nil {
		http.Error(w, "Ошибка генерации access токена", http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		http.Error(w, "Ошибка генерации нового refresh токена", http.StatusInternalServerError)
		return
	}

	_ = models.RevokeRefreshToken(r.Context(), storedToken.Token)

	expiresAt := time.Now().Add(time.Duration(config.RefreshTokenDays) * 24 * time.Hour)

	err = models.SaveRefreshToken(r.Context(), storedToken.UserID, newRefreshToken, expiresAt)
	if err != nil {
		http.Error(w, "Ошибка сохранения нового refresh токена", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  expiresAt,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		_ = models.RevokeRefreshToken(r.Context(), cookie.Value)

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
		})
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Вы вышли из системы")
}