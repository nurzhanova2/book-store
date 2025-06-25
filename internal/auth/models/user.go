package models

import (
    "context"
    "database/sql"
    "time"

    "book-store/internal/auth/config"
    "book-store/internal/auth/utils"

    "github.com/google/uuid"
)

type User struct {
    ID        uuid.UUID `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Role      string    `json:"role"`
    Password  string    `json:"password"`
    CreatedAt time.Time `json:"created_at"`
}

// Проверка, существует ли email
func EmailExists(db *sql.DB, email string) (bool, error) {
    var exists bool
    query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
    err := db.QueryRowContext(context.Background(), query, email).Scan(&exists)
    return exists, err
}

// Получить роль пользователя по UUID
func GetUserRole(ctx context.Context, userID uuid.UUID) (string, error) {
    var role string
    err := config.DB.QueryRow(ctx, "SELECT role FROM users WHERE id = $1", userID).Scan(&role)
    return role, err
}

// Получить всех пользователей
func GetAllUsers() ([]User, error) {
    query := "SELECT id, username, email, role, created_at FROM users"
    rows, err := config.DB.Query(context.Background(), query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var u User
        if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.CreatedAt); err != nil {
            return nil, err
        }
        users = append(users, u)
    }

    return users, nil
}

// Создать пользователя
func CreateUser(user User) error {
    hashedPassword, err := utils.HashPassword(user.Password)
    if err != nil {
        return err
    }

    query := `INSERT INTO users (username, email, password, role, created_at)
              VALUES ($1, $2, $3, $4, $5)`
    _, err = config.DB.Exec(
        context.Background(),
        query,
        user.Username,
        user.Email,
        hashedPassword,
        user.Role,
        time.Now(),
    )
    return err
}

// Обновить пользователя по UUID
func UpdateUserByID(id uuid.UUID, updatedUser User) error {
    query := `UPDATE users
              SET username = $1, email = $2, role = $3, updated_at = NOW()
              WHERE id = $4`
    _, err := config.DB.Exec(
        context.Background(),
        query,
        updatedUser.Username,
        updatedUser.Email,
        updatedUser.Role,
        id,
    )
    return err
}

// Удалить пользователя по UUID
func DeleteUserByID(id uuid.UUID) error {
    query := `DELETE FROM users WHERE id = $1`
    _, err := config.DB.Exec(context.Background(), query, id)
    return err
}