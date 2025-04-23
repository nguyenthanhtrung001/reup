package usecase

import (
	"context"

	"github.com/nguyenthanhtrung001/reup/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (uc implUseCase) getProfile(ctx context.Context, user models.User) ProfileOutput {
	profile := ProfileOutput{
		ID:       user.ID.Hex(),
		Email:    user.Email,
		Phone:    user.Phone,
		FullName: user.FullName,
		Group: groupConfig{
			ID:   user.GroupID.Hex(),
			Name: user.GroupName,
			Role: user.GroupRole,
		},
		Credit:     user.Credit,
		CreditUsed: user.CreditUsed,
		Discount:   user.Discount,
	}
	if user.RefererID.Hex() != primitive.NilObjectID.Hex() {
		refererUser, err := uc.GetUserById(ctx, models.Scope{}, user.RefererID.Hex())
		if err == nil {
			profile.Referer = refererInfo{
				ID:   refererUser.ID.Hex(),
				Name: refererUser.FullName,
			}
		}
	}
	return profile
}
