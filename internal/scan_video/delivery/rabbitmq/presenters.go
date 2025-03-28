package rabbitmq

type BookMsg struct {
	Msg string
}
type ScanDouyinVideosMsg struct {
	Mid        int64
	SecUserID  string
	VideoCount int
	NewFlag    bool
	Group      int
	DomainAPI  string
}
