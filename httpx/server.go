package httpx

import (
	"context"
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"github.com/samber/lo"
)

// Server HTTP 服务器
type Server struct {
	adapter     Adapter
	generator   *RouterGenerator
	middleware  []MiddlewareFunc
	basePath    string
	routes      []RouteInfo
	logger      *slog.Logger
	printRoutes bool
	humaOpts    HumaOptions
}

// ServerOption Server 配置选项
type ServerOption func(*Server)

// WithAdapter 设置适配器
func WithAdapter(adapter Adapter) ServerOption {
	return func(s *Server) {
		s.adapter = adapter
	}
}

// WithAdapterName 通过名称设置适配器
func WithAdapterName(name string) ServerOption {
	return func(s *Server) {
		adapter, err := Create(name)
		if err == nil {
			s.adapter = adapter
		}
	}
}

// WithGenerator 设置路由生成器
func WithGenerator(gen *RouterGenerator) ServerOption {
	return func(s *Server) {
		s.generator = gen
	}
}

// WithBasePath 设置基础路径
func WithBasePath(path string) ServerOption {
	return func(s *Server) {
		s.basePath = path
		s.generator.opts.BasePath = path
	}
}

// WithMiddleware 注册中间件
func WithMiddleware(middlewares ...MiddlewareFunc) ServerOption {
	return func(s *Server) {
		s.middleware = append(s.middleware, middlewares...)
	}
}

// WithLogger 设置日志记录器
func WithLogger(logger *slog.Logger) ServerOption {
	return func(s *Server) {
		s.logger = logger
	}
}

// WithPrintRoutes 设置是否打印路由
func WithPrintRoutes(enabled bool) ServerOption {
	return func(s *Server) {
		s.printRoutes = enabled
	}
}

// WithHuma 配置 Huma OpenAPI 文档
func WithHuma(opts HumaOptions) ServerOption {
	return func(s *Server) {
		s.humaOpts = opts
	}
}

// NewServer 创建 HTTP 服务器
func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		generator: NewRouterGenerator(),
		logger:    slog.Default(),
		routes:    make([]RouteInfo, 0),
	}

	lo.ForEach(opts, func(opt ServerOption, index int) {
		opt(s)
	})

	if s.adapter == nil {
		s.adapter = NewStdHTTPAdapter()
	}

	// 如果启用了 Huma，配置适配器
	if s.humaOpts.Enabled {
		s.configureAdapterHuma()
	}

	if len(s.middleware) > 0 {
		s.adapter.Use(s.middleware...)
	}

	return s
}

// configureAdapterHuma 配置适配器的 Huma 支持
func (s *Server) configureAdapterHuma() {
	switch adapter := s.adapter.(type) {
	case *StdHTTPAdapter:
		adapter.WithHuma(s.humaOpts)
	case *GinAdapter:
		adapter.WithHuma(s.humaOpts)
	case *EchoAdapter:
		adapter.WithHuma(s.humaOpts)
	case *FiberAdapter:
		adapter.WithHuma(s.humaOpts)
	}
}

// Register 注册 endpoint
func (s *Server) Register(endpoints ...interface{}) error {
	for _, endpoint := range endpoints {
		result := s.generator.Generate(endpoint)
		routes, err := result.Get()
		if err != nil {
			return err
		}

		lo.ForEach(routes, func(route RouteInfo, _ int) {
			s.adapter.Handle(route.Method, route.Path, s.wrapHandler(endpoint, route))
			s.routes = append(s.routes, route)

			// 同步到 Huma OpenAPI
			s.registerHumaRoute(route)
		})
	}

	s.printRoutesIfEnabled()
	return nil
}

// registerHumaRoute 注册 Huma OpenAPI 路由
func (s *Server) registerHumaRoute(route RouteInfo) {
	operationID := strings.ToUpper(route.Method) + "-" + strings.ReplaceAll(strings.TrimPrefix(route.Path, "/"), "/", "-")

	// 调用适配器的 RegisterHumaRoute 方法
	switch adapter := s.adapter.(type) {
	case *StdHTTPAdapter:
		adapter.RegisterHumaRoute(route.Method, route.Path, operationID)
	case *GinAdapter:
		adapter.RegisterHumaRoute(route.Method, route.Path, operationID)
	case *EchoAdapter:
		adapter.RegisterHumaRoute(route.Method, route.Path, operationID)
	case *FiberAdapter:
		adapter.RegisterHumaRoute(route.Method, route.Path, operationID)
	}
}

// HumaInput Huma 输入结构
type HumaInput struct{}

// HumaOutput Huma 输出结构
type HumaOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

// RegisterWithPrefix 注册 endpoint 并添加路径前缀
func (s *Server) RegisterWithPrefix(prefix string, endpoints ...interface{}) error {
	opts := DefaultGeneratorOptions()
	opts.BasePath = s.basePath + prefix
	gen := NewRouterGenerator(opts)

	for _, endpoint := range endpoints {
		result := gen.Generate(endpoint)
		routes, err := result.Get()
		if err != nil {
			return err
		}

		lo.ForEach(routes, func(route RouteInfo, _ int) {
			s.adapter.Handle(route.Method, route.Path, s.wrapHandler(endpoint, route))
			s.routes = append(s.routes, route)
			s.registerHumaRoute(route)
		})
	}

	s.printRoutesIfEnabled()
	return nil
}

