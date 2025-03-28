package models

import "time"

type Channel struct {
	ID                        int        `bson:"id"`
	Username                  string     `bson:"username"`
	Passwd                    string     `bson:"passwd"`
	Recover                   string     `bson:"recover"`
	Live                      int        `bson:"live"`
	ChannelID                 string     `bson:"channel_id"`
	ChannelTitle              string     `bson:"channel_title"`
	ChannelThumb              string     `bson:"channel_thumb"`
	View                      int        `bson:"view"`
	Subscribe                 int        `bson:"subscribe"`
	LastUpdate                int64      `bson:"last_update"`
	LastGet                   int64      `bson:"last_get"`
	UpdatedAt                 time.Time  `bson:"updated_at"`
	CreatedAt                 time.Time  `bson:"created_at"`
	LastScanNote              int64      `bson:"last_scan_note"`
	LastView                  int        `bson:"last_view"`
	LastSubscribe             int        `bson:"last_subscribe"`
	Computer                  string     `bson:"computer"`
	Confirm                   int        `bson:"confirm"`
	NextToday                 int        `bson:"next_today"`
	CopyRight                 int        `bson:"copy_right"`
	MaxUploadDay              int        `bson:"max_upload_day"`
	CheckFrom                 int        `bson:"check_from"`
	CheckTo                   int        `bson:"check_to"`
	CheckCopyright            int        `bson:"check_copyright"`
	UploadToday               int        `bson:"upload_today"`
	Ads                       int        `bson:"ads"`
	NeedCheckCopyright        int        `bson:"need_check_copyright"`
	VideoLeft                 int        `bson:"video_left"`
	Proxy                     string     `bson:"proxy"`
	IsUpload                  int        `bson:"is_upload"`
	Priority                  int        `bson:"priority"`
	LastUpload                int64      `bson:"last_upload"`
	GetToday                  int        `bson:"get_today"`
	NewOnly                   int        `bson:"new_only"`
	IsBanned                  *int       `bson:"is_banned,omitempty"`
	TimeBanned                *time.Time `bson:"time_banned,omitempty"`
	DouyinUploadNewLimit      *int       `bson:"douyin_upload_new_limit,omitempty"`
	DouyinErrorUploadDayLimit *int       `bson:"douyin_error_upload_day_limit,omitempty"`
	DouyinErrorUploadToday    *int       `bson:"douyin_error_upload_today,omitempty"`
	DouyinUploadOldLimit      *int       `bson:"douyin_upload_old_limit,omitempty"`
	DouyinUploadOldToday      *int       `bson:"douyin_upload_old_today,omitempty"`
}
