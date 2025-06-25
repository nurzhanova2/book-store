package services

import (
    "book-store/internal/book-store/models"
    "book-store/internal/book-store/repositories"
    "book-store/internal/book-store/handlers/dto"
    "context"
)

type BookService struct {
    repo repositories.BookRepository
}

func NewBookService(repo repositories.BookRepository) *BookService {
    return &BookService{repo: repo}
}

func (s *BookService) GetBookByID(ctx context.Context, id int64) (*models.Book, error) {
    if id <= 0 {
        return nil, ErrInvalidID
    }

    return s.repo.GetByID(ctx, id)
}

func (s *BookService) CreateBook(ctx context.Context, book *models.Book) (int64, error) {
    if book.Title == "" || book.Author == "" {
        return 0, ErrInvalidInput
    }

    return s.repo.Create(ctx, book)
}

func (s *BookService) UpdateBook(ctx context.Context, book *models.Book) error {
    if book.ID <= 0 {
        return ErrInvalidID
    }
    if book.Title == "" || book.Author == "" {
        return ErrInvalidInput
    }

    return s.repo.Update(ctx, book)
}

func (s *BookService) DeleteBook(ctx context.Context, id int64) error {
    if id <= 0 {
        return ErrInvalidID
    }

    return s.repo.Delete(ctx, id)
}

// Новый метод
func (s *BookService) GetBooks(ctx context.Context, filter dto.BookFilterRequest) ([]*models.Book, error) {
	repoFilter := repositories.BookFilter{
		Title:  filter.Title,
		Author: filter.Author,
		Genre:  filter.Genre,
		Year:   filter.Year,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	}
	return s.repo.GetList(ctx, repoFilter)
}