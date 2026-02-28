//go:build !no_echo

package httpx

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/DaiYuANg/toolkit4go/httpx/huma"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

// EchoAdapter Echo 框架适配器
type EchoAdapter struct {
	engine     *echo.Echo
	group      *echo.Group
	middleware []MiddlewareFunc
	logger     *slog.Logger
	huma       *huma.Service
}

// NewEchoAdapter 创建 Echo 适配器
func NewEchoAdapter(engine ...*echo.Echo) *EchoAdapter {
	var eng *echo.Echo
	if len(engine) > 0 {
		eng = engine[0]
	} else {
		eng = echo.New()
	}

	adapter := &EchoAdapter{
		engine: eng,
		group:  nil,
		logger: slog.Default(),
	}

	Register("echo", func() Adapter {
		return NewEchoAdapter()
	})

	return adapter
}

// WithHuma 启用 Huma OpenAPI 文档
func (a *EchoAdapter) WithHuma(opts HumaOptions) *EchoAdapter {
	api := humaecho.New(a.engine, huma.DefaultConfig(opts.Title, opts.Version))
	a.huma = huma.NewService(api, opts.Title, opts.Version, opts.Description)
	return a
}

// WithLogger 设置日志记录器
func (a *EchoAdapter) WithLogger(logger *slog.Logger) *EchoAdapter {
	a.logger = logger
	return a
}

// Name 返回适配器名称
func (a *EchoAdapter) Name() string {
	return "echo"
}

// Handle 注册处理函数
func (a *EchoAdapter) Handle(method, path string, handler HandlerFunc) {
	if a.group != nil {
		a.group.Add(method, path, a.echoHandler(handler))
	} else {
		a.engine.Add(method, path, a.echoHandler(handler))
	}
}

// Group 创建路由组
func (a *EchoAdapter) Group(prefix string) Adapter {
	var g *echo.Group
	if a.group != nil {
		g = a.group.Group(prefix)
	} else {
		g = a.engine.Group(prefix)
	}

	lo.ForEach(a.middleware, func(mw MiddlewareFunc, _ int) {
		g.Use(a.echoMiddleware(mw))
	})

	return &EchoAdapter{
		engine: a.engine,
		group:  g,
		logger: a.logger,
		huma:   a.huma,
	}
}

// Use 注册中间件
func (a *EchoAdapter) Use(middlewares ...MiddlewareFunc) {
	a.middleware = append(a.middleware, middlewares...)
	lo.ForEach(middlewares, func(mw MiddlewareFunc, _ int) {
		if a.group != nil {
			a.group.Use(a.echoMiddleware(mw))
		} else {
			a.engine.Use(a.echoMiddleware(mw))
		}
	})
}

// ServeHTTP 实现 http.Handler 接口
func (a *EchoAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.engine.ServeHTTP(w, r)
}

// Engine 返回 Echo 引擎
func (a *EchoAdapter) Engine() *echo.Echo {
	return a.engine
}

// echoHandler 包装处理函数
func (a *EchoAdapter) echoHandler(handler HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := handler(c.Request().Context(), c.Response(), c.Request()); err != nil {
			a.logger.Error("Handler error",
				slog.String("method", c.Request().Method),
				slog.String("path", c.Request().URL.Path),
				slog.String("error", err.Error()),
			)
			return err
		}
		return nil
	}
}

// echoMiddleware 将 httpx 中间件转换为 Echo 中间件
func (a *EchoAdapter) echoMiddleware(mw MiddlewareFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			wrappedNext := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
				return next(c)
			}
			handler := mw(wrappedNext)
			_ = handler(c.Request().Context(), c.Response(), c.Request())
			return nil
		}
	}
}

// HumaService 返回 Huma 服务
func (a *EchoAdapter) HumaService() *huma.Service {
	return a.huma
}

// HasHuma 检查是否启用了 Huma
func (a *EchoAdapter) HasHuma() bool {
	return a.huma != nil
}

// RegisterHumaRoute 注册路由到 Huma
func (a *EchoAdapter) RegisterHumaRoute(method, path, operationID string) {
	if a.huma == nil {
		return
	}

	huma.Register(a.huma.API(), method, path, operationID, func(ctx context.Context, input *struct{}) (*struct {
		Body struct {
			Operation string `json:"operation"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				Operation string `json:"operation"`
			}
		}{}
		resp.Body.Operation = operationID
		return resp, nil
	})
}
