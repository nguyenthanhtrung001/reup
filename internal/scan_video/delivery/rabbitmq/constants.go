package rabbitmq

import "github.com/nguyenthanhtrung001/reup/pkg/rabbitmq"

const (
	ScanVideoNewExcName = "scan_video_new_exc"
	ScanVideoNewQueue   = "scan_video_new_queue"

	ScanVideoOldExcName = "scan_video_old_exc"
	ScanVideoOldQueue   = "scan_video_old_queue"

	ScanVideoManualExcName = "scan_video_manual_exc"
	ScanVideoManualQueue   = "scan_video_manual_queue"
)

var (
	ScanVideoNewExc = rabbitmq.ExchangeArgs{
		Name:       ScanVideoNewExcName,
		Type:       rabbitmq.ExchangeTypeFanout,
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
	}
	ScanVideoOldExc = rabbitmq.ExchangeArgs{
		Name:       ScanVideoOldExcName,
		Type:       rabbitmq.ExchangeTypeFanout,
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
	}
	ScanVideoManualExc = rabbitmq.ExchangeArgs{
		Name:       ScanVideoManualExcName,
		Type:       rabbitmq.ExchangeTypeFanout,
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
	}
)
