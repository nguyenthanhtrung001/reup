package mongo

import "book-store/pkg/mongo"

const (
	userTokenCollection = "user_tokens"
)

func (repo implRepository) getUserTokenCollection() mongo.Collection {
	return repo.database.Collection(userTokenCollection)
}
