package utils

import (
    "crypto/rand"
    "encoding/hex"
    "errors"
    "fmt"
    "os"
    "time"

    "book-store/internal/auth/config"

    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID string `json:"user_id"`
    jwt.RegisteredClaims
}

// Генерация JWT access-токена
func GenerateJWT(userID string) (string, error) {
    // Получаем длительность из конфига (в минутах)
    expiresIn := time.Duration(config.AccessTokenMinutes) * time.Minute

    // Используем UTC-время, вычитаем 1 сек. — защита от microdrift
    issuedAt := time.Now().UTC().Add(-1 * time.Second)
    expiresAt := issuedAt.Add(expiresIn)


    claims := Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            IssuedAt:  jwt.NewNumericDate(issuedAt),
            ExpiresAt: jwt.NewNumericDate(expiresAt),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return "", errors.New("JWT_SECRET не задан в .env")
    }

    return token.SignedString([]byte(secret))
}


// Проверка и извлечение userID из JWT
func ParseJWT(tokenStr string) (string, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        fmt.Println("❌ JWT_SECRET не задан")
        return "", errors.New("JWT_SECRET не задан в .env")
    }

    fmt.Println("🔐 JWT_SECRET (из env):", secret[:5], "...") // обрезаем для безопасности

    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })
    if err != nil {
        fmt.Println("❌ Ошибка парсинга токена:", err)
        return "", err
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        fmt.Println("❌ Невалидный токен. Claims:", claims, "Valid:", token.Valid)
        return "", errors.New("невалидный токен или claims")
    }

    fmt.Println("✅ JWT валиден. UserID:", claims.UserID)
    return claims.UserID, nil
}

// Генерация криптостойкого refresh токена
func GenerateRefreshToken() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}
