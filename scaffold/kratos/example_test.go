package kratos

import (
	"log"
	"time"
)

func ExampleNewCleaner() {
	cleaner := NewCleaner()
	defer cleaner.Clean()

	cleaner.AddFunc(func() {
		// do something
	})
}

func ExampleNewTaskerManager() {
	m := NewTaskerManager()
	m.AddFunc("test", "a test task", func() {
		// do something
	})
	m.AddFunc("test2", "test2 desc", func() {
		// do something
	}, TaskOptions{
		TimeOut:        time.Second * 10,
		EnableRecovery: true,
	})
	err := m.Run("test")
	if err != nil {
		log.Fatal(err)
	}
}
