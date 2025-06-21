package middleware

import (
	"context"
	"net/http"
	"strings"

	"go-auth-app/internal/utils"
)

type contextKey string

const userKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Токен не передан", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := utils.ParseJWT(tokenStr)
		if err != nil {
			http.Error(w, "Невалидный токен: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserKey() interface{} {
	return userKey
}