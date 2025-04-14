package models

type AppSetting struct {
	ID                     int    `bson:"id"`
	VideoPerDay            int    `bson:"video_per_day"`
	Cookies                string `bson:"cookies"`
	KeyAuto                string `bson:"key_auto"`
	KeyManual              string `bson:"key_manual"`
	MaxNextDay             int    `bson:"max_next_day"`
	ScanHistoryLineNumber  int    `bson:"scan_history_line_number"`
	DouyinErrorUploadCount int    `bson:"douyin_error_upload_count"`
	DouyinVideoPerDay      int    `bson:"douyin_video_per_day"`
	DouyinManualScanStatus int    `bson:"douyin_manual_scan_status"`
}
