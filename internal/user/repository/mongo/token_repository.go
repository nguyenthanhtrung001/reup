package mongo

import "github.com/nguyenthanhtrung001/reup/pkg/mongo"

const (
	userTokenCollection = "user_tokens"
)

func (repo implRepository) getUserTokenCollection() mongo.Collection {
	return repo.database.Collection(userTokenCollection)
}
