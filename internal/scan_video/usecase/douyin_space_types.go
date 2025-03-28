package usecase

import "time"

type FindDouyinOldSpacesInput struct {
	ChannelUsernames     []string // Danh sách usernames cần tìm kiếm
	ScanNumbers          int64    // Số lượng kết quả cần trả về
	DouyinLastScanIsZero bool
}
type CreateBiliSpaceInput struct {
	Mid                int64     `json:"mid" bson:"mid"`
	SpaceTitle         string    `json:"space_title" bson:"space_title"`
	SpaceThumb         string    `json:"space_thumb" bson:"space_thumb"`
	LastScan           int64     `json:"last_scan" bson:"last_scan"`
	Note               string    `json:"note" bson:"note"`
	Username           string    `json:"username" bson:"username"`
	CurrentPage        int       `json:"current_page" bson:"current_page"`
	ScanError          int       `json:"scan_error" bson:"scan_error"`
	UpdatedAt          time.Time `json:"updated_at" bson:"updated_at"`
	CreatedAt          time.Time `json:"created_at" bson:"created_at"`
	CountVideo         int       `json:"count_video" bson:"count_video"`
	CountUpload        int       `json:"count_upload" bson:"count_upload"`
	CountGet           int       `json:"count_get" bson:"count_get"`
	CountNext          int       `json:"count_next" bson:"count_next"`
	LastUpload         int64     `json:"last_upload" bson:"last_upload"`
	Type               string    `json:"type" bson:"type"`
	DouyinLink         string    `json:"douyin_link" bson:"douyin_link"`
	IxiguaLastScan     *int64    `json:"ixigua_last_scan,omitempty" bson:"ixigua_last_scan,omitempty"`
	FirstJoinScanCount *int64    `json:"first_join_scan_count,omitempty" bson:"first_join_scan_count,omitempty"`
	DouyinLastScan     *int64    `json:"douyin_last_scan,omitempty" bson:"douyin_last_scan,omitempty"`
	DouyinWaitScan     *int      `json:"douyin_wait_scan,omitempty" bson:"douyin_wait_scan,omitempty"`
}
