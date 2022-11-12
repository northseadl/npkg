package kratos

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

type ZeroConf struct {
	// Stdout 是否输出到控制台, 会影响性能, 建议在生产环境中关闭
	Stdout bool

	// Level 日志输出级别
	Level ZeroLevel

	// FileName 文件路径, 默认 log/trace.log
	FileName string

	// MaxSize 单个日志文件最大尺寸, 单位 MB, 默认为 100
	MaxSize int

	// MaxAge 保留旧日志(根据文件时间戳)的最大时间, 单位 days, 默认不删除旧日志
	MaxAge int

	// MaxBackups 保留的旧日志文件的最大数量, 默认不删除旧日志
	MaxBackups int

	// LocalTime 是否使用本地时间为轮转日志添加时间戳, 默认使用UTC时间
	RotateLocalTime bool

	// 是否启用压缩, 启用则旧日志会被压缩为gzip, 默认不压缩
	RotateCompress bool
}

type ZeroLogger struct {
	logger *zerolog.Logger
	sync   func() error
}

type ZeroLevel int8

const (
	ZeroDebugLevel = ZeroLevel(zerolog.DebugLevel)
	ZeroInfoLevel  = ZeroLevel(zerolog.InfoLevel)
	ZeroWarnLevel  = ZeroLevel(zerolog.WarnLevel)
	ZeroErrorLevel = ZeroLevel(zerolog.ErrorLevel)
	ZeroFatalLevel = ZeroLevel(zerolog.FatalLevel)
	ZeroPanicLevel = ZeroLevel(zerolog.PanicLevel)
	ZeroNoLevel    = ZeroLevel(zerolog.NoLevel)
	ZeroDisabled   = ZeroLevel(zerolog.Disabled)

	ZeroTraceLevel = ZeroLevel(zerolog.TraceLevel)
)

type loggerWrapper struct {
	writer func(p []byte) (n int, err error)
	closer func() error
	err    error
	n      int
	std    io.Writer
	logger io.Writer
}

func (l *loggerWrapper) Close() error {
	return l.closer()
}

func (l *loggerWrapper) Write(p []byte) (n int, err error) {
	return l.writer(p)
}

var ErrWriteFailed = errors.New("logger write failed")
var ErrInvalidErrValue = errors.New("err key has invalid err value")

// NewZeroLogger 创建 zero logger 的 kratos logger
func NewZeroLogger(conf *ZeroConf) *ZeroLogger {
	// stage: create rotate logger
	lumberjackLogger := lumberjack.Logger{
		Filename:   conf.FileName,
		MaxSize:    conf.MaxSize,
		MaxAge:     conf.MaxAge,
		MaxBackups: conf.MaxBackups,
		LocalTime:  conf.RotateLocalTime,
		Compress:   conf.RotateCompress,
	}

	// stage: create zero logger
	var out loggerWrapper
	out.logger = &lumberjackLogger
	out.closer = func() error {
		return out.logger.(io.Closer).Close()
	}
	if conf.Stdout {
		// stdout && logger
		out.std = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
		out.writer = func(p []byte) (n int, err error) {
			out.n, out.err = out.std.Write(p)
			if out.err != nil {
				return 0, errors.Wrap(err, ErrWriteFailed.Error())
			}
			out.n, out.err = out.logger.Write(p)
			if out.err != nil {
				out.err = errors.Wrap(err, ErrWriteFailed.Error())
			}
			return out.n, err
		}
	} else {
		// logger
		out.writer = func(p []byte) (n int, err error) {
			out.n, out.err = out.logger.Write(p)
			if out.err != nil {
				out.err = errors.Wrap(err, ErrWriteFailed.Error())
			}
			return out.n, err
		}
	}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	// timestamp & stack 会损耗一定的性能, 但是这是有必要的
	zeroLogger := zerolog.New(&out).With().Stack().Timestamp().Logger().Level(zerolog.Level(conf.Level))

	// wrap to kratos logger
	return &ZeroLogger{
		logger: &zeroLogger,
		sync: func() error {
			return out.Close()
		},
	}
}

const (
	errKey = "err"
)

func (l *ZeroLogger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.logger.Warn().Msg(fmt.Sprint("Key-values must appear in pairs: ", keyvals))
		return nil
	}

	var e *zerolog.Event
	switch level {
	case log.LevelDebug:
		e = l.logger.Debug()
		break
	case log.LevelInfo:
		e = l.logger.Info()
		break
	case log.LevelWarn:
		e = l.logger.Warn()
		break
	case log.LevelError:
		e = l.logger.Error()
		break
	case log.LevelFatal:
		e = l.logger.Fatal()
		break
	}
	for i := 0; i < len(keyvals); i += 2 {
		if keyvals[i] == errKey {
			err, ok := keyvals[i+1].(error)
			if ok {
				e = e.Err(err)
			} else {
				return ErrInvalidErrValue
			}
		} else {
			e = e.Interface(keyvals[i].(string), keyvals[i+1])
		}
	}
	e.Send()
	return nil
}

func (l *ZeroLogger) Sync() error {
	return l.sync()
}
