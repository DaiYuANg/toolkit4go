//go:build !no_fiber

package httpx

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/DaiYuANg/toolkit4go/httpx/huma"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

// FiberAdapter Fiber v2 框架适配器
type FiberAdapter struct {
	app        *fiber.App
	group      fiber.Router
	middleware []MiddlewareFunc
	logger     *slog.Logger
	huma       *huma.Service
	humaOpts   HumaOptions
}

// NewFiberAdapter 创建 Fiber 适配器
func NewFiberAdapter(app ...*fiber.App) *FiberAdapter {
	var a *fiber.App
	if len(app) > 0 {
		a = app[0]
	} else {
		a = fiber.New()
	}

	adapter := &FiberAdapter{
		app:    a,
		group:  a,
		logger: slog.Default(),
	}

	Register("fiber", func() Adapter {
		return NewFiberAdapter()
	})

	return adapter
}

// WithHuma 启用 Huma OpenAPI 文档
func (a *FiberAdapter) WithHuma(opts HumaOptions) *FiberAdapter {
	a.humaOpts = opts

	api := humafiber.New(a.app, huma.DefaultConfig(opts.Title, opts.Version))
	a.huma = huma.NewService(api, opts.Title, opts.Version, opts.Description)

	// Fiber 需要直接注册路由到 app
	a.registerHumaDocs()

	return a
}

// registerHumaDocs 注册 Huma 文档路由到 Fiber
func (a *FiberAdapter) registerHumaDocs() {
	if a.huma == nil {
		return
	}

	// OpenAPI JSON
	a.app.Get("/openapi.json", func(c *fiber.Ctx) error {
		c.Type("json")
		return c.JSON(a.huma.API().OpenAPI())
	})

	// Swagger UI
	a.app.Get("/docs", func(c *fiber.Ctx) error {
		c.Type("html")
		return c.SendString(a.swaggerUIHTML())
	})
}

// swaggerUIHTML 生成 Swagger UI HTML
func (a *FiberAdapter) swaggerUIHTML() string {
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>` + a.humaOpts.Title + `</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
        SwaggerUIBundle({url: "/openapi.json", dom_id: '#swagger-ui'});
    </script>
</body>
</html>`
}

// WithLogger 设置日志记录器
func (a *FiberAdapter) WithLogger(logger *slog.Logger) *FiberAdapter {
	a.logger = logger
	return a
}

// Name 返回适配器名称
func (a *FiberAdapter) Name() string {
	return "fiber"
}

// Handle 注册处理函数
func (a *FiberAdapter) Handle(method, path string, handler HandlerFunc) {
	a.group.Add(method, path, a.wrapHandler(handler))
}

// Group 创建路由组
func (a *FiberAdapter) Group(prefix string) Adapter {
	fiberGroup := a.group.Group(prefix)
	lo.ForEach(a.middleware, func(mw MiddlewareFunc, _ int) {
		fiberGroup.Use(a.fiberMiddleware(mw))
	})

	return &FiberAdapter{
		app:    a.app,
		group:  fiberGroup,
		logger: a.logger,
		huma:   a.huma,
	}
}

// Use 注册中间件
func (a *FiberAdapter) Use(middlewares ...MiddlewareFunc) {
	a.middleware = append(a.middleware, middlewares...)
	lo.ForEach(middlewares, func(mw MiddlewareFunc, _ int) {
		a.group.Use(a.fiberMiddleware(mw))
	})
}

// ServeHTTP 实现 http.Handler 接口
func (a *FiberAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

// App 返回 Fiber 应用
func (a *FiberAdapter) App() *fiber.App {
	return a.app
}

// wrapHandler 包装处理函数
func (a *FiberAdapter) wrapHandler(handler HandlerFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		w := &fiberResponseWriter{ctx: c}
		r := convertFiberRequest(c)

		if err := handler(c.Context(), w, r); err != nil {
			a.logger.Error("Handler error",
				slog.String("method", c.Method()),
				slog.String("path", c.Path()),
				slog.String("error", err.Error()),
			)
			return err
		}
		return nil
	}
}

// convertFiberRequest 转换 Fiber 请求
func convertFiberRequest(c *fiber.Ctx) *http.Request {
	u := &url.URL{
		Path:     c.Path(),
		RawQuery: string(c.Request().URI().QueryString()),
	}

	header := make(http.Header)
	c.Request().Header.VisitAll(func(k, v []byte) {
		header.Add(string(k), string(v))
	})

	return &http.Request{
		Method:        c.Method(),
		URL:           u,
		Header:        header,
		Body:          io.NopCloser(bytes.NewReader(c.Body())),
		ContentLength: int64(len(c.Body())),
		Host:          string(c.Request().Header.Host()),
		RemoteAddr:    c.IP(),
	}
}

// fiberResponseWriter 适配 Fiber 响应
type fiberResponseWriter struct {
	ctx        *fiber.Ctx
	statusCode int
	header     http.Header
}

func (w *fiberResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *fiberResponseWriter) Write(b []byte) (int, error) {
	w.ctx.Response().SetBody(b)
	return len(b), nil
}

func (w *fiberResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ctx.Status(statusCode)
}

// fiberMiddleware 将 httpx 中间件转换为 Fiber 中间件
func (a *FiberAdapter) fiberMiddleware(mw MiddlewareFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		next := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			return c.Next()
		}
		w := &fiberResponseWriter{ctx: c, header: make(http.Header)}
		r := convertFiberRequest(c)
		handler := mw(next)
		_ = handler(c.Context(), w, r)
		return nil
	}
}

// HumaService 返回 Huma 服务
func (a *FiberAdapter) HumaService() *huma.Service {
	return a.huma
}

// HasHuma 检查是否启用了 Huma
func (a *FiberAdapter) HasHuma() bool {
	return a.huma != nil
}

// RegisterHumaRoute 注册路由到 Huma
func (a *FiberAdapter) RegisterHumaRoute(method, path, operationID string) {
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
