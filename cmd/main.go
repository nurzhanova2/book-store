package main

import (
    "log"
    "net/http"
    _ "go-auth-app/cmd/docs"
    "go-auth-app/internal/config"
    "go-auth-app/internal/routers"
    "github.com/joho/godotenv"
    _ "go-auth-app/internal/handlers/auth"
    _ "go-auth-app/internal/handlers/admin"
    _ "go-auth-app/internal/handlers/users"
)

// @title Go Auth API
// @version 1.0
// @description Это документация API для админ-панели на Go
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization


func main() {
    if err := godotenv.Load(); err != nil {
	log.Println(".env файл не найден, переменные окружения должны быть заданы вручную")
}    
    if err := config.ConnectDB(); err != nil {
    log.Fatalf("Ошибка подключения к базе данных: %v", err)
    }

    config.LoadTokenConfig()

    r := routers.SetupRoutes()

    log.Println("Сервер запущен на http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}