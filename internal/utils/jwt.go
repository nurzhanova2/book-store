package utils

import (
	"errors"
	"os"
	"time"
	"crypto/rand"
    "encoding/hex"
	"go-auth-app/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// Claims содержит только user_id + срок действия токена
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// Создание JWT токена
func GenerateJWT(userID int) (string, error) {
    expiresIn := time.Duration(config.AccessTokenMinutes) * time.Minute

    claims := Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return "", errors.New("JWT_SECRET не задан в .env")
    }

    return token.SignedString([]byte(secret))
}

// Проверка и извлечение user_id из токена
func ParseJWT(tokenStr string) (int, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return 0, errors.New("JWT_SECRET не задан в .env")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return 0, errors.New("невалидный токен или claims")
	}

	return claims.UserID, nil
}

func GenerateRefreshToken() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}