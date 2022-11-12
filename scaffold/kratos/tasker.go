package kratos

import (
	"github.com/northseadl/go-utils/maps"
	"time"
)

type TaskOptions struct {
	TimeOut time.Duration
}

type TaskFunc struct {
	Name    string
	Desc    string
	task    func()
	timeout time.Duration
}

type TaskerManager struct {
	tasks map[string]TaskFunc
}

func NewTaskerManager() TaskerManager {
	return TaskerManager{
		tasks: make(map[string]TaskFunc),
	}
}

func (t *TaskerManager) AddFunc(name string, desc string, task func(), opts ...TaskOptions) {
	if t.tasks == nil {
		t.tasks = make(map[string]TaskFunc)
	}
	t.tasks[name] = TaskFunc{
		Name: name,
		Desc: desc,
		task: task,
	}
}

func (t *TaskerManager) Run(name string) error {
	return nil
}

func (t *TaskerManager) List() []TaskFunc {
	return maps.Values(t.tasks)
}
