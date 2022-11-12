package qmgo

import (
	"context"
	"fmt"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoConf struct {
	DSN         string
	Debug       bool
	MaxPoolSize uint64
	MinPoolSize uint64
}

type DebugLogger interface {
	Debug(msg string) error
}

func NewClient(ctx context.Context, conf *MongoConf, logger DebugLogger) (*qmgo.Client, error) {
	var opt mongoOptions.ClientOptions
	var monitor event.CommandMonitor
	if conf.Debug {
		monitor.Started = func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			cmdBytes, _ := bson.Marshal(startedEvent.Command)
			_ = logger.Debug(fmt.Sprintf("[rid-%d] %s", startedEvent.RequestID, string(cmdBytes)))
		}
		monitor.Succeeded = func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {
			replyBytes, _ := bson.Marshal(succeededEvent.Reply)
			_ = logger.Debug(fmt.Sprintf("[rid-%d] %s; (dur %s)", succeededEvent.RequestID, string(replyBytes), time.Duration(succeededEvent.DurationNanos).String()))
		}
		opt.SetMonitor(&monitor)
	}
	if conf.MaxPoolSize == 0 {
		conf.MaxPoolSize = 20
	}
	opt.SetMaxPoolSize(conf.MaxPoolSize)
	if conf.MinPoolSize == 0 {
		conf.MinPoolSize = 1
	}
	opt.SetMinPoolSize(conf.MinPoolSize)

	return qmgo.NewClient(ctx, &qmgo.Config{Uri: conf.DSN})

}
