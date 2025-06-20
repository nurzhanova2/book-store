package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"go-auth-app/internal/config"
	"go-auth-app/internal/handlers"
	"go-auth-app/internal/middleware"
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
	log.Println("✅ Роут /register подключен") 
    http.HandleFunc("/login", handlers.LoginHandler)
	http.Handle("/profile", middleware.AuthMiddleware(http.HandlerFunc(handlers.ProfileHandler)))

	log.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
