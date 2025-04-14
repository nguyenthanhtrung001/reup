package mogo

import (
	"context"
	"fmt"
	"log"
	"time"

	"reup/internal/models"
	"reup/internal/scan_video/repository"
	"reup/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	biliSpaceCollection = "bili_space"
)

func (repo implRepository) getBiliSpaceCollection() mongo.Collection {
	return repo.database.Collection(biliSpaceCollection)
}

func (repo implRepository) FindDouyinOldSpaces(ctx context.Context, opts *repository.FindDouyinOldSpacesOptions) ([]models.BiliSpace, error) {
	// Lấy collection bili_space
	col := repo.getBiliSpaceCollection()

	// Tạo filter cơ bản
	filter := repo.buildFilter(opts)
	// Tách việc tạo filter ra hàm riêng để tái sử dụng
	if opts.DouyinLastScanIsZero != nil {
		if *opts.DouyinLastScanIsZero {
			// Lọc các bản ghi có douyin_last_scan == 0
			filter["douyin_last_scan"] = 0
		} else {
			// Lọc các bản ghi có douyin_last_scan khác 0 hoặc rỗng
			filter["douyin_last_scan"] = bson.M{"$ne": 0}
		}
	}
	// Tạo options để sắp xếp và giới hạn số lượng kết quả
	findOptions := options.Find()
	findOptions.SetSort(bson.D{
		{Key: "douyin_last_scan", Value: 1}, // Sắp xếp theo douyin_last_scan (asc)
		{Key: "last_scan", Value: 1},        // Sắp xếp theo last_scan (asc)
	})
	findOptions.SetLimit(opts.ScanNumbers) // Giới hạn số lượng kết quả

	// Thực hiện truy vấn
	cursor, err := col.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("could not find douyin old spaces: %v", err)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			// Log error when closing the cursor
			log.Printf("Error closing cursor: %v", err)
		}
	}()

	// Duyệt qua các kết quả trả về
	var douyinOldSpaces []models.BiliSpace
	for cursor.Next(ctx) {
		var biliSpace models.BiliSpace
		if err := cursor.Decode(&biliSpace); err != nil {
			return nil, fmt.Errorf("cursor decode error: %v", err)
		}
		douyinOldSpaces = append(douyinOldSpaces, biliSpace)
	}

	return douyinOldSpaces, nil
}

func (repo implRepository) CreateBiliSpace(ctx context.Context, input repository.CreateBiliSpaceInput) (*models.BiliSpace, error) {
	col := repo.getBiliSpaceCollection()

	// Chuyển đổi từ input sang model BiliSpace
	biliSpace := models.BiliSpace{
		Mid:                input.Mid,
		SpaceTitle:         input.SpaceTitle,
		SpaceThumb:         input.SpaceThumb,
		LastScan:           input.LastScan,
		Note:               input.Note,
		Username:           input.Username,
		CurrentPage:        input.CurrentPage,
		ScanError:          input.ScanError,
		UpdatedAt:          input.UpdatedAt,
		CreatedAt:          input.CreatedAt,
		CountVideo:         input.CountVideo,
		CountUpload:        input.CountUpload,
		CountGet:           input.CountGet,
		CountNext:          input.CountNext,
		LastUpload:         input.LastUpload,
		Type:               input.Type,
		DouyinLink:         input.DouyinLink,
		IxiguaLastScan:     input.IxiguaLastScan,
		FirstJoinScanCount: input.FirstJoinScanCount,
		DouyinLastScan:     input.DouyinLastScan,
		DouyinWaitScan:     input.DouyinWaitScan,
	}

	// Thực hiện insert BiliSpace mới vào MongoDB
	_, err := col.InsertOne(ctx, biliSpace)
	if err != nil {
		return nil, fmt.Errorf("error inserting bili space: %v", err)
	}

	return &biliSpace, nil
}

func (repo implRepository) UpdateBiliSpaceF(ctx context.Context, mid int64, updatedFields map[string]interface{}) error {
	col := repo.getBiliSpaceCollection()

	filter := bson.M{"mid": mid}

	update := bson.M{
		"$set": updatedFields,
	}

	_, err := col.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating bili space: %v", err)
	}

	return nil
}

func (repo implRepository) UpdateBiliSpace(ctx context.Context, mid int64, douyinWaitScan int, scanError int) error {
	col := repo.getBiliSpaceCollection()

	_, err := col.UpdateOne(ctx, bson.M{"mid": mid}, bson.M{
		"$set": bson.M{
			"douyin_last_scan": time.Now().Unix(),
			"last_scan":        time.Now().Unix(),
			"douyin_wait_scan": douyinWaitScan,
			"current_page":     0,
			"scan_error":       scanError,
		},
	})

	if err != nil {
		return fmt.Errorf("error updating bili space: %v", err)
	}

	return nil
}
