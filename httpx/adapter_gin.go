//go:build !no_gin

package httpx

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/DaiYuANg/toolkit4go/httpx/huma"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

// GinAdapter Gin 框架适配器
type GinAdapter struct {
	engine     *gin.Engine
	group      *gin.RouterGroup
	middleware []MiddlewareFunc
	logger     *slog.Logger
	huma       *huma.Service
}

// NewGinAdapter 创建 Gin 适配器
func NewGinAdapter(engine ...*gin.Engine) *GinAdapter {
	var eng *gin.Engine
	if len(engine) > 0 {
		eng = engine[0]
	} else {
		eng = gin.New()
	}

	adapter := &GinAdapter{
		engine: eng,
		group:  &eng.RouterGroup,
		logger: slog.Default(),
	}

	Register("gin", func() Adapter {
		return NewGinAdapter()
	})

	return adapter
}

// WithHuma 启用 Huma OpenAPI 文档
func (a *GinAdapter) WithHuma(opts HumaOptions) *GinAdapter {
	api := humagin.New(a.engine, huma.DefaultConfig(opts.Title, opts.Version))
	a.huma = huma.NewService(api, opts.Title, opts.Version, opts.Description)
	return a
}

// WithLogger 设置日志记录器
func (a *GinAdapter) WithLogger(logger *slog.Logger) *GinAdapter {
	a.logger = logger
	return a
}

// Name 返回适配器名称
func (a *GinAdapter) Name() string {
	return "gin"
}

// Handle 注册处理函数
func (a *GinAdapter) Handle(method, path string, handler HandlerFunc) {
	a.group.Handle(method, path, a.wrapHandler(handler))
}

// Group 创建路由组
func (a *GinAdapter) Group(prefix string) Adapter {
	ginGroup := a.group.Group(prefix)
	lo.ForEach(a.middleware, func(mw MiddlewareFunc, _ int) {
		ginGroup.Use(a.ginMiddleware(mw))
	})

	return &GinAdapter{
		engine: a.engine,
		group:  ginGroup,
		logger: a.logger,
		huma:   a.huma,
	}
}

// Use 注册中间件
func (a *GinAdapter) Use(middlewares ...MiddlewareFunc) {
	a.middleware = append(a.middleware, middlewares...)
	lo.ForEach(middlewares, func(mw MiddlewareFunc, _ int) {
		a.group.Use(a.ginMiddleware(mw))
	})
}

// ServeHTTP 实现 http.Handler 接口
func (a *GinAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.engine.ServeHTTP(w, r)
}

// Engine 返回 Gin 引擎
func (a *GinAdapter) Engine() *gin.Engine {
	return a.engine
}

// wrapHandler 包装处理函数
func (a *GinAdapter) wrapHandler(handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := handler(c.Request.Context(), c.Writer, c.Request); err != nil {
			a.logger.Error("Handler error",
				slog.String("method", c.Request.Method),
				slog.String("path", c.Request.URL.Path),
				slog.String("error", err.Error()),
			)
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
}

// ginMiddleware 将 httpx 中间件转换为 Gin 中间件
func (a *GinAdapter) ginMiddleware(mw MiddlewareFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		next := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			c.Next()
			return nil
		}
		wrapped := mw(next)
		_ = wrapped(c.Request.Context(), c.Writer, c.Request)
	}
}

// HumaService 返回 Huma 服务
func (a *GinAdapter) HumaService() *huma.Service {
	return a.huma
}

// HasHuma 检查是否启用了 Huma
func (a *GinAdapter) HasHuma() bool {
	return a.huma != nil
}

// RegisterHumaRoute 注册路由到 Huma
func (a *GinAdapter) RegisterHumaRoute(method, path, operationID string) {
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
