package utils

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // 14 — уровень сложности (чем выше — тем медленнее, но безопаснее)
	return string(bytes), err
}

func ValidatePassword(password string) error {
	var (
		hasMinLen  = false
		hasNumber  = false
		hasUpper   = false
		hasSpecial = false
	)

	if len(password) >= 8 {
		hasMinLen = true
	}

	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	if !hasMinLen || !hasNumber || !hasUpper || !hasSpecial {
		return errors.New("пароль должен содержать минимум 8 символов, цифру, заглавную букву и спецсимвол")
	}

	return nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

