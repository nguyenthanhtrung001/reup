package mongo

import (
	"context"
	"time"

	"book-store/internal/models"
	"book-store/internal/user/repository"
	"book-store/pkg/jwt"
	"book-store/pkg/mongo"
	"book-store/pkg/paginator"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	userCollection = "users"
)

func (repo implRepository) getUserCollection() mongo.Collection {
	return repo.database.Collection(userCollection)
}

func (repo implRepository) CreateNewUser(ctx context.Context, opt repository.RegisterOptions) error {
	col := repo.getUserCollection()

	groupID, err := primitive.ObjectIDFromHex(opt.GroupID)
	if err != nil {
		repo.l.Warnf(ctx, "user.repository.mongo.RegisterFilter.primitive.ObjectIDFromHex: %v", err)
		return err
	}

	newUser := models.User{
		ID:        primitive.NewObjectID(),
		Email:     opt.Email,
		Phone:     opt.Phone,
		Password:  opt.Password,
		FullName:  opt.FullName,
		Verified:  opt.Verified,
		GroupID:   groupID,
		GroupRole: opt.GroupRole,
		GroupName: opt.GroupName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Check referer_id
	// if opt.RefererID != "" {
	// 	refererID, err := primitive.ObjectIDFromHex(opt.RefererID)
	// 	if err != nil {
	// 		repo.l.Warnf(ctx, "user.repository.mongo.RegisterFilter.primitive.ObjectIDFromHex: %v", err)
	// 		return err
	// 	}
	// 	newUser.RefererID = refererID
	// }

	if _, err := col.InsertOne(ctx, newUser); err != nil {
		repo.l.Warnf(ctx, "user.repository.mongo.RegisterFilter.InsertOne: %v", err)
		return err
	}

	return nil
}

func (repo implRepository) GetUser(ctx context.Context, do repository.DetailOptions) (models.User, error) {
	col := repo.getUserCollection()

	filter, err := repo.buildGetUserQuery(ctx, do, true)
	if err != nil {
		repo.l.Warnf(ctx, "user.repository.mongo.GetUserByEmail.buildGetUserByEmailQuery: %v", err)
		return models.User{}, err
	}
	cursor := col.FindOne(ctx, filter)

	var user models.User

	if err := cursor.Decode(&user); err != nil {
		repo.l.Warnf(ctx, "user.repository.mongo.GetUserByEmailOrId.Decode: %v", err)
		return models.User{}, err
	}

	return user, nil
}

func (repo implRepository) GenerateToken(ctx context.Context, opt repository.TokenOptions) (string, error) {
	payload := jwt.NewPayload(jwt.UserField{
		UserID:    opt.UserID,
		GroupID:   opt.GroupID,
		GroupRole: opt.GroupRole,
	}, 30*24*time.Hour)
	token, t, err := repo.jwtManager.CreateToken(payload)
	if err != nil {
		repo.l.Warnf(ctx, "user.repository.mongo.GenerateToken.jwtManager.CreateToken: %v", err)
		return "", err
	}

	col := repo.getUserTokenCollection()

	newToken := primitive.M{
		"user_id":   opt.UserID,
		"token":     token,
		"expiredAt": primitive.NewDateTimeFromTime(t),
	}

	_, err = col.InsertOne(ctx, newToken)
	if err != nil {
		repo.l.Warnf(ctx, "user.repository.mongo.GenerateToken.InsertOne: %v", err)
		return "", err
	}

	return token, nil
}

func (repo implRepository) UpdateUser(ctx context.Context, opt repository.UpdateOptions) error {
	col := repo.getUserCollection()

	filter, err := repo.buildUpdateUserQuery(ctx, opt, true)
	if err != nil {
		repo.l.Warnf(ctx, "user.repository.mongo.UpdateUser.buildUpdateUserQuery: %v", err)
		return err
	}
	update := repo.buildUpdateUserContent(opt)
	_, err = col.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		repo.l.Warnf(ctx, "user.repository.mongo.UpdateUser.UpdateOne: %v", err)
		return err
	}

	return nil
}

func (repor implRepository) IsExist(ctx context.Context, email string) bool {
	col := repor.getUserCollection()

	return col.FindOne(ctx, bson.M{"email": email}).Decode(&models.User{}) == nil
}

func (repo implRepository) List(ctx context.Context, sc models.Scope, opt repository.ListOptions) ([]models.User, paginator.Paginator, error) {
	col := repo.getUserCollection()

	filter, err := repo.buildListQuery(ctx, sc, true, opt)
	if err != nil {
		repo.l.Warnf(ctx, "user.repository.List.buildListQuery: %v", err)
		return nil, paginator.Paginator{}, err
	}

	cursor, err := col.Find(ctx, filter, options.Find().
		SetSkip(opt.PaginatorQuery.Offset()).
		SetLimit(opt.PaginatorQuery.Limit).
		SetSort(bson.D{
			{Key: "created_at", Value: -1},
			{Key: "_id", Value: -1},
		}),
	)
	if err != nil {
		repo.l.Warnf(ctx, "user.repository.List.col.Find: %v", err)
		return nil, paginator.Paginator{}, err
	}

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		repo.l.Warnf(ctx, "user.repository.List.cursor.All: %v", err)
		return nil, paginator.Paginator{}, err
	}

	total, err := col.CountDocuments(ctx, filter)
	if err != nil {
		repo.l.Warnf(ctx, "user.repository.List.col.CountDocuments: %v", err)
		return nil, paginator.Paginator{}, err
	}

	return users, paginator.Paginator{
		Total:       total,
		Count:       int64(len(users)),
		PerPage:     opt.PaginatorQuery.Limit,
		CurrentPage: opt.PaginatorQuery.Page,
	}, nil
}

