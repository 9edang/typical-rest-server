package repository

import (
	"context"
	"time"
)

// Book represented database model
type Book struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"required"`
	Author    string    `json:"author" validate:"required"`
	UpdatedAt time.Time `json:"update_at"`
	CreatedAt time.Time `json:"created_at"`
}

// BookRepo to get book data from database [mock]
type BookRepo interface {
	FindOne(context.Context, int64) (*Book, error)
	Find(context.Context) ([]*Book, error)
	Create(context.Context, *Book) (*Book, error)
	Delete(context.Context, int64) error
	Update(context.Context, Book) error
}

// NewBookRepo return new instance of BookRepo [constructor]
func NewBookRepo(impl CachedBookRepoImpl) BookRepo {
	return &impl
}
