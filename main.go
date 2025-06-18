package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"go-auth-app/config"
	"go-auth-app/handlers"
	"go-auth-app/middleware"
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
	http.Handle("/profile", middleware.AuthMiddleware(http.HandlerFunc(handlers.ProfileHandler)))

	log.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
