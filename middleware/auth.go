package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"github.com/golang-jwt/jwt/v5"
)

type key string

const userKey key = "userID"

// AuthMiddleware проверяет JWT и передаёт user_id в контекст
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Получаем заголовок Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Отсутствует токен", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 2. Проверяем токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверяем алгоритм подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Невалидный токен", http.StatusUnauthorized)
			return
		}

		// 3. Извлекаем user_id
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Ошибка в токене", http.StatusUnauthorized)
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "user_id не найден", http.StatusUnauthorized)
			return
		}

		// 4. Передаём user_id в контекст запроса
		ctx := context.WithValue(r.Context(), userKey, int(userIDFloat))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


func UserKey() interface{} {
	return userKey
}
