package httpx

import (
	"context"
	"net/http"

	"github.com/samber/lo"
)

// HTTPMethod HTTP 方法常量
const (
	MethodGet     = http.MethodGet
	MethodPost    = http.MethodPost
	MethodPut     = http.MethodPut
	MethodDelete  = http.MethodDelete
	MethodPatch   = http.MethodPatch
	MethodHead    = http.MethodHead
	MethodOptions = http.MethodOptions
)

// RouteInfo 路由信息
type RouteInfo struct {
	Method      string   `json:"method"`
	Path        string   `json:"path"`
	HandlerName string   `json:"handler_name"`
	Comment     string   `json:"comment,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// String 返回路由的字符串表示
func (r RouteInfo) String() string {
	return r.Method + " " + r.Path + " -> " + r.HandlerName
}

// Endpoint HTTP 端点接口
type Endpoint interface {
	Routes() []RouteInfo
}

// HandlerFunc 通用处理函数签名
type HandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// Adapter HTTP 框架适配器接口
type Adapter interface {
	Name() string
	Handle(method, path string, handler HandlerFunc)
	Group(prefix string) Adapter
	Use(middlewares ...MiddlewareFunc)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// MiddlewareFunc 中间件函数签名
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

// AdapterFactory 适配器工厂函数
type AdapterFactory func() Adapter

// Registry 全局适配器注册表
var Registry = make(map[string]AdapterFactory)

// Register 注册适配器工厂
func Register(name string, factory AdapterFactory) {
	Registry[name] = factory
}

// Create 创建适配器实例
func Create(name string) (Adapter, error) {
	factory, ok := Registry[name]
	if !ok {
		return nil, ErrAdapterNotFound
	}
	return factory(), nil
}

// RegisteredAdapters 返回所有已注册的适配器名称
func RegisteredAdapters() []string {
	return lo.Keys(Registry)
}

// HumaOptions Huma OpenAPI 配置选项
type HumaOptions struct {
	// Enabled 是否启用 OpenAPI 文档
	Enabled bool
	// Title API 标题
	Title string
	// Version API 版本
	Version string
	// Description API 描述
	Description string
	// DocsPath 文档路径（默认 /docs）
	DocsPath string
	// OpenAPIPath OpenAPI JSON 路径（默认 /openapi.json）
	OpenAPIPath string
}

// DefaultHumaOptions 默认 Huma 配置
func DefaultHumaOptions() HumaOptions {
	return HumaOptions{
		Enabled:     false,
		Title:       "My API",
		Version:     "1.0.0",
		Description: "API Documentation",
		DocsPath:    "/docs",
		OpenAPIPath: "/openapi.json",
	}
}
