package usecase

import (
	"time"

	"reup/internal/models"
	"reup/pkg/paginator"
)

type RegisterInput struct {
	FullName  string
	Email     string
	Phone     string
	Password  string
	Verified  bool
	RefererID string
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	Token string
	User  ProfileOutput
}

type UpdateProfileInput struct {
	Name     string
	Priority bool
}

type UpdateUserByIDInput struct {
	Name        string
	LastOrderAt *time.Time
}

type UpdatePasswordInput struct {
	OldPassword string
	NewPassword string
}

type groupConfig struct {
	ID   string
	Name string
	Role string
}
type refererInfo struct {
	ID   string
	Name string
}
type ProfileOutput struct {
	ID         string
	Email      string
	Phone      string
	FullName   string
	Group      groupConfig
	Referer    refererInfo
	Credit     float64
	CreditUsed float64
	Discount   map[int]float64
}

type CreateActivityInput struct {
	UserID string
	Type   string
}

type ListInputFilter struct {
	IDs  []string
	Name string
}
type ListInput struct {
	Filter         ListInputFilter
	PaginatorQuery paginator.PaginatorQuery
}

type ListOutput struct {
	Users     []models.User
	Paginator paginator.Paginator
}

type UpdateByAdminInput struct {
	Email    string
	Phone    string
	FullName string
	GroupID  string
	Discount *map[int]float64
	Priority bool
}

type RefundCreditByOrderUserInput struct {
	UserID      string
	Quantity    int
	Note        string
	ServiceType string
	ServiceID   int
	OrderID     int
	Amount      float64
}
