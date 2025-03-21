package job

import (
	"reup/pkg/cron"
)

func (h Handler) Register() []cron.JobInfo {
	return []cron.JobInfo{
		{CronTime: "* * * * *", Handler: h.testJobBook},
	}
}
