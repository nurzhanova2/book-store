package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Создаём JWT-токен
func GenerateJWT(userID int) (string, error) {
	// Создаём claims — полезную информацию внутри токена
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // срок действия — 24 часа
	}

	// Создаём токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Секрет берём из переменных окружения
	secret := os.Getenv("JWT_SECRET")

	// Подписываем токен
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
