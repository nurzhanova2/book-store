package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"go-auth-app/internal/config"
	"go-auth-app/internal/models"
	"go-auth-app/internal/utils"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	var userID int
	var storedHash string

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

	if err := models.SaveRefreshToken(r.Context(), userID, refreshToken, expiresAt); err != nil {
		http.Error(w, "Ошибка сохранения refresh токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
}