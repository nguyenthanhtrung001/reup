package usecase

type HandleDouyinWebhookInput struct {
	ChannelID string      `json:"channel_id"`
	Posts     []VideoData `json:"posts"`
}

type VideoData struct {
	AwemeID     string    `json:"aweme_id"`
	CoverURL    string    `json:"cover_url"`
	Duration    int64     `json:"duration"`
	Description string    `json:"description"`
	CreateTime  int64     `json:"create_time"`
	VideoURI    string    `json:"video_uri"`
	VideoTag    []TagData `json:"video_tag"`
}

type TagData struct {
	TagName string `json:"tag_name"`
}
