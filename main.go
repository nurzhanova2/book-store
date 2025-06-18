package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"go-auth-app/config"
	"go-auth-app/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	if err := config.ConnectDB(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/register", handlers.RegisterHandler)
    http.HandleFunc("/login", handlers.LoginHandler)

	log.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
