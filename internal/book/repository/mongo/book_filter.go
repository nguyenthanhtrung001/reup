package mogo

import (
	"context"
	"reup/internal/book/repository"
	"reup/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

func (repo implRepository) buildListQuery(ctx context.Context, sc models.Scope, withSoftDel bool, opt repository.ListOptions) (bson.M, error) {
	query := bson.M{}
	// if sc.GroupRole
	if sc.IsAdmin() {
		// Giới hạn truy vấn theo vai trò người dùng
		repo.l.Infof(ctx, "User is admin, applying role-based query")
	}

	if title := opt.Filter.Title; title != "" {
		query["title"] = bson.M{"$regex": title, "$options": "i"}
	}
	if author := opt.Filter.Author; author != "" {
		query["author"] = bson.M{"$regex": author, "$options": "i"}
	}
	// xóa mềm nếu
	if withSoftDel {
		query["isDeleted"] = bson.M{"$ne": true} // Chỉ lấy các bản ghi không bị xóa (isDeleted != true)
	}
	return query, nil
}
