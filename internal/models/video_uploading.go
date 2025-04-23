package models

import "time"

type VideoUploading struct {
	VideoID   int64     `bson:"video_id" json:"video_id"`
	Username  string    `bson:"username" json:"username"`
	StartTime int64     `bson:"start_time" json:"start_time"`
	Mid       int64     `bson:"mid" json:"mid"`
	Computer  string    `bson:"computer" json:"computer"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
