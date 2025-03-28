package models

import "time"

type DouyinScanChannel struct {
	ID             int64     `bson:"id"`
	SecUID         string    `bson:"sec_uid"`
	TotalFollower  int       `bson:"total_follower"`
	TotalFavorited int       `bson:"total_favorited"`
	CreatedAt      time.Time `bson:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at"`
}
