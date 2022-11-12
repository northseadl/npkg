package kratos

import (
	"context"
	"errors"
	kerrs "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/northseadl/go-utils/slices"
	"strings"
)

// Credential 权限中间件的 context value 的 key
type Credential struct{}

// AuthCredential AuthMeta
type AuthCredential struct {
	Uid   string
	Perms []string
}

// AuthPermsMatchOptions operation => perms[]
type AuthPermsMatchOptions map[string][]string

type AuthPermMatchBuilder interface {
	Match(opt string, perms ...string) AuthPermMatchBuilder
	Build() AuthPermsMatchOptions
}

type authPermMatchBuilder struct {
	hashMap map[string][]string
}

func (b *authPermMatchBuilder) Match(opt string, perms ...string) AuthPermMatchBuilder {
	if len(perms) > 0 {
		b.hashMap[opt] = perms
	}
	return b
}

func (b *authPermMatchBuilder) Build() AuthPermsMatchOptions {
	return b.hashMap
}

func NewAuthPermMatchBuilder() AuthPermMatchBuilder {
	return &authPermMatchBuilder{
		hashMap: make(AuthPermsMatchOptions),
	}
}

// AuthPermsMiddleware Metadata权限匹配器
func AuthPermsMiddleware(opts AuthPermsMatchOptions) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				uid := tr.RequestHeader().Get("x-auth-uid")
				perms := strings.Split(tr.RequestHeader().Get("x-auth-perms"), ",")

				op := tr.Operation()
				if opPerms, ok2 := opts[op]; ok2 {
					if !slices.Contains(perms, opPerms...) {
						return nil, kerrs.Unauthorized("PERMISSION_DENIED", "权限不足")
					}
				}

				if uid != "" {
					ctx = context.WithValue(ctx, Credential{}, AuthCredential{
						Uid:   uid,
						Perms: perms,
					})
				}
			}
			return handler(ctx, req)
		}
	}
}

var ErrMissAuthCredential = errors.New("miss auth credential")

func FromAuthContext(ctx context.Context) (*AuthCredential, error) {
	credential, ok := ctx.Value(Credential{}).(AuthCredential)
	if !ok {
		return nil, ErrMissAuthCredential
	}
	return &credential, nil
}

func WrapErrForKratos(err error) error {
	return kerrs.Unauthorized("UNAUTHORIZED", err.Error())
}
