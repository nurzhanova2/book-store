package models

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` 
	CreatedAt time.Time `json:"created_at"`
}

func EmailExists(db *sql.DB, email string) (bool, error) {
	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err := db.QueryRowContext(context.Background(), query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
