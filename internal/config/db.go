package config

import (
    "context"
    "fmt"
    "os"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() error {
    dbURL := os.Getenv("DATABASE_URL")

    if dbURL == "" {
        return fmt.Errorf("переменная окружения DATABASE_URL не установлена")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    pool, err := pgxpool.New(ctx, dbURL)
    if err != nil {
        return fmt.Errorf("ошибка при создании пула подключений: %v", err)
    }

    err = pool.Ping(ctx)
    if err != nil {
        return fmt.Errorf("не удалось подключиться к базе: %v", err)
    }

    DB = pool
    fmt.Println("Успешное подключение к базе данных")
    return nil
}
