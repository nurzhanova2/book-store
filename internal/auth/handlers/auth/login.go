package handlers

import (
    "encoding/json"
    "net/http"
    "time"

    "book-store/internal/auth/config"
    "book-store/internal/auth/models"
    "book-store/internal/auth/utils"
)

// LoginHandler godoc
// @Summary Авторизация пользователя
// @Description Логин по email и паролю. Возвращает access_token.
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginInput true "Данные для логина"
// @Success 200 {object} map[string]string
// @Failure 401 {string} string "Неверный email или пароль"
// @Router /login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    var input LoginInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Невалидный JSON", http.StatusBadRequest)
        return
    }

    var userID string

    var storedHash string

    query := `SELECT id, password FROM users WHERE email = $1`
    err := config.DB.QueryRow(r.Context(), query, input.Email).Scan(&userID, &storedHash)
    
    // Отладочная информация:
    println("===[DEBUG: LoginHandler]================")
    println("Email:", input.Email)
    println("Введённый пароль:", input.Password)
    println("Хеш из БД:", storedHash)

    if err != nil {
        println("Ошибка при запросе к БД:", err.Error())
    }

    if !utils.CheckPasswordHash(input.Password, storedHash) {
        println("❌ Пароль НЕ совпал (bcrypt.CompareHashAndPassword)")
    } else {
        println("✅ Пароль совпал (bcrypt ok)")
    }
    println("========================================")

    if err != nil || !utils.CheckPasswordHash(input.Password, storedHash) {
        http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
        return
    }

    _, err = config.DB.Exec(
        r.Context(),
        `UPDATE users SET last_login = NOW() WHERE id = $1`,
        userID,
    )

    accessToken, err := utils.GenerateJWT(userID)
    if err != nil {
        http.Error(w, "Ошибка генерации access токена", http.StatusInternalServerError)
        return
    }

    refreshToken, err := utils.GenerateRefreshToken()
    if err != nil {
        http.Error(w, "Ошибка генерации refresh токена", http.StatusInternalServerError)
        return
    }

    expiresAt := time.Now().Add(time.Duration(config.RefreshTokenDays) * 24 * time.Hour)

    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    refreshToken,
        HttpOnly: true,
        Secure:   true,
        Path:     "/",
        Expires:  expiresAt,
    })

    if err := models.SaveRefreshToken(r.Context(), userID, refreshToken, expiresAt); err != nil {
        http.Error(w, "Ошибка сохранения refresh токена", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "access_token": accessToken,
    })
}