package repositories

import (
	"context"
	"book-store/internal/book-store/models"

    "github.com/jackc/pgx/v5/pgxpool"

)

type BookRepositoryPGX struct {
    db *pgxpool.Pool
}

func NewBookRepositoryPGX(db *pgxpool.Pool) *BookRepositoryPGX {
    return &BookRepositoryPGX{db: db}
}

func (r *BookRepositoryPGX) GetByID(ctx context.Context, id int64) (*models.Book, error) {
	query := `SELECT id, title, author, genre, year, quantity, created_at, updated_at FROM books WHERE id = $1`
	var book models.Book
	err := r.db.QueryRow(ctx, query, id).Scan(
		&book.ID, &book.Title, &book.Author, &book.Genre, &book.Year,
		&book.Quantity, &book.CreatedAt, &book.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// Создание книги
func (r *BookRepositoryPGX) Create(ctx context.Context, book *models.Book) (int64, error) {
	query := `
        INSERT INTO books (title, author, genre, year, quantity)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`
	var id int64
	err := r.db.QueryRow(ctx, query,
		book.Title, book.Author, book.Genre, book.Year, book.Quantity,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Обновление книги
func (r *BookRepositoryPGX) Update(ctx context.Context, book *models.Book) error {
	query := `
        UPDATE books
        SET title = $1, author = $2, genre = $3, year = $4, quantity = $5, updated_at = CURRENT_TIMESTAMP
        WHERE id = $6`
	_, err := r.db.Exec(ctx, query,
		book.Title, book.Author, book.Genre, book.Year, book.Quantity, book.ID,
	)
	return err
}

// Удаление книги
func (r *BookRepositoryPGX) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM books WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *BookRepositoryPGX) GetList(ctx context.Context, filter BookFilter) ([]*models.Book, error) {
	query := `
		SELECT id, title, author, genre, year, quantity, created_at, updated_at
		FROM books
		ORDER BY id
		LIMIT $1 OFFSET $2`
	
	rows, err := r.db.Query(ctx, query, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*models.Book

	for rows.Next() {
		var book models.Book
		err := rows.Scan(
			&book.ID, &book.Title, &book.Author, &book.Genre, &book.Year,
			&book.Quantity, &book.CreatedAt, &book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}