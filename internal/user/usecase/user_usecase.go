package usecase

import (
	"context"

	"reup/internal/models"
	"reup/internal/user/repository"
	"reup/pkg/jwt"
)

func (uc implUseCase) Register(ctx context.Context, input RegisterInput) error {
	hashedPassword := hashPassword(input.Password)

	if uc.repo.IsExist(ctx, input.Email) {
		uc.l.Errorf(ctx, "user.usecase.Register.repo.IsExist: %v", ErrUserAlreadyExist)
		return ErrUserAlreadyExist
	}
	mUser := models.User{}
	defaultGroup := mUser.GetDefaultGroup()
	ro := repository.RegisterOptions{
		Email:     input.Email,
		Phone:     input.Phone,
		Password:  hashedPassword,
		FullName:  input.FullName,
		Verified:  true,
		GroupID:   defaultGroup.GroupID,
		GroupRole: defaultGroup.GroupRole,
		GroupName: defaultGroup.GroupName,
	}

	// Check referer_id
	if input.RefererID != "" {
		ro.RefererID = input.RefererID
	}

	err := uc.repo.CreateNewUser(ctx, ro)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.Register.repo.Register: %v", err)
		return err
	}

	return nil
}

func (uc implUseCase) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
	do := repository.DetailOptions{
		Email: input.Email,
	}

	user, err := uc.repo.GetUser(ctx, do)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.Login.repo.GetUserByEmail: %v", err)
		return LoginOutput{}, ErrUserNotFound
	}

	if user.ID.Hex() == "" {
		uc.l.Errorf(ctx, "user.usecase.Login.repo.GetUserByEmail: %v", ErrUserNotFound)
		return LoginOutput{}, ErrUserNotFound
	}

	if !checkPassword(input.Password, user.Password) {
		uc.l.Errorf(ctx, "user.usecase.Login.repo.GetUserByEmail: %v", ErrInvalidValidation)
		return LoginOutput{}, ErrInvalidValidation
	}

	to := repository.TokenOptions{
		UserID:    user.ID.Hex(),
		GroupID:   user.GroupID.Hex(),
		GroupRole: user.GroupRole,
	}

	token, err := uc.repo.GenerateToken(ctx, to)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.Login.repo.GenerateToken: %v", err)
		return LoginOutput{}, err
	}
	output := LoginOutput{
		Token: token,
		User:  uc.getProfile(ctx, user),
	}

	return output, nil
}

func (uc implUseCase) Profile(ctx context.Context, sc models.Scope) (ProfileOutput, error) {
	do := repository.DetailOptions{
		UserID: sc.UserID,
	}

	user, err := uc.repo.GetUser(ctx, do)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.Profile.repo.GetUserByID: %v", err)
		return ProfileOutput{}, err
	}

	output := uc.getProfile(ctx, user)

	return output, nil
}

func (uc implUseCase) UpdateProfile(ctx context.Context, sc models.Scope, input UpdateProfileInput) error {
	do := repository.DetailOptions{
		UserID: sc.UserID,
	}

	user, err := uc.repo.GetUser(ctx, do)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.UpdateProfile.repo.GetUserByID: %v", err)
		return err
	}

	if user.ID.Hex() == "" {
		uc.l.Errorf(ctx, "user.usecase.UpdateProfile.repo.GetUserByID: %v", ErrUserNotFound)
		return ErrUserNotFound
	}

	uo := repository.UpdateOptions{
		ID:       user.ID.Hex(),
		FullName: input.Name,
		Priority: input.Priority,
	}

	err = uc.repo.UpdateUser(ctx, uo)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.UpdateProfile.repo.UpdateUser: %v", err)
		return err
	}

	return nil
}

func (uc implUseCase) UpdatePassword(ctx context.Context, sc models.Scope, input UpdatePasswordInput) error {
	do := repository.DetailOptions{
		UserID: sc.UserID,
	}

	user, err := uc.repo.GetUser(ctx, do)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.UpdatePassword.repo.GetUserByID: %v", err)
		return err
	}

	if user.ID.Hex() == "" {
		uc.l.Errorf(ctx, "user.usecase.UpdatePassword.repo.GetUserByID: %v", ErrUserNotFound)
		return ErrUserNotFound
	}

	if !checkPassword(input.OldPassword, user.Password) {
		uc.l.Errorf(ctx, "user.usecase.UpdatePassword.repo.GetUserByID: %v", ErrInvalidValidation)
		return ErrInvalidValidation
	}

	hashedPassword := hashPassword(input.NewPassword)

	uo := repository.UpdateOptions{
		ID:       user.ID.Hex(),
		Password: hashedPassword,
	}

	err = uc.repo.UpdateUser(ctx, uo)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.UpdatePassword.repo.UpdateUser: %v", err)
		return err
	}

	return nil
}

func (uc implUseCase) GetUserById(ctx context.Context, sc models.Scope, id string) (models.User, error) {
	do := repository.DetailOptions{
		UserID: id,
	}
	user, err := uc.repo.GetUser(ctx, do)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.GetUserById.repo.GetUser: %v", err)
		return models.User{}, ErrUserNotFound
	}
	return user, nil
}

