package handlers

import (
    "book-store/internal/book-store/handlers/dto"
    "book-store/internal/book-store/models"
    "book-store/internal/book-store/handlers/services"
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
)

type BookHandler struct {
    service *services.BookService
}

func NewBookHandler(service *services.BookService) *BookHandler {
    return &BookHandler{service: service}
}

func (h *BookHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
    idStr := mux.Vars(r)["id"]
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid book ID", http.StatusBadRequest)
        return
    }

    book, err := h.service.GetBookByID(r.Context(), id)
    if err != nil {
        http.Error(w, "Book not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
    var req dto.CreateBookRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    book := &models.Book{
        Title:    req.Title,
        Author:   req.Author,
        Genre:    req.Genre,
        Year:     req.Year,
        Quantity: req.Quantity,
    }

    id, err := h.service.CreateBook(r.Context(), book)
    if err != nil {
        http.Error(w, "Failed to create book", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "id": id,
    })
}

// Обновить книгу по ID (PUT /books/{id})
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
    idStr := mux.Vars(r)["id"]
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Неверный ID книги", http.StatusBadRequest)
        return
    }

    var req dto.UpdateBookRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Неверное тело запроса", http.StatusBadRequest)
        return
    }

    if req.Title == nil || req.Author == nil || req.Genre == nil || req.Year == nil || req.Quantity == nil {
        http.Error(w, "Отсутствуют обязательные поля", http.StatusBadRequest)
        return
    }

    book := &models.Book{
        ID:       id,
        Title:    *req.Title,
        Author:   *req.Author,
        Genre:    *req.Genre,
        Year:     *req.Year,
        Quantity: *req.Quantity,
    }

    err = h.service.UpdateBook(r.Context(), book)
    if err != nil {
        http.Error(w, "Не удалось обновить книгу", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// Удалить книгу по ID (DELETE /books/{id})
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
    idStr := mux.Vars(r)["id"]
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Неверный ID книги", http.StatusBadRequest)
        return
    }

    err = h.service.DeleteBook(r.Context(), id)
    if err != nil {
        http.Error(w, "Не удалось удалить книгу", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// Получить список книг с фильтрацией, пагинацией (GET /books)
func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filter := dto.BookFilterRequest{
		Title:  q.Get("title"),
		Author: q.Get("author"),
		Genre:  q.Get("genre"),
	}

	if yearStr := q.Get("year"); yearStr != "" {
		if year, err := strconv.Atoi(yearStr); err == nil {
			filter.Year = &year
		}
	}
	if limitStr := q.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	} else {
		filter.Limit = 10
	}
	if offsetStr := q.Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filter.Offset = offset
		}
	}

	books, err := h.service.GetBooks(r.Context(), filter)
	if err != nil {
		http.Error(w, "Не удалось получить книги", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}