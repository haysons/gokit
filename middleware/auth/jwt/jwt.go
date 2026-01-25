package jwt

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/haysons/gokit/errors"
	"github.com/haysons/gokit/middleware"
	"github.com/haysons/gokit/transport"
)

type authKey struct{}

const (
	bearerWord string = "Bearer"

	bearerFormat string = "Bearer %s"

	authorizationKey string = "Authorization"
)

var (
	ErrMissingJwtToken        = errors.NewUnauthorized(10001, "获取 jwt header 失败", "JWT token is missing")
	ErrMissingKeyFunc         = errors.NewUnauthorized(10002, "秘钥获取函数缺失", "keyFunc is missing")
	ErrTokenInvalid           = errors.NewUnauthorized(10003, "jwt 验证失败", "Token is invalid")
	ErrTokenExpired           = errors.NewUnauthorized(10004, "jwt 已过期", "JWT token has expired")
	ErrTokenParseFail         = errors.NewUnauthorized(10005, "jwt 解析失败", "Fail to parse JWT token ")
	ErrUnSupportSigningMethod = errors.NewUnauthorized(10006, "jwt 签名方式异常", "Wrong signing method")
	ErrWrongContext           = errors.NewUnauthorized(10007, "middleware context 异常", "Wrong context for middleware")
	ErrNeedTokenProvider      = errors.NewUnauthorized(10008, "获取 token 秘钥失败", "Token provider is missing")
	ErrSignToken              = errors.NewUnauthorized(10009, "生成 token 失败", "Can not sign token.Is the key correct?")
	ErrGetKey                 = errors.NewUnauthorized(10010, "获取 token 秘钥失败", "Can not get key while signing token")
)

type Option func(*options)

type options struct {
	signingMethod jwt.SigningMethod
	claims        func() jwt.Claims
	tokenHeader   map[string]any
}

// WithSigningMethod 指定 jwt 的签名方法
func WithSigningMethod(method jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

// WithClaims 指定 jwt 自定义 claim
func WithClaims(f func() jwt.Claims) Option {
	return func(o *options) {
		o.claims = f
	}
}

// WithTokenHeader 指定 jwt 自定义 tokenHeader
func WithTokenHeader(header map[string]any) Option {
	return func(o *options) {
		o.tokenHeader = header
	}
}

// Server 自 server 端 header 中解析并验证 jwt
func Server(keyFunc jwt.Keyfunc, opts ...Option) middleware.Middleware {
	o := &options{
		signingMethod: jwt.SigningMethodHS256,
	}
	for _, opt := range opts {
		opt(o)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				if keyFunc == nil {
					return nil, ErrMissingKeyFunc
				}
				auths := strings.SplitN(tr.RequestHeader().Get(authorizationKey), " ", 2)
				if len(auths) != 2 || !strings.EqualFold(auths[0], bearerWord) {
					return nil, ErrMissingJwtToken
				}
				jwtToken := auths[1]
				var (
					tokenInfo *jwt.Token
					err       error
				)
				if o.claims != nil {
					tokenInfo, err = jwt.ParseWithClaims(jwtToken, o.claims(), keyFunc)
				} else {
					tokenInfo, err = jwt.Parse(jwtToken, keyFunc)
				}
				if err != nil {
					if errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrTokenUnverifiable) {
						return nil, ErrTokenInvalid
					}
					if errors.Is(err, jwt.ErrTokenNotValidYet) || errors.Is(err, jwt.ErrTokenExpired) {
						return nil, ErrTokenExpired
					}
					return nil, ErrTokenParseFail
				}

				if !tokenInfo.Valid {
					return nil, ErrTokenInvalid
				}
				if tokenInfo.Method != o.signingMethod {
					return nil, ErrUnSupportSigningMethod
				}
				// 将 token 中包含的自定义数据放入 context 之中
				ctx = InjectContext(ctx, tokenInfo.Claims)
				return handler(ctx, req)
			}
			return nil, ErrWrongContext
		}
	}
}

// Client 自 client 端 header 中注入 jwt
func Client(keyProvider jwt.Keyfunc, opts ...Option) middleware.Middleware {
	claims := jwt.RegisteredClaims{}
	o := &options{
		signingMethod: jwt.SigningMethodHS256,
		claims:        func() jwt.Claims { return claims },
	}
	for _, opt := range opts {
		opt(o)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			if keyProvider == nil {
				return nil, ErrNeedTokenProvider
			}
			token := jwt.NewWithClaims(o.signingMethod, o.claims())
			if o.tokenHeader != nil {
				for k, v := range o.tokenHeader {
					token.Header[k] = v
				}
			}
			key, err := keyProvider(token)
			if err != nil {
				return nil, ErrGetKey
			}
			tokenStr, err := token.SignedString(key)
			if err != nil {
				return nil, ErrSignToken
			}
			// 自 client header 中添加 token 信息
			if clientContext, ok := transport.FromClientContext(ctx); ok {
				clientContext.RequestHeader().Set(authorizationKey, fmt.Sprintf(bearerFormat, tokenStr))
				return handler(ctx, req)
			}
			return nil, ErrWrongContext
		}
	}
}

// InjectContext 将 jwt 自定义信息注入 context 之中
func InjectContext(ctx context.Context, info jwt.Claims) context.Context {
	return context.WithValue(ctx, authKey{}, info)
}

// FromContext 自 context 中提取 jwt 自定义信息
func FromContext(ctx context.Context) (token jwt.Claims, ok bool) {
	token, ok = ctx.Value(authKey{}).(jwt.Claims)
	return
}
