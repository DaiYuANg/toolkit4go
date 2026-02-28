// Package adapteroptions 提供适配器配置选项
package adapteroptions

import (
	"log/slog"

	"github.com/DaiYuANg/toolkit4go/httpx"
)

// StdHTTPOptions StdHTTPAdapter 配置
type StdHTTPOptions struct {
	Logger     *slog.Logger
	Huma       httpx.HumaOptions
	EnableHuma bool
}

// DefaultStdHTTPOptions 默认配置
func DefaultStdHTTPOptions() *StdHTTPOptions {
	return &StdHTTPOptions{
		Logger: slog.Default(),
	}
}

// StdHTTPOption 配置选项函数
type StdHTTPOption func(*StdHTTPOptions)

// WithLogger 设置日志记录器
func WithStdLogger(logger *slog.Logger) StdHTTPOption {
	return func(o *StdHTTPOptions) {
		o.Logger = logger
	}
}

// WithHuma 启用 Huma OpenAPI
func WithStdHuma(title, version, description string) StdHTTPOption {
	return func(o *StdHTTPOptions) {
		o.EnableHuma = true
		o.Huma = httpx.HumaOptions{
			Enabled:     true,
			Title:       title,
			Version:     version,
			Description: description,
		}
	}
}

// Build 构建适配器
func (o *StdHTTPOptions) Build() *httpx.StdHTTPAdapter {
	adapter := httpx.NewStdHTTPAdapter().WithLogger(o.Logger)
	if o.EnableHuma {
		adapter.WithHuma(o.Huma)
	}
	return adapter
}

// GinOptions GinAdapter 配置
type GinOptions struct {
	Logger     *slog.Logger
	Huma       httpx.HumaOptions
	EnableHuma bool
	Mode       string // debug, release, test
}

// DefaultGinOptions 默认配置
func DefaultGinOptions() *GinOptions {
	return &GinOptions{
		Logger: slog.Default(),
		Mode:   "release",
	}
}

// GinOption 配置选项函数
type GinOption func(*GinOptions)

// WithGinLogger 设置日志记录器
func WithGinLogger(logger *slog.Logger) GinOption {
	return func(o *GinOptions) {
		o.Logger = logger
	}
}

// WithGinHuma 启用 Huma OpenAPI
func WithGinHuma(title, version, description string) GinOption {
	return func(o *GinOptions) {
		o.EnableHuma = true
		o.Huma = httpx.HumaOptions{
			Enabled:     true,
			Title:       title,
			Version:     version,
			Description: description,
		}
	}
}

// WithMode 设置 Gin 模式
func WithMode(mode string) GinOption {
	return func(o *GinOptions) {
		o.Mode = mode
	}
}

// Build 构建适配器
func (o *GinOptions) Build() *httpx.GinAdapter {
	adapter := httpx.NewGinAdapter().WithLogger(o.Logger)
	if o.EnableHuma {
		adapter.WithHuma(o.Huma)
	}
	return adapter
}

// EchoOptions EchoAdapter 配置
type EchoOptions struct {
	Logger     *slog.Logger
	Huma       httpx.HumaOptions
	EnableHuma bool
	HideBanner bool
}

// DefaultEchoOptions 默认配置
func DefaultEchoOptions() *EchoOptions {
	return &EchoOptions{
		Logger:     slog.Default(),
		HideBanner: true,
	}
}

// EchoOption 配置选项函数
type EchoOption func(*EchoOptions)

// WithEchoLogger 设置日志记录器
func WithEchoLogger(logger *slog.Logger) EchoOption {
	return func(o *EchoOptions) {
		o.Logger = logger
	}
}

// WithEchoHuma 启用 Huma OpenAPI
func WithEchoHuma(title, version, description string) EchoOption {
	return func(o *EchoOptions) {
		o.EnableHuma = true
		o.Huma = httpx.HumaOptions{
			Enabled:     true,
			Title:       title,
			Version:     version,
			Description: description,
		}
	}
}

// WithHideBanner 设置是否隐藏 Banner
func WithHideBanner(enabled bool) EchoOption {
	return func(o *EchoOptions) {
		o.HideBanner = enabled
	}
}

// Build 构建适配器
func (o *EchoOptions) Build() *httpx.EchoAdapter {
	adapter := httpx.NewEchoAdapter().WithLogger(o.Logger)
	if o.EnableHuma {
		adapter.WithHuma(o.Huma)
	}
	return adapter
}

// FiberOptions FiberAdapter 配置
type FiberOptions struct {
	Logger       *slog.Logger
	Huma         httpx.HumaOptions
	EnableHuma   bool
	Prefork      bool
	ServerHeader string
}

// DefaultFiberOptions 默认配置
func DefaultFiberOptions() *FiberOptions {
	return &FiberOptions{
		Logger:       slog.Default(),
		Prefork:      false,
		ServerHeader: "httpx-fiber",
	}
}

// FiberOption 配置选项函数
type FiberOption func(*FiberOptions)

// WithFiberLogger 设置日志记录器
func WithFiberLogger(logger *slog.Logger) FiberOption {
	return func(o *FiberOptions) {
		o.Logger = logger
	}
}

// WithFiberHuma 启用 Huma OpenAPI
func WithFiberHuma(title, version, description string) FiberOption {
	return func(o *FiberOptions) {
		o.EnableHuma = true
		o.Huma = httpx.HumaOptions{
			Enabled:     true,
			Title:       title,
			Version:     version,
			Description: description,
		}
	}
}

// WithPrefork 设置是否启用 Prefork 模式
func WithPrefork(enabled bool) FiberOption {
	return func(o *FiberOptions) {
		o.Prefork = enabled
	}
}

// WithServerHeader 设置 Server Header
func WithServerHeader(header string) FiberOption {
	return func(o *FiberOptions) {
		o.ServerHeader = header
	}
}

// Build 构建适配器
func (o *FiberOptions) Build() *httpx.FiberAdapter {
	adapter := httpx.NewFiberAdapter().WithLogger(o.Logger)
	if o.EnableHuma {
		adapter.WithHuma(o.Huma)
	}
	return adapter
}