func (uc implUseCase) UpdateUserByID(ctx context.Context, sc models.Scope, input UpdateUserByIDInput, id string) (models.User, error) {
	do := repository.DetailOptions{
		UserID: id,
	}

	user, err := uc.repo.GetUser(ctx, do)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.UpdateProfile.repo.GetUserByID: %v", err)
		return models.User{}, err
	}

	if user.ID.Hex() == "" {
		uc.l.Errorf(ctx, "user.usecase.UpdateProfile.repo.GetUserByID: %v", ErrUserNotFound)
		return models.User{}, ErrUserNotFound
	}

	uo := repository.UpdateOptions{
		ID:          user.ID.Hex(),
		FullName:    input.Name,
		LastOrderAt: input.LastOrderAt,
	}

	err = uc.repo.UpdateUser(ctx, uo)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.UpdateProfile.repo.UpdateUser: %v", err)
		return models.User{}, err
	}

	return user, nil
}

func (uc implUseCase) CreateActivity(ctx context.Context, sc models.Scope, input CreateActivityInput) error {
	createOpt := repository.CreateActivityOptions{
		UserID: sc.UserID,
		Type:   input.Type,
	}
	err := uc.repo.CreateActivity(ctx, createOpt)
	if err != nil {
		return err
	}

	return nil
}

func (uc implUseCase) List(ctx context.Context, sc models.Scope, input ListInput) (ListOutput, error) {
	do := repository.ListOptions{
		PaginatorQuery: input.PaginatorQuery,
		Filter: repository.ListFilterOptions{
			IDs:  input.Filter.IDs,
			Name: input.Filter.Name,
		},
	}

	users, paginator, err := uc.repo.List(ctx, sc, do)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.List.repo.List: %v", err)
		return ListOutput{}, err
	}

	return ListOutput{
		Users:     users,
		Paginator: paginator,
	}, nil
}

func (uc implUseCase) ListByIDS(ctx context.Context, sc models.Scope, user_ids []string) ([]models.User, error) {
	users, err := uc.repo.ListByIDS(ctx, sc, user_ids)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.List.repo.List: %v", err)
		return []models.User{}, err
	}

	return users, nil
}

func (uc implUseCase) UpdateByAdmin(ctx context.Context, sc models.Scope, input UpdateByAdminInput, id string) (models.User, error) {

	opt := repository.UpdateAdminOptions{
		FullName: input.FullName,
		Email:    input.Email,
		Phone:    input.Phone,
		Priority: &input.Priority,
	}

	checkUser, err := uc.repo.GetUser(ctx, repository.DetailOptions{UserID: id})
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.UpdateByAdmin.repo.GetUser: %v", err)
		return models.User{}, err
	}
	if input.GroupID != "" {
		if checkUser.GroupID.Hex() != input.GroupID {
			group, err := uc.repo.GetGroup(ctx, sc, input.GroupID)
			if err != nil {
				uc.l.Errorf(ctx, "user.usecase.UpdateByAdmin.repo.GetGroupByID: %v", err)
				return models.User{}, err
			}
			opt.GroupID = input.GroupID
			opt.GroupName = group.Name
			opt.GroupRole = group.Role
		}
	}

	opt.Discount = input.Discount

	user, err := uc.repo.UpdateByAdmin(ctx, sc, opt, id)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.UpdateByAdmin.repo.UpdateByAdmin: %v", err)
		return models.User{}, err
	}

	return user, nil
}

func (uc implUseCase) UpdateCreditUser(ctx context.Context, sc models.Scope, credit float64, id string) error {
	err := uc.repo.UpdateCreditUser(ctx, sc, credit, id)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.UpdateCredit.repo.UpdateCredit: %v", err)
		return err
	}

	return nil
}

func (uc implUseCase) UpdateCreditUserByOrder(ctx context.Context, sc models.Scope, credit float64, id string) error {
	err := uc.repo.UpdateCreditUserByOrder(ctx, sc, credit, id)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.UpdateCredit.repo.UpdateCreditByOrder: %v", err)
		return err
	}

	return nil
}

func (uc implUseCase) CreateApiKeyUser(ctx context.Context, sc models.Scope, id string) (models.User, error) {
	user, err := uc.GetUserById(ctx, sc, id)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.CreateApiKeyUser.GetUserByID: %v", err)
		return models.User{}, err
	}

	uScope := models.Scope{
		UserID:    user.ID.Hex(),
		GroupID:   user.GroupID.Hex(),
		GroupRole: user.GroupRole,
	}
	apiKey, err := jwt.CreateApiKey(uScope, uc.encrypter)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.CreateApiKeyUser.CreateApiKey: %v", err)
		return models.User{}, err
	}

	err = uc.repo.UpdateApiKeyUser(ctx, sc, id, apiKey)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.CreateApiKeyUser.UpdateApiKey: %v", err)
		return models.User{}, err
	}

	user.ApiKey = apiKey
	return user, nil
}
