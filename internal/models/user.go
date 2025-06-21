package models

import (
	"context"
	"database/sql"
	"time"
	"go-auth-app/internal/config"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Password  string    `json:"-"` // не возвращается в JSON
	CreatedAt time.Time `json:"created_at"`
}

// Проверка существования email
func EmailExists(db *sql.DB, email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err := db.QueryRowContext(context.Background(), query, email).Scan(&exists)
	return exists, err
}

// Получение роли пользователя
func GetUserRole(ctx context.Context, userID int) (string, error) {
	var role string
	err := config.DB.QueryRow(ctx, "SELECT role FROM users WHERE id = $1", userID).Scan(&role)
	return role, err
}