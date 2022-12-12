package kratos

import "github.com/go-kratos/kratos/v2/log"

type DebugLogger struct {
	h *log.Helper
}

func (l *DebugLogger) Debug(msg string) error {
	l.h.Debug(msg)
	return nil
}

func NewDebugLogger(logger log.Logger) *DebugLogger {
	return &DebugLogger{
		h: log.NewHelper(logger),
	}
}
