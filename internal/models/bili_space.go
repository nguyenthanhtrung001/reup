package models

import "time"

type BiliSpace struct {
	Mid                int64     `bson:"mid"`
	SpaceTitle         string    `bson:"space_title"`
	SpaceThumb         string    `bson:"space_thumb"`
	LastScan           int64     `bson:"last_scan"`
	Note               string    `bson:"note"`
	Username           string    `bson:"username"`
	CurrentPage        int       `bson:"current_page"`
	ScanError          int       `bson:"scan_error"`
	UpdatedAt          time.Time `bson:"updated_at"`
	CreatedAt          time.Time `bson:"created_at"`
	CountVideo         int       `bson:"count_video"`
	CountUpload        int       `bson:"count_upload"`
	CountGet           int       `bson:"count_get"`
	CountNext          int       `bson:"count_next"`
	LastUpload         int64     `bson:"last_upload"`
	Type               string    `bson:"type"`
	DouyinLink         string    `bson:"douyin_link"`
	IxiguaLastScan     *int64    `bson:"ixigua_last_scan,omitempty"`
	FirstJoinScanCount *int64    `bson:"first_join_scan_count,omitempty"`
	DouyinLastScan     *int64    `bson:"douyin_last_scan,omitempty"`
	DouyinWaitScan     *int      `bson:"douyin_wait_scan,omitempty"`
}
