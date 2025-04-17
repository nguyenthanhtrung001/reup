package rabbitmq

type BookMsg struct {
	Msg string
}

type ScanDouyinVideosMsg struct {
	ArrChannel []ArrChannel
	IsScanFull bool
	Group      int
}
type ArrChannel struct {
	SpaceId   int64
	ChannelId []string
}
