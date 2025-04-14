package job

import (
	"context"
)

func (h Handler) JobScanVideo() {
	ctx := context.Background()
	h.l.Info(ctx, "Scan video ")
	h.uc.ScanDouyinVideoSheduler()

}

func (h Handler) JobScanVideoFullPage() {
	ctx := context.Background()
	h.l.Info(ctx, "Scan video new ")
	h.uc.ScanDouyinVideoFullPageSheduler()

}

func (h Handler) JobScanVideoManual() {
	ctx := context.Background()
	h.l.Info(ctx, "Scan video manual ")
	h.uc.ScanDouyinVideoManualSheduler()

}
