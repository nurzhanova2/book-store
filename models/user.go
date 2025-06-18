package models

import "time"

// User представляет модель пользователя в базе данных
// с полями ID, Username, Email, Password и CreatedAt.
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` 
	// Используем тег "-" для того, чтобы не выводить пароль в JSON ответах
	CreatedAt time.Time `json:"created_at"`
}
