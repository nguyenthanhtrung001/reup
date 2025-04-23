package mongo

import (
	"context"

	"github.com/nguyenthanhtrung001/reup/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ObjectIDFromHexOrNil returns an ObjectID from the provided hex representation.
func ObjectIDFromHexOrNil(id string) primitive.ObjectID {
	objID, _ := primitive.ObjectIDFromHex(id)
	return objID
}

func ObjectIDsFromHexs(ids []string) ([]primitive.ObjectID, error) {
	objIDs := make([]primitive.ObjectID, len(ids))
	for i, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objIDs[i] = objID
	}
	return objIDs, nil
}

func BuildScopeWithRoleQuery(ctx context.Context, sc models.Scope) bson.M {
	filter := bson.M{}
	if sc.GroupRole == (models.ListGroupRole["reseller"]) {
		priUserID, err := primitive.ObjectIDFromHex(sc.UserID)
		if err != nil {
			return filter
		}
		filter["referer_id"] = priUserID
	}

	return filter
}

func BuildQueryWithSoftDelete(query bson.M) bson.M {
	query["deleted_at"] = nil
	return query
}

type SortOption struct {
	Field string
	Order int
}

func BuildSorts(opts []SortOption) bson.D {
	sorts := bson.D{}

	for _, s := range opts {
		sorts = append(sorts, bson.E{Key: s.Field, Value: s.Order})
	}

	return sorts
}
