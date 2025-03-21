package repository

import (
	"book-store/pkg/paginator"
	"time"
)

type CreateOption struct {
	Title       string
	Author      string
	PublishedAt string
}

type UpdateOpition struct {
	ID          string
	Title       string
	Author      string
	PublishedAt string
}

type DetailOutputOption struct {
	ID          string    `bson:"_id"`
	Title       string    `bson:"title"`
	Author      string    `bson:"author"`
	PublishedAt string    `bson:"published_at"`
	CreateAt    time.Time `bson:"create_at"`
	UpdateAt    time.Time `bson:"update_at"`
}
type ListFilterOptions struct {
	Title  string
	Author string
}

type ListOptions struct {
	Filter         ListFilterOptions
	PaginatorQuery paginator.PaginatorQuery
}
