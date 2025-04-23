package usecase

import "time"

type CreateBiliVideoInput struct {
	VideoID           int64     `json:"video_id" bson:"video_id"`
	Mid               int64     `json:"mid" bson:"mid"`
	VideoThumb        string    `json:"video_thumb" bson:"video_thumb"`
	UploadTitle       string    `json:"upload_title" bson:"upload_title"`
	UploadDescription string    `json:"upload_description" bson:"upload_description"`
	UploadKeyword     string    `json:"upload_keyword" bson:"upload_keyword"`
	Duration          int       `json:"duration" bson:"duration"`
	DateMake          time.Time `json:"date_make" bson:"date_make"`
	CountGet          int       `json:"count_get" bson:"count_get"`
	Next              int       `json:"next" bson:"next"`
	DownloadFail      int       `json:"download_fail" bson:"download_fail"`
	Type              string    `json:"type" bson:"type"`
	UpdatedAt         time.Time `json:"updated_at" bson:"updated_at"`
	CreatedAt         time.Time `json:"created_at" bson:"created_at"`
}
