package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nguyenthanhtrung001/reup/internal/models"
	"github.com/nguyenthanhtrung001/reup/pkg/curl"
)

func (uc implUseCase) ScanDouyinHomePageWithoutRequest() ([]models.DouyinScanChannel, error) {
	ctx := context.Background()
	uc.l.Info(ctx, "(job) scan douyin home page is in progress...")

	// Gửi request đến API
	url := "http://ixigua.sscapi.co/api/v1/douyin/get-multi-channels/"
	headers := map[string]string{
		"accept": "application/json",
	}

	response, err := curl.Get(url, headers)
	if err != nil {
		uc.l.Error(ctx, "Error sending request: %v", err)
		return nil, err
	}

	// Nếu phản hồi thành công
	if response != "" {
		// Parse JSON data
		var data map[string]interface{}
		err := json.Unmarshal([]byte(response), &data)
		if err != nil {
			uc.l.Error(ctx, "Error parsing JSON response: %v", err)
			return nil, err
		}

		channelIds := data["data"].([]interface{})
		var remainingArr []models.DouyinScanChannel

		uc.l.Info(ctx, "=========================channel_ids=====================")
		uc.l.Info(ctx, channelIds)
		uc.l.Info(ctx, "=========================channel_ids=====================")

		for _, cid := range channelIds {
			cidStr := fmt.Sprintf("%v", cid) // Convert integer to string
			// Check if the sec_uid exists in the database (Assume that we have a `FindDouyinScanChannel` method)
			existingChannel, err := uc.repo.FindDouyinScanChannelBySecUID(ctx, cidStr)
			if err != nil {
				uc.l.Error(ctx, "Error checking if Douyin scan channel exists: %v", err)
				return nil, err
			}
			if existingChannel != nil {
				continue // Ignore this loop round
			}

			// Create a new DouyinScanChannel object
			now := time.Now()
			item := models.DouyinScanChannel{
				SecUID:    cidStr,
				CreatedAt: now,
				UpdatedAt: now,
			}

			// Append to remaining array
			remainingArr = append(remainingArr, item)
		}

		// If there are remaining channels, insert them into the database
		if len(remainingArr) >= 1 {
			err := uc.repo.InsertDouyinScanChannels(ctx, remainingArr)
			if err != nil {
				uc.l.Error(ctx, "Error inserting Douyin scan channels: %v", err)
				return nil, err
			}
			log.Printf("Stored %d Douyin scan channel ids", len(remainingArr))
		}

		return remainingArr, nil
	} else {
		log.Println("Scan douyin home page failed")
		return nil, fmt.Errorf("failed to fetch data")
	}
}
