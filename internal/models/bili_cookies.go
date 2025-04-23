package models

import "time"

type BiliCookies struct {
	ID          int       `bson:"id" json:"id"`                     // id (INT)
	CookiesList string    `bson:"cookies_list" json:"cookies_list"` // cookies_list (TEXT)
	GetFail     int       `bson:"get_fail" json:"get_fail"`         // get_fail (INT)
	LastGet     time.Time `bson:"last_get" json:"last_get"`         // last_get (DATETIME)
	LastUpdate  int64     `bson:"last_update" json:"last_update"`   // last_update (BIGINT)
}
