package usecase

import (
	"context"

	"reup/internal/models"
	"reup/internal/user/repository"
	"reup/pkg/encrypter"
	"reup/pkg/log"
)

type UseCase interface {
	UserUseCase
}

type UserUseCase interface {
	Register(ctx context.Context, input RegisterInput) error
	Login(ctx context.Context, input LoginInput) (LoginOutput, error)
	Profile(ctx context.Context, sc models.Scope) (ProfileOutput, error)
	UpdateProfile(ctx context.Context, sc models.Scope, input UpdateProfileInput) error
	UpdatePassword(ctx context.Context, sc models.Scope, input UpdatePasswordInput) error
	GetUserById(ctx context.Context, sc models.Scope, id string) (models.User, error)
	CreateActivity(ctx context.Context, sc models.Scope, input CreateActivityInput) error
	UpdateUserByID(ctx context.Context, sc models.Scope, input UpdateUserByIDInput, id string) (models.User, error)

	// Admin User
	List(ctx context.Context, sc models.Scope, input ListInput) (ListOutput, error)
	ListByIDS(ctx context.Context, sc models.Scope, user_ids []string) ([]models.User, error)
	UpdateByAdmin(ctx context.Context, sc models.Scope, input UpdateByAdminInput, id string) (models.User, error)
	CreateApiKeyUser(ctx context.Context, sc models.Scope, id string) (models.User, error)
}

type implUseCase struct {
	l         log.Logger
	repo      repository.Repository
	encrypter encrypter.Encrypter
}

func New(l log.Logger, repo repository.Repository,
	encrypter encrypter.Encrypter,

) UseCase {
	return implUseCase{
		l:         l,
		repo:      repo,
		encrypter: encrypter,
	}
}
