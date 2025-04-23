package models

import "time"

type DouyinCookies struct {
	ID        int64     `bson:"id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Cookies   string    `bson:"cookies" json:"cookies"`
	UsedCount int       `bson:"used_count" json:"used_count"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
