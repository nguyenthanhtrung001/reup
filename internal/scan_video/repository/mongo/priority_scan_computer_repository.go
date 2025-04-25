package mogo

import (
	"context"
	"fmt"

	"github.com/nguyenthanhtrung001/reup/internal/models"
	"github.com/nguyenthanhtrung001/reup/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	priorityCollection = "priority_scan_computer"
)

func (repo implRepository) getPriorityCollection() mongo.Collection {
	return repo.database.Collection(priorityCollection)
}

func (repo implRepository) FindAllPriorityScanComputers(ctx context.Context) ([]models.PriorityScanComputer, error) {
	col := repo.getPriorityCollection()

	// Tạo options để lấy tất cả các bản ghi
	findOptions := options.Find()

	// Thực hiện truy vấn để lấy tất cả các bản ghi có ComputerName khác rỗng
	cursor, err := col.Find(ctx, bson.M{"computer_name": bson.M{"$ne": ""}}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("could not find priority scan computers: %v", err)
	}
	defer cursor.Close(ctx)

	// Duyệt qua các kết quả trả về
	var computers []models.PriorityScanComputer
	for cursor.Next(ctx) {
		var computer models.PriorityScanComputer
		if err := cursor.Decode(&computer); err != nil {
			return nil, fmt.Errorf("cursor decode error: %v", err)
		}
		computers = append(computers, computer)
	}

	return computers, nil
}
