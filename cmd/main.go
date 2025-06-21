package main

import (
    "log"
    "net/http"
    "go-auth-app/internal/config"
    "go-auth-app/internal/routers"
    "github.com/joho/godotenv"
)

func main() {
    _ = godotenv.Load()
    _ = config.ConnectDB()
    config.LoadTokenConfig()

    r := routers.SetupRoutes()

    log.Println("Сервер запущен на http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}