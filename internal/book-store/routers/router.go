package bookrouters

import (
    "net/http"

    "github.com/gorilla/mux"
    "book-store/internal/book-store/handlers"
    authMiddleware "book-store/internal/auth/middleware"
)

func SetupBookRoutes(bookHandler *handlers.BookHandler) http.Handler {
    r := mux.NewRouter()

    // Авторизация — применяется ко всем маршрутам
    r.Use(authMiddleware.AuthMiddleware)

    // Публичные маршруты
r.HandleFunc("/", bookHandler.GetBooks).Methods("GET")

// Админские маршруты (с проверкой роли)
adminRouter := r.NewRoute().Subrouter()
adminRouter.Use(authMiddleware.RoleMiddleware("admin"))

adminRouter.HandleFunc("/", bookHandler.CreateBook).Methods("POST")
adminRouter.HandleFunc("/{id}", bookHandler.UpdateBook).Methods("PUT")
adminRouter.HandleFunc("/{id}", bookHandler.DeleteBook).Methods("DELETE")  

    return r
}
