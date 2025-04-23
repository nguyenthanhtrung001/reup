package models

import "time"

type BiliVideo struct {
	VideoID           int64     `bson:"video_id"`
	Mid               int64     `bson:"mid"`
	VideoThumb        string    `bson:"video_thumb"`
	UploadTitle       string    `bson:"upload_title"`
	UploadDescription string    `bson:"upload_description"`
	UploadKeyword     string    `bson:"upload_keyword"`
	Duration          int       `bson:"duration"`
	DateMake          time.Time `bson:"date_make"`
	CountGet          int       `bson:"count_get"`
	Next              int       `bson:"next"`
	DownloadFail      int       `bson:"download_fail"`
	Type              string    `bson:"type"`
	UpdatedAt         time.Time `bson:"updated_at"`
	CreatedAt         time.Time `bson:"created_at"`
	DownloadURL       string    `bson:"download_url"`
	Tags              string    `bson:"tags"`
}
