package usecase

import (
	"reup/internal/models"
	"reup/pkg/paginator"
	"time"
)

type BookInput struct {
	Title       string
	Author      string
	PublishedAt string
}

type UpdateInput struct {
	ID          string
	Title       string
	Author      string
	PublishedAt string
}
type BookOutput struct {
	ID          string
	Title       string
	Author      string
	PublishedAt string
	CreateAt    time.Time
	UpdateAt    time.Time
}
type ListFilterOptions struct {
	Title  string
	Author string
}

type ListBookInput struct {
	Filter         ListFilterOptions
	PaginatorQuery paginator.PaginatorQuery
}

type ListBookOutput struct {
	Books     []models.Book
	Paginator paginator.Paginator
}