func (repo implRepository) ListByIDS(ctx context.Context, sc models.Scope, user_ids []string) ([]models.User, error) {
	col := repo.getUserCollection()

	object_user_ids := []primitive.ObjectID{}

	for _, v := range user_ids {
		id, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			repo.l.Warnf(ctx, "user.repository.mongo.ListByIDS.primitive.ObjectIDFromHex: %v", err)
		} else {
			object_user_ids = append(object_user_ids, id)
		}
	}

	if len(object_user_ids) == 0 {
		return []models.User{}, nil
	}

	filter := bson.M{
		"_id": bson.M{
			"$in": object_user_ids,
		},
	}

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		repo.l.Warnf(ctx, "user.repository.mongo.ListByIDS.col.Find: %v", err)
		return []models.User{}, err
	}

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		repo.l.Warnf(ctx, "user.repository.mongo.ListByIDS.cursor.All: %v", err)
		return []models.User{}, err
	}

	return users, nil
}

func (repo implRepository) UpdateByAdmin(ctx context.Context, sc models.Scope, opt repository.UpdateAdminOptions, id string) (models.User, error) {
	col := repo.getUserCollection()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		repo.l.Warnf(ctx, "user.repository.mongo.UpdateAdmin.primitive.ObjectIDFromHex: %v", err)
		return models.User{}, err
	}
	filter := bson.M{"_id": objID}

	set := bson.M{}
	if opt.Email != "" {
		set["email"] = opt.Email
	}
	if opt.Phone != "" {
		set["phone"] = opt.Phone
	}
	if opt.FullName != "" {
		set["fullname"] = opt.FullName
	}

	if opt.Priority != nil {
		set["priority"] = opt.Priority
	}

	if opt.GroupID != "" && opt.GroupRole != "" && opt.GroupName != "" {
		groupID, err := primitive.ObjectIDFromHex(opt.GroupID)
		if err != nil {
			repo.l.Warnf(ctx, "user.repository.mongo.UpdateAdmin.primitive.ObjectIDFromHex: %v", err)
			return models.User{}, err
		}
		set["group_id"] = groupID
		set["group_role"] = opt.GroupRole
		set["group_name"] = opt.GroupName
	}

	if opt.Discount != nil {
		set["discount"] = *opt.Discount
	}

	_, err = col.UpdateOne(ctx, filter, bson.M{"$set": set})
	if err != nil {
		repo.l.Warnf(ctx, "user.repository.mongo.UpdateAdmin.UpdateOne: %v", err)
		return models.User{}, err
	}

	return repo.GetUser(ctx, repository.DetailOptions{UserID: id})
}

func (repo implRepository) UpdateCreditUser(ctx context.Context, sc models.Scope, credit float64, id string) error {
	col := repo.getUserCollection()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		repo.l.Warnf(ctx, "user.repository.mongo.UpdateCredit.primitive.ObjectIDFromHex: %v", err)
		return err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$inc": bson.M{"credit": credit}}
	_, err = col.UpdateOne(ctx, filter, update)
	if err != nil {
		repo.l.Warnf(ctx, "user.repository.mongo.UpdateCredit.UpdateOne: %v", err)
		return err
	}

	return nil
}

func (repo implRepository) UpdateCreditUserByOrder(ctx context.Context, sc models.Scope, credit float64, id string) error {
	col := repo.getUserCollection()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		repo.l.Warnf(ctx, "user.repository.mongo.UpdateCreditByOrder.primitive.ObjectIDFromHex: %v", err)
		return err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$inc": bson.M{
			"credit":      -credit,
			"credit_used": credit,
		},
	}
	_, err = col.UpdateOne(ctx, filter, update)
	if err != nil {
		repo.l.Warnf(ctx, "user.repository.mongo.UpdateCreditByOrder.UpdateOne: %v", err)
		return err
	}

	return nil
}

func (repo implRepository) UpdateApiKeyUser(ctx context.Context, sc models.Scope, id string, apiKey string) error {
	col := repo.getUserCollection()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		repo.l.Warnf(ctx, "user.repository.mongo.UpdateApiKeyUser.primitive.ObjectIDFromHex: %v", err)
		return err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"api_key": apiKey}}
	_, err = col.UpdateOne(ctx, filter, update)
	if err != nil {
		repo.l.Warnf(ctx, "user.repository.mongo.UpdateApiKeyUser.UpdateOne: %v", err)
		return err
	}

	return nil
}
