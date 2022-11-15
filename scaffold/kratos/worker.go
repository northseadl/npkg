package kratos

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"sync"
	"time"
)

// Worker 是一个为了执行持续任务的任务管理器, 且带有定时执行功能
type Worker struct {
	cron  *cron.Cron
	works []func()
	h     *log.Helper
}

func NewWorker(l *log.Helper) Worker {
	return Worker{
		h:    l,
		cron: cron.New(),
	}
}

func (w *Worker) AddFunc(f func()) *Worker {
	w.works = append(w.works, f)
	return w
}

func (w *Worker) AddCronFunc(spec string, f func()) (cron.EntryID, error) {
	return w.cron.AddFunc(spec, f)
}

func (w *Worker) AddCronJob(spec string, job cron.Job) (cron.EntryID, error) {
	return w.cron.AddJob(spec, job)
}

func (w *Worker) RemoveEntry(id cron.EntryID) {
	w.cron.Remove(id)
}

func (w *Worker) Run() {
	// run cron
	w.cron.Start()
	defer w.cron.Stop()

	// run works
	group := new(sync.WaitGroup)
	for _, work := range w.works {
		group.Add(1)
		go func(work func(), group *sync.WaitGroup) {
			defer func() {
				if err := recover(); err != nil {
					w.h.Error(errors.WithStack(err.(error)))
				}
			}()

			work()
			group.Done()
		}(work, group)
	}
	group.Wait()

	for {
		if len(w.cron.Entries()) > 0 {
			<-time.After(time.Minute * 30)
		} else {
			break
		}
	}
}
