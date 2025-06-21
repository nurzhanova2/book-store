package config

import (
	"log"
	"os"
	"strconv"
)

var AccessTokenMinutes int
var RefreshTokenDays int

func LoadTokenConfig() {
	var err error

	AccessTokenMinutes, err = strconv.Atoi(os.Getenv("ACCESS_TOKEN_LIFETIME_MINUTES"))
	if err != nil {
		log.Println("ACCESS_TOKEN_LIFETIME_MINUTES не задан, используется значение по умолчанию: 15")
		AccessTokenMinutes = 15
	}

	RefreshTokenDays, err = strconv.Atoi(os.Getenv("REFRESH_TOKEN_LIFETIME_DAYS"))
	if err != nil {
		log.Println("REFRESH_TOKEN_LIFETIME_DAYS не задан, используется значение по умолчанию: 7")
		RefreshTokenDays = 7
	}
}