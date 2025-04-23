package usecase

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"reup/internal/models"
)

func (uc implUseCase) HandleDouyinWebhook(c interface{}, data HandleDouyinWebhookInput) (int, error) {
	linkDouyin := "https://www.douyin.com/user/" + data.ChannelID
	ctx := c.(context.Context)
	billiSpace, err := uc.repo.GetBiliSpaceByDouyinLink(ctx, linkDouyin)
	if err != nil || billiSpace == nil {
		log.Printf("Error fetching BiliSpace: %v", err)
		return 0, errors.New("Error fetching BiliSpace")
	}

	videoIds := extractVideoIds(data.Posts)
	existingVideos, err := uc.repo.GetExistingVideoIds(ctx, videoIds)
	if err != nil {
		log.Printf("Error fetching existing video IDs: %v", err)
		return 0, errors.New("Error fetching existing video IDs")
	}

	var insertData []models.BiliVideo
	for _, video := range data.Posts {
		if !contains(existingVideos, video.AwemeID) {
			tags := extractTags(video.VideoTag)
			videoObj := models.BiliVideo{
				VideoID:           convertStringToInt64(video.AwemeID),
				Mid:               billiSpace.Mid,
				VideoThumb:        video.CoverURL,
				UploadTitle:       video.Description,
				Duration:          int(video.Duration / 1000),
				UploadDescription: video.Description,
				Type:              "douyin",
				UpdatedAt:         time.Now(),
				CreatedAt:         time.Now(),
				DateMake:          time.Unix(video.CreateTime, 0),
				DownloadURL:       "https://aweme.snssdk.com/aweme/v1/play/?video_id=" + video.VideoURI + "&line=0&ratio=1080&media_type=4&vr_type=0&improve_bitrate=0&is_play_url=1&is_support_h265=0&source=PackSourceEnum_PUBLISH",
				Tags:              tags,
			}
			insertData = append(insertData, videoObj)
		}
	}

	if len(insertData) > 0 {
		if err := uc.repo.InsertBiliVideos(ctx, insertData); err != nil {
			log.Printf("Error inserting videos: %v", err)
			return 0, errors.New("Error inserting videos")
		}
	}

	now := time.Now().Unix()
	billiSpace.DouyinLastScan = &now
	billiSpace.LastScan = time.Now().Unix()
	if err := uc.repo.UpdateBiliSpace(ctx, billiSpace.Mid, int(*billiSpace.DouyinLastScan), int(billiSpace.LastScan)); err != nil {
		log.Printf("Error updating BiliSpace: %v", err)
		return 0, errors.New("Error updating BiliSpace")
	}

	return len(insertData), nil
}

func extractVideoIds(posts []VideoData) []string {
	var videoIds []string
	for _, post := range posts {
		videoIds = append(videoIds, post.AwemeID)
	}
	return videoIds
}

func contains(videoMap map[string]struct{}, item string) bool {
	_, exists := videoMap[item]
	return exists
}

func convertStringToInt64(s string) int64 {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Printf("Error converting string to int64: %v", err)
		return 0
	}
	return id
}

func extractTags(tags []TagData) string {
	var tagNames []string
	for _, tag := range tags {
		tagNames = append(tagNames, tag.TagName)
	}
	return strings.Join(tagNames, ", ")
}
