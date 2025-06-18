package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword хэширует пароль перед сохранением
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // 14 — уровень сложности (чем выше — тем медленнее, но безопаснее)
	return string(bytes), err
}
