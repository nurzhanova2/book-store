package services

import "errors"

var (
	ErrInvalidID    = errors.New("invalid book ID")
	ErrInvalidInput = errors.New("missing required fields")
)
