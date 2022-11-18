package kratos

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
)

var ErrUnknownRequest = errors.InternalServer("UNKNOWN", "未知请求错误")

func RecoveryZeroMiddleware(logger log.Logger) func(handler middleware.Handler) middleware.Handler {
	h := log.NewHelper(logger)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			defer func() {
				if rerr := recover(); rerr != nil {
					h.WithContext(ctx).Errorw("req", req, "err", rerr)
					err = ErrUnknownRequest
				}
			}()
			return handler(ctx, req)
		}
	}
}
