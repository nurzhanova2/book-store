package handlers

import (
	"fmt"
	"net/http"
	"time"

	"go-auth-app/internal/models"
)

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