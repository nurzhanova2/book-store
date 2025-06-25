package models

import "time"

type Book struct {
	ID          int64     `db:"id"`
	Title       string    `db:"title"`
	Author      string    `db:"author"`
	Genre       string    `db:"genre"`
	Year        int       `db:"year"`
	Quantity    int       `db:"quantity"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}