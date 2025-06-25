package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "time"

    bookhandlers "book-store/internal/book-store/handlers"
    bookservices "book-store/internal/book-store/handlers/services"
    bookrepos "book-store/internal/book-store/repositories"
    bookrouters "book-store/internal/book-store/routers"

    authrouters "book-store/internal/auth/routers"
    authconfig "book-store/internal/auth/config" 

    "github.com/joho/godotenv"
    "github.com/jackc/pgx/v5/pgxpool"
)

func main() {
    defer func() {
        if r := recover(); r != nil {
            log.Fatalf("panic: %v", r)
        }
    }()

    // Загрузка .env
    if err := godotenv.Load(); err != nil {
        log.Println(".env файл не найден")
    }

    // Подключение к БД
    dbURL := os.Getenv("DATABASE_URL")
    pool, err := pgxpool.New(context.Background(), dbURL)
    if err != nil {
        log.Fatalf("Ошибка подключения к БД: %v", err)
    }
    defer pool.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    if err := pool.Ping(ctx); err != nil {
        log.Fatalf("PostgreSQL не отвечает: %v", err)
    }
    log.Println("Успешное подключение к PostgreSQL")

    authconfig.DB = pool

    authconfig.LoadTokenConfig()

    bookRepo := bookrepos.NewBookRepositoryPGX(pool)
    bookService := bookservices.NewBookService(bookRepo)
    bookHandler := bookhandlers.NewBookHandler(bookService)
    bookRouter := bookrouters.SetupBookRoutes(bookHandler)

    authRouter := authrouters.SetupRoutes()

    mux := http.NewServeMux()
    mux.Handle("/auth/", http.StripPrefix("/auth", authRouter))
    mux.Handle("/books/", http.StripPrefix("/books", bookRouter))

    mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("pong"))
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Сервер запущен на http://localhost:%s", port)
    log.Fatal(http.ListenAndServe(":"+port, mux))
}