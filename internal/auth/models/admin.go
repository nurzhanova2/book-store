package models

import (
    "context"
    "time"

    "book-store/internal/auth/config"

    "github.com/google/uuid"
)

type DashboardData struct {
    TotalUsers     int           `json:"total_users"`
    ActiveSessions int           `json:"active_sessions"`
    LastLogins     []LastLogin   `json:"last_logins"`
}

type LastLogin struct {
    UserID    uuid.UUID `json:"user_id"`   
    Email     string    `json:"email"`
    LastLogin time.Time `json:"last_login"`
}

func GetDashboardData() (*DashboardData, error) {
    ctx := context.Background()
    var data DashboardData

    // Общее количество пользователей
    err := config.DB.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&data.TotalUsers)
    if err != nil {
        return nil, err
    }

    // Активные сессии
    err = config.DB.QueryRow(ctx, `
        SELECT COUNT(*) FROM refresh_tokens
        WHERE revoked = FALSE AND expires_at > NOW()
    `).Scan(&data.ActiveSessions)
    if err != nil {
        return nil, err
    }

    // Последние входы
    rows, err := config.DB.Query(ctx, `
        SELECT id, email, last_login FROM users
        WHERE last_login IS NOT NULL
        ORDER BY last_login DESC
        LIMIT 5
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var l LastLogin
        if err := rows.Scan(&l.UserID, &l.Email, &l.LastLogin); err != nil {
            return nil, err
        }
        data.LastLogins = append(data.LastLogins, l)
    }

    return &data, nil
}