package repository

import (
	"time"

	"book-store/pkg/paginator"
)

type RegisterOptions struct {
	Email     string
	Phone     string
	Password  string
	FullName  string
	Verified  bool
	GroupID   string
	GroupRole string
	GroupName string
	RefererID string
}

type DetailOptions struct {
	Email  string
	UserID string
}

type TokenOptions struct {
	UserID    string
	GroupID   string
	GroupRole string
}

type UpdateOptions struct {
	ID          string
	Email       string
	Phone       string
	FullName    string
	Password    string
	LastOrderAt *time.Time
	Priority    bool
}

type CreateActivityOptions struct {
	UserID string
	Type   string
}

type ListFilterOptions struct {
	IDs  []string
	Name string
}
type ListOptions struct {
	Filter         ListFilterOptions
	PaginatorQuery paginator.PaginatorQuery
}

type UpdateAdminOptions struct {
	Email     string
	Phone     string
	FullName  string
	GroupID   string
	GroupRole string
	GroupName string
	Discount  *map[int]float64
	Priority  *bool
}
