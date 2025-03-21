package mongo

import (
	"context"

	"reup/internal/models"
	"reup/internal/user/repository"
	pkgMongo "reup/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo implRepository) buildGetUserQuery(ctx context.Context, do repository.DetailOptions, withSoftDel ...bool) (bson.M, error) {
	query := bson.M{}

	if do.Email != "" {
		query["email"] = do.Email
	}

	if do.UserID != "" {
		userID, err := primitive.ObjectIDFromHex(do.UserID)
		if err != nil {
			repo.l.Warnf(ctx, "user.repository.mongo.buildGetUserQuery.primitive.ObjectIDFromHex: %v", err)
			return bson.M{}, err
		}
		query["_id"] = userID
	}

	if len(withSoftDel) > 0 && withSoftDel[0] {
		query = pkgMongo.BuildQueryWithSoftDelete(query)
	}

	return query, nil

}

func (repo implRepository) buildUpdateUserQuery(ctx context.Context, uo repository.UpdateOptions, withSoftDel ...bool) (bson.M, error) {
	query := bson.M{}

	if uo.ID != "" {
		userID, err := primitive.ObjectIDFromHex(uo.ID)
		if err != nil {
			repo.l.Warnf(ctx, "user.repository.mongo.buildUpdateUserQuery.primitive.ObjectIDFromHex: %v", err)
			return bson.M{}, err
		}
		query["_id"] = userID
	}

	if len(withSoftDel) > 0 && withSoftDel[0] {
		query = pkgMongo.BuildQueryWithSoftDelete(query)
	}

	return query, nil
}

func (repo implRepository) buildUpdateUserContent(uo repository.UpdateOptions) bson.M {
	update := bson.M{}

	if uo.FullName != "" {
		update["fullname"] = uo.FullName
	}

	if uo.Email != "" {
		update["email"] = uo.Email
	}

	if uo.Phone != "" {
		update["phone"] = uo.Phone
	}

	if uo.Password != "" {
		update["password"] = uo.Password
	}

	if uo.LastOrderAt != nil {
		update["last_order_at"] = *uo.LastOrderAt
	}

	return update
}

func (repo implRepository) buildListQuery(ctx context.Context, sc models.Scope, withSoftDel bool, opt repository.ListOptions) (bson.M, error) {
	query := pkgMongo.BuildScopeWithRoleQuery(ctx, sc)

	if withSoftDel {
		query = pkgMongo.BuildQueryWithSoftDelete(query)
	}

	if len(opt.Filter.IDs) > 0 {
		ids, err := pkgMongo.ObjectIDsFromHexs(opt.Filter.IDs)
		if err != nil {
			repo.l.Warnf(ctx, "user.repository.mongo.buildListQuery.pkgMongo.ObjectIDsFromHexs: %v", err)
			return bson.M{}, err
		}

		query["_id"] = bson.M{"$in": ids}
	}

	if opt.Filter.Name != "" {
		query["fullname"] = bson.M{"$regex": opt.Filter.Name, "$options": "i"}
	}

	return query, nil
}
