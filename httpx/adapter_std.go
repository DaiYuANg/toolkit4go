package httpx

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/DaiYuANg/toolkit4go/httpx/huma"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/samber/lo"
)

// StdHTTPAdapter 标准 net/http 库适配器
type StdHTTPAdapter struct {
	mux        *http.ServeMux
	router     *chi.Mux
	middleware []MiddlewareFunc
	routes     map[string]http.HandlerFunc
	mu         sync.RWMutex
	logger     *slog.Logger
	huma       *huma.Service
}

// NewStdHTTPAdapter 创建标准 HTTP 适配器
func NewStdHTTPAdapter() *StdHTTPAdapter {
	router := chi.NewMux()
	mux := http.NewServeMux()

	adapter := &StdHTTPAdapter{
		mux:    mux,
		router: router,
		routes: make(map[string]http.HandlerFunc),
		logger: slog.Default(),
	}

	Register("std", func() Adapter {
		return NewStdHTTPAdapter()
	})

	return adapter
}

// WithHuma 启用 Huma OpenAPI 文档
func (a *StdHTTPAdapter) WithHuma(opts HumaOptions) *StdHTTPAdapter {
	api := humachi.New(a.router, huma.DefaultConfig(opts.Title, opts.Version))
	a.huma = huma.NewService(api, opts.Title, opts.Version, opts.Description)
	return a
}

// WithLogger 设置日志记录器
func (a *StdHTTPAdapter) WithLogger(logger *slog.Logger) *StdHTTPAdapter {
	a.logger = logger
	return a
}

// Name 返回适配器名称
func (a *StdHTTPAdapter) Name() string {
	return "std"
}

// Handle 注册处理函数
func (a *StdHTTPAdapter) Handle(method, path string, handler HandlerFunc) {
	a.mu.Lock()
	defer a.mu.Unlock()

	wrappedHandler := a.wrapHandler(handler)
	key := method + " " + path
	a.routes[key] = wrappedHandler

	a.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		key := r.Method + " " + r.URL.Path
		a.mu.RLock()
		handler, ok := a.routes[key]
		a.mu.RUnlock()

		if !ok {
			a.mu.RLock()
			found := lo.SomeBy(lo.Keys(a.routes), func(k string) bool {
				return strings.HasSuffix(k, " "+path)
			})
			a.mu.RUnlock()

			if found {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		handler(w, r)
	})
}

// Group 创建路由组
func (a *StdHTTPAdapter) Group(prefix string) Adapter {
	return &StdHTTPAdapter{
		mux:        a.mux,
		router:     a.router,
		middleware: a.middleware,
		routes:     a.routes,
		mu:         a.mu,
		logger:     a.logger,
		huma:       a.huma,
	}
}

// Use 注册中间件
func (a *StdHTTPAdapter) Use(middlewares ...MiddlewareFunc) {
	a.middleware = append(a.middleware, middlewares...)
}

// ServeHTTP 实现 http.Handler 接口
func (a *StdHTTPAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

// wrapHandler 包装处理函数
func (a *StdHTTPAdapter) wrapHandler(handler HandlerFunc) http.HandlerFunc {
	h := handler
	for i := len(a.middleware) - 1; i >= 0; i-- {
		h = a.middleware[i](h)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(r.Context(), w, r); err != nil {
			a.logger.Error("Handler error",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("error", err.Error()),
			)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// HumaService 返回 Huma 服务
func (a *StdHTTPAdapter) HumaService() *huma.Service {
	return a.huma
}

// HasHuma 检查是否启用了 Huma
func (a *StdHTTPAdapter) HasHuma() bool {
	return a.huma != nil
}

// RegisterHumaRoute 注册路由到 Huma
func (a *StdHTTPAdapter) RegisterHumaRoute(method, path, operationID string) {
	if a.huma == nil {
		return
	}

	huma.Register(a.huma.API(), method, path, operationID, func(ctx context.Context, input *struct{}) (*struct {
		Body struct {
			Operation string `json:"operation"`
			Method    string `json:"method"`
			Path      string `json:"path"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				Operation string `json:"operation"`
				Method    string `json:"method"`
				Path      string `json:"path"`
			}
		}{}
		resp.Body.Operation = operationID
		resp.Body.Method = method
		resp.Body.Path = path
		return resp, nil
	})
}

// MiddlewareLogger 日志中间件
func MiddlewareLogger(next HandlerFunc) HandlerFunc {
	logger := slog.Default()
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		logger.Debug("Request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)
		return next(ctx, w, r)
	}
}

// MiddlewareRecovery 恢复中间件
func MiddlewareRecovery(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Default().Error("Panic recovered",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Any("panic", rec),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		return next(ctx, w, r)
	}
}

// MiddlewareCORS CORS 中间件
func MiddlewareCORS(allowedOrigins ...string) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			origin := r.Header.Get("Origin")
			allowed := lo.Ternary(len(allowedOrigins) == 0, true,
				lo.Contains(allowedOrigins, "*") || lo.Contains(allowedOrigins, origin),
			)

			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return nil
			}
			return next(ctx, w, r)
		}
	}
}
