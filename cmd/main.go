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
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	if err := config.ConnectDB(); err != nil {
		log.Fatal(err)
	}

	config.LoadTokenConfig()

	// Auth endpoints
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/auth/logout", handlers.LogoutHandler)
	http.HandleFunc("/auth/refresh", handlers.RefreshHandler)

	// Protected endpoints
	http.Handle("/profile", middleware.AuthMiddleware(http.HandlerFunc(handlers.ProfileHandler)))
	http.Handle("/admin", middleware.RoleMiddleware("admin")(http.HandlerFunc(handlers.AdminHandler)))

	log.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
