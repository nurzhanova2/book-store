package handlers

import (
	"encoding/json"
	"net/http"
	"go-auth-app/middleware" // импортируем middleware, чтобы получить userKey
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserKey()).(int)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Добро пожаловать!",
		"user_id": userID,
	})
}
