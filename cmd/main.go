package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
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

    r := mux.NewRouter()

    // Auth endpoints
    r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
    r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
    r.HandleFunc("/auth/logout", handlers.LogoutHandler).Methods("POST")
    r.HandleFunc("/auth/refresh", handlers.RefreshHandler).Methods("POST")

    // Protected route
    r.Handle("/profile", middleware.AuthMiddleware(http.HandlerFunc(handlers.ProfileHandler)))

    // Admin routes
    admin := r.PathPrefix("/admin").Subrouter()
    admin.Use(middleware.RoleMiddleware("admin"))

    admin.HandleFunc("", handlers.AdminHandler).Methods("GET")
    admin.HandleFunc("/users", handlers.GetAllUsers).Methods("GET")
    admin.HandleFunc("/users", handlers.CreateUser).Methods("POST")
    admin.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")
    admin.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")

    log.Println("Сервер запущен на http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
