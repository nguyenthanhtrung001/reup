package job

import (
	"github.com/nguyenthanhtrung001/reup/pkg/cron"
)

func (h Handler) Register() []cron.JobInfo {
	return []cron.JobInfo{
		{CronTime: "* * * * *", Handler: h.JobScanVideo},
		{CronTime: "* * * * *", Handler: h.JobScanVideoFullPage},
		{CronTime: "* * * * *", Handler: h.JobScanVideoManual},
	}
}
