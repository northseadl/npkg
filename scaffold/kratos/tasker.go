package kratos

import (
	"context"
	"github.com/northseadl/go-utils/maps"
	"github.com/pkg/errors"
	"log"
	"runtime/debug"
	"time"
)

type TaskOptions struct {
	TimeOut        time.Duration
	EnableRecovery bool
}

type Task interface {
	Run()
}

type TaskFunc struct {
	Name           string
	Desc           string
	task           Task
	timeout        time.Duration
	enableRecovery bool
}

// TaskerManager 是一个HashMap任务管理器, 日志输出到std, 由k8s管理
type TaskerManager struct {
	tasks map[string]TaskFunc
}

func NewTaskerManager() *TaskerManager {
	return &TaskerManager{
		tasks: make(map[string]TaskFunc),
	}
}

func (t *TaskerManager) AddFunc(name string, desc string, task Task, opts ...TaskOptions) {
	if t.tasks == nil {
		t.tasks = make(map[string]TaskFunc)
	}
	it := TaskFunc{
		Name: name,
		Desc: desc,
		task: task,
	}
	if len(opts) > 0 {
		opt := &opts[0]
		if opt.TimeOut > 0 {
			it.timeout = opt.TimeOut
		}
		if opt.EnableRecovery {
			it.enableRecovery = true
		}
	}
	t.tasks[name] = it
}

func (t *TaskerManager) Run(name string) error {
	task, ok := t.tasks[name]
	if !ok {
		return errors.New("没有该指令")
	}

	// 运行任务, 指定超时
	cancel := func() {}
	ctx := context.Background()
	done := make(chan struct{}, 1)
	if task.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, task.timeout)
	}
	defer cancel()

	go func(ctx context.Context) {
		log.Printf("开始执行 %s", task.Name)
		start := time.Now()
		// 是否启用Recovery
		if task.enableRecovery {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("执行中断, 遇到错误 %s", err.(error).Error())
					debug.PrintStack()
				}
			}()
		}

		// 执行任务
		task.task.Run()

		// 任务结束
		log.Printf("执行结束, 本次任务耗时 %.2f 秒", time.Since(start).Seconds())
		done <- struct{}{}
	}(ctx)

	select {
	case <-done:
	case <-ctx.Done():
		log.Println("执行超时, 任务失败")
	}

	return nil
}

func (t *TaskerManager) List() []TaskFunc {
	return maps.Values(t.tasks)
}
