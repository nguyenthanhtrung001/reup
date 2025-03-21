package repository

import (
	"context"

	"book-store/internal/models"
	"book-store/pkg/paginator"
)

type Repository interface {
	UserRepository
	GroupRepository
}

type UserRepository interface {
	CreateNewUser(ctx context.Context, ro RegisterOptions) error
	GetUser(ctx context.Context, do DetailOptions) (models.User, error)
	GenerateToken(ctx context.Context, to TokenOptions) (string, error)
	UpdateUser(ctx context.Context, uo UpdateOptions) error
	IsExist(ctx context.Context, email string) bool
	CreateActivity(ctx context.Context, opt CreateActivityOptions) error

	// Admin User
	List(ctx context.Context, sc models.Scope, opt ListOptions) ([]models.User, paginator.Paginator, error)
	ListByIDS(ctx context.Context, sc models.Scope, user_ids []string) ([]models.User, error)
	UpdateByAdmin(ctx context.Context, sc models.Scope, opt UpdateAdminOptions, id string) (models.User, error)
	UpdateCreditUser(ctx context.Context, sc models.Scope, credit float64, id string) error
	UpdateCreditUserByOrder(ctx context.Context, sc models.Scope, credit float64, id string) error
	UpdateApiKeyUser(ctx context.Context, sc models.Scope, id string, apiKey string) error
}
type GroupRepository interface {
	GetGroup(ctx context.Context, sc models.Scope, id string) (models.Group, error)
}
