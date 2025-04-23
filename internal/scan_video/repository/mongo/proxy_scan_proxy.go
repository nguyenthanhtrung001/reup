package mogo

import (
	"context"
	"fmt"
	"time"

	"github.com/nguyenthanhtrung001/reup/internal/models"
	"github.com/nguyenthanhtrung001/reup/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	proxyScanCollection = "proxy_scan"
)

func (repo implRepository) getProxyScanCollection() mongo.Collection {
	return repo.database.Collection(proxyScanCollection)
}

// DoneAllProxyScan implements repository.Repository.
func (repo implRepository) DoneAllProxyScan(ctx context.Context) error {

	col := repo.getProxyScanCollection()
	_, err := col.UpdateMany(ctx, bson.M{
		"live":   1,
		"enable": 1,
		"status": 0,
	}, bson.M{
		"$set": bson.M{
			"status": 1,
		},
	})
	if err != nil {
		return fmt.Errorf("error updating all proxy scan: %v", err)
	}
	return nil
}

// DoneProxyScan implements repository.Repository.
func (repo implRepository) DoneProxyScan(ctx context.Context, proxyIP string) error {

	col := repo.getProxyScanCollection()
	_, err := col.UpdateOne(ctx, bson.M{
		"proxy_ip": proxyIP,
	}, bson.M{
		"$set": bson.M{
			"status": 1,
		},
	})
	if err != nil {
		return fmt.Errorf("error updating proxy scan: %v", err)
	}
	return nil

}

// GetAllProxyScan implements repository.Repository.
func (repo implRepository) GetAllProxyScan(ctx context.Context) ([]models.ProxyScan, error) {
	col := repo.getProxyScanCollection()
	cursor, err := col.Find(ctx, bson.M{
		"live":   1,
		"enable": 1,
		"status": 1,
	})
	if err != nil {
		repo.l.Error(ctx, "error fetching all proxy: %v", err)
		return nil, fmt.Errorf("error fetching all proxy: %v", err)
	}
	defer cursor.Close(ctx)

	var proxies []models.ProxyScan
	for cursor.Next(ctx) {
		var proxy models.ProxyScan
		err := cursor.Decode(&proxy)
		if err != nil {
			repo.l.Error(ctx, "error decoding proxy: %v", err)
			return nil, fmt.Errorf("error decoding proxy: %v", err)
		}
		proxies = append(proxies, proxy)
	}
	_, err = col.UpdateMany(ctx, bson.M{
		"live":   1,
		"enable": 1,
		"status": 1,
	}, bson.M{
		"$inc": bson.M{"count_used": 1},
		"$set": bson.M{
			"last_call": time.Now(),
			"status":    0,
		},
	})
	if err != nil {
		repo.l.Error(ctx, "error updating proxy: %v", err)
		return nil, fmt.Errorf("error updating proxy: %v", err)
	}

	return proxies, nil
}

// GetProxyScanRandom implements repository.Repository.
func (repo implRepository) GetProxyScanRandom(ctx context.Context) (*models.ProxyScan, error) {
	col := repo.getProxyScanCollection()

	findOptions := options.FindOne().SetSort(bson.D{{Key: "last_call", Value: 1}})

	var proxy models.ProxyScan
	err := col.FindOneWithOpt(ctx, bson.M{
		"live":   1,
		"enable": 1,
		"status": 1,
	}, findOptions).Decode(&proxy)

	if err != nil {
		return nil, fmt.Errorf("error fetching random proxy: %v", err)
	}

	_, err = col.UpdateOne(ctx,
		bson.M{"_id": proxy.ID},
		bson.M{
			"$inc": bson.M{"count_used": 1},
			"$set": bson.M{
				"last_call": time.Now(),
				"status":    0,
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("error updating proxy: %v", err)
	}

	return &proxy, nil
}

// InsertProxyScan implements repository.Repository.
func (repo implRepository) InsertProxyScan(ctx context.Context, proxyIP string) error {
	col := repo.getProxyScanCollection()

	// Khởi tạo proxy mới
	proxyScan := models.ProxyScan{
		ID:        primitive.NewObjectID().Hex(), // Nếu cần string ID
		ProxyIP:   proxyIP,
		Live:      1,
		Enable:    1,
		Status:    1,
		CountUsed: 0,
		LastCall:  time.Now(),
		Username:  "",
	}

	// Insert vào MongoDB
	result, err := col.InsertOne(ctx, proxyScan)
	if err != nil {
		fmt.Printf("❌ Insert thất bại: %v\n", err)
		return fmt.Errorf("error inserting proxy scan: %v", err)
	}

	fmt.Printf("✅ Insert thành công! InsertedID: %v | ProxyIP: %s\n", result, proxyScan.ProxyIP)
	return nil
}