// printRoutesIfEnabled 打印路由
func (s *Server) printRoutesIfEnabled() {
	if !s.printRoutes {
		return
	}

	s.logger.Info("Registered routes", slog.Int("count", len(s.routes)))
	lo.ForEach(s.routes, func(route RouteInfo, _ int) {
		s.logger.Info("  "+route.String(),
			slog.String("method", route.Method),
			slog.String("path", route.Path),
			slog.String("handler", route.HandlerName),
		)
	})
}

// GetRoutes 返回所有路由
func (s *Server) GetRoutes() []RouteInfo {
	return lo.Map(s.routes, func(route RouteInfo, _ int) RouteInfo {
		return route
	})
}

// GetRoutesByMethod 按方法过滤路由
func (s *Server) GetRoutesByMethod(method string) []RouteInfo {
	return lo.Filter(s.routes, func(route RouteInfo, _ int) bool {
		return route.Method == method
	})
}

// GetRoutesByPath 按路径过滤路由
func (s *Server) GetRoutesByPath(prefix string) []RouteInfo {
	return lo.Filter(s.routes, func(route RouteInfo, _ int) bool {
		return len(prefix) == 0 || strings.HasPrefix(route.Path, prefix)
	})
}

// HasRoute 检查路由是否存在
func (s *Server) HasRoute(method, path string) bool {
	return lo.SomeBy(s.routes, func(route RouteInfo) bool {
		return route.Method == method && route.Path == path
	})
}

// RouteCount 返回路由数量
func (s *Server) RouteCount() int {
	return len(s.routes)
}

// wrapHandler 包装 handler
func (s *Server) wrapHandler(endpoint interface{}, route RouteInfo) HandlerFunc {
	v := reflect.ValueOf(endpoint)
	method := v.MethodByName(route.HandlerName)

	if !method.IsValid() {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			return NewError(http.StatusInternalServerError, "handler not found")
		}
	}

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		methodType := method.Type()
		args := make([]reflect.Value, methodType.NumIn())

		for i := 0; i < methodType.NumIn(); i++ {
			paramType := methodType.In(i)
			if paramType.Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
				args[i] = reflect.ValueOf(ctx)
			}
			if paramType.Implements(reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()) {
				args[i] = reflect.ValueOf(w)
			}
			if paramType == reflect.TypeOf(&http.Request{}) {
				args[i] = reflect.ValueOf(r)
			}
		}

		results := method.Call(args)
		if len(results) > 0 {
			if err, ok := results[0].Interface().(error); ok && err != nil {
				return err
			}
		}
		return nil
	}
}

// Handler 返回 http.Handler
func (s *Server) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.adapter.ServeHTTP(w, r)
	})
}

// ServeHTTP 实现 http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Handler().ServeHTTP(w, r)
}

// ListenAndServe 启动服务器
func (s *Server) ListenAndServe(addr string) error {
	s.logger.Info("Starting server",
		slog.String("address", addr),
		slog.String("adapter", s.adapter.Name()),
		slog.Int("routes", len(s.routes)),
	)

	// 如果启用了 Huma，需要组合路由（Fiber 除外）
	if s.HasHuma() {
		// Fiber 已经自己处理了 OpenAPI 路由
		if _, ok := s.adapter.(*FiberAdapter); ok {
			return s.startFiberServer(addr)
		}

		mux := http.NewServeMux()

		// 注册应用路由
		mux.Handle("/", s.adapter)

		// 注册 Huma OpenAPI 路由
		switch adapter := s.adapter.(type) {
		case *StdHTTPAdapter:
			if svc := adapter.HumaService(); svc != nil {
				svc.RegisterHandler(mux, "/docs", "/openapi.json")
			}
		case *GinAdapter:
			if svc := adapter.HumaService(); svc != nil {
				svc.RegisterHandler(mux, "/docs", "/openapi.json")
			}
		case *EchoAdapter:
			if svc := adapter.HumaService(); svc != nil {
				svc.RegisterHandler(mux, "/docs", "/openapi.json")
			}
		}

		return http.ListenAndServe(addr, mux)
	}

	return http.ListenAndServe(addr, s.Handler())
}

// startFiberServer 启动 Fiber 服务器
func (s *Server) startFiberServer(addr string) error {
	if adapter, ok := s.adapter.(*FiberAdapter); ok {
		return adapter.App().Listen(addr)
	}
	return nil
}

// ListenAndServeContext 启动服务器（支持 context）
func (s *Server) ListenAndServeContext(ctx context.Context, addr string) error {
	server := &http.Server{
		Addr:    addr,
		Handler: s.Handler(),
	}

	s.logger.Info("Starting server with context", slog.String("address", addr))

	go func() {
		<-ctx.Done()
		s.logger.Info("Shutting down server")
		server.Shutdown(context.Background())
	}()

	return server.ListenAndServe()
}

// Adapter 返回适配器
func (s *Server) Adapter() Adapter {
	return s.adapter
}

// Logger 返回日志记录器
func (s *Server) Logger() *slog.Logger {
	return s.logger
}

// HasHuma 检查是否启用了 Huma
func (s *Server) HasHuma() bool {
	switch adapter := s.adapter.(type) {
	case *StdHTTPAdapter:
		return adapter.HasHuma()
	case *GinAdapter:
		return adapter.HasHuma()
	case *EchoAdapter:
		return adapter.HasHuma()
	case *FiberAdapter:
		return adapter.HasHuma()
	default:
		return false
	}
}
