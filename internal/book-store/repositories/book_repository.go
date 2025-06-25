package repositories

import (
	"context"

	"book-store/internal/book-store/models"
)

type BookRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Book, error)
	GetList(ctx context.Context, filter BookFilter) ([]*models.Book, error)
	Create(ctx context.Context, book *models.Book) (int64, error)
	Update(ctx context.Context, book *models.Book) error
	Delete(ctx context.Context, id int64) error
}


type BookFilter struct {
	Title    string
	Author   string
	Genre    string
	Year     *int
	Limit    int
	Offset   int
}
