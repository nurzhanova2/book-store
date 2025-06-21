package models

import (
	"context"
	"time"

	"go-auth-app/internal/config"
)

type RefreshToken struct {
	Token     string
	UserID    int
	ExpiresAt time.Time
	Revoked   bool
}

func GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	query := `
        SELECT token, user_id, expires_at, revoked
        FROM refresh_tokens
        WHERE token = $1
    `

	var rt RefreshToken
	err := config.DB.QueryRow(ctx, query, token).Scan(&rt.Token, &rt.UserID, &rt.ExpiresAt, &rt.Revoked)
	if err != nil {
		return nil, err
	}

	return &rt, nil
}

func SaveRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	query := `
        INSERT INTO refresh_tokens (token, user_id, expires_at)
        VALUES ($1, $2, $3)
    `
	_, err := config.DB.Exec(ctx, query, token, userID, expiresAt)
	return err
}

func RevokeRefreshToken(ctx context.Context, token string) error {
	query := `
        UPDATE refresh_tokens
        SET revoked = TRUE
        WHERE token = $1
    `
	_, err := config.DB.Exec(ctx, query, token)
	return err
}