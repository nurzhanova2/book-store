package middleware

import (
    "net/http"
    "strings"
    "book-store/internal/auth/models"
    "book-store/internal/auth/utils"
    "github.com/google/uuid"
)

func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if !strings.HasPrefix(authHeader, "Bearer ") {
                http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
                return
            }

            tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
            userIDStr, err := utils.ParseJWT(tokenStr)
            if err != nil {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }

            userID, err := uuid.Parse(userIDStr) // ✅ вот здесь происходит нужный каст
            if err != nil {
                http.Error(w, "Invalid user ID format", http.StatusBadRequest)
                return
            }

            role, err := models.GetUserRole(r.Context(), userID)
            if err != nil {
                http.Error(w, "Ошибка получения роли", http.StatusInternalServerError)
                return
            }

            for _, allowed := range allowedRoles {
                if role == allowed {
                    next.ServeHTTP(w, r)
                    return
                }
            }

            http.Error(w, "Forbidden", http.StatusForbidden)
        })
    }
}