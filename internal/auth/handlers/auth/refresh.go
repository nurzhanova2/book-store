package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"book-store/internal/auth/config"
	"book-store/internal/auth/models"
	"book-store/internal/auth/utils"
)

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Refresh token отсутствует", http.StatusUnauthorized)
		return
	}

	storedToken, err := models.GetRefreshToken(r.Context(), cookie.Value)
	if err != nil || storedToken.Revoked || storedToken.ExpiresAt.Before(time.Now()) {
		http.Error(w, "Недействительный или истёкший refresh токен", http.StatusUnauthorized)
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