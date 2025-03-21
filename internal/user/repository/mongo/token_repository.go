package mongo

import "reup/pkg/mongo"

const (
	userTokenCollection = "user_tokens"
)

func (repo implRepository) getUserTokenCollection() mongo.Collection {
	return repo.database.Collection(userTokenCollection)
}
