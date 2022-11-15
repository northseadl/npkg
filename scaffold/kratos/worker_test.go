package kratos

import (
	"github.com/go-kratos/kratos/v2/log"
	"testing"
)

func TestWorker_Run(t *testing.T) {
	h := log.NewHelper(log.DefaultLogger)
	worker := NewWorker(h)
	worker.AddFunc(func() {
		// work
	})
	_, _ = worker.AddCronFunc("* 0/1 * * *", func() {
		// cron func
	})
	worker.Run()
}
