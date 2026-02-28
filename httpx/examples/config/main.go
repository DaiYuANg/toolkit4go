package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/DaiYuANg/toolkit4go/httpx"
	"github.com/DaiYuANg/toolkit4go/httpx/options"
	"github.com/DaiYuANg/toolkit4go/httpx/options/adapteroptions"
	"github.com/DaiYuANg/toolkit4go/logx"
)

// UserEndpoint 用户端点
type UserEndpoint struct {
	httpx.BaseEndpoint
}

// ListUsers 获取用户列表
func (e *UserEndpoint) ListUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	e.Success(w, map[string]interface{}{
		"users": []string{"Alice", "Bob", "Charlie"},
	})
	return nil
}

// GetUser 获取单个用户
func (e *UserEndpoint) GetUser(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := e.GetQuery(r, "id", "1")
	e.Success(w, map[string]interface{}{
		"user": map[string]string{"id": id, "name": "User" + id},
	})
	return nil
}

func main() {
	// 创建 logger
	logger, _ := logx.New(logx.WithConsole(true))
	defer logger.Close()
	slogLogger := logx.NewSlog(logger)

	// 示例 1: 使用 ServerOptions 配置
	fmt.Println("=== Example 1: Using ServerOptions ===")
	serverOpts := options.DefaultServerOptions()
	serverOpts.Logger = slogLogger
	serverOpts.BasePath = "/api"
	serverOpts.PrintRoutes = true
	serverOpts.HumaEnabled = true
	serverOpts.HumaTitle = "My API"
	serverOpts.HumaVersion = "1.0.0"
	serverOpts.HumaDescription = "API Documentation"
	serverOpts.ReadTimeout = 15 * time.Second
	serverOpts.WriteTimeout = 15 * time.Second
	serverOpts.IdleTimeout = 60 * time.Second
	serverOpts.EnablePanicRecover = true
	serverOpts.EnableAccessLog = true
	serverOpts.Middlewares = append(serverOpts.Middlewares,
		// 添加自定义中间件
		func(next httpx.HandlerFunc) httpx.HandlerFunc {
			return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
				fmt.Println("Custom middleware:", r.Method, r.URL.Path)
				return next(ctx, w, r)
			}
		},
	)

	server1 := httpx.NewServer(serverOpts.Build()...)
	_ = server1.Register(&UserEndpoint{})

	// 示例 2: 使用 Adapter Options (Gin)
	fmt.Println("\n=== Example 2: Using Gin Adapter Options ===")
	ginAdapter := adapteroptions.DefaultGinOptions().Build()
	server2 := httpx.NewServer(
		httpx.WithAdapter(ginAdapter),
		httpx.WithLogger(slogLogger),
		httpx.WithPrintRoutes(true),
	)
	_ = server2.Register(&UserEndpoint{})

	// 示例 3: 使用 Adapter Options (Echo)
	fmt.Println("\n=== Example 3: Using Echo Adapter Options ===")
	echoAdapter := adapteroptions.DefaultEchoOptions().Build()
	server3 := httpx.NewServer(
		httpx.WithAdapter(echoAdapter),
		httpx.WithLogger(slogLogger),
		httpx.WithPrintRoutes(true),
	)
	_ = server3.Register(&UserEndpoint{})

	// 示例 4: 使用 Adapter Options (Fiber)
	fmt.Println("\n=== Example 4: Using Fiber Adapter Options ===")
	fiberAdapter := adapteroptions.DefaultFiberOptions().Build()
	server4 := httpx.NewServer(
		httpx.WithAdapter(fiberAdapter),
		httpx.WithLogger(slogLogger),
		httpx.WithPrintRoutes(true),
	)
	_ = server4.Register(&UserEndpoint{})

	// 示例 5: 使用 HTTP Client Options
	fmt.Println("\n=== Example 5: Using HTTP Client Options ===")
	clientOpts := &options.HTTPClientOptions{
		Timeout: 30 * time.Second,
	}
	client := clientOpts.Build()
	fmt.Printf("HTTP Client Timeout: %v\n", client.Timeout)

	// 示例 6: 使用 Context Options
	fmt.Println("\n=== Example 6: Using Context Options ===")
	ctxOpts := &options.ContextOptions{
		Timeout:       5 * time.Second,
		CancelOnPanic: true,
	}
	ctxOpts = options.WithContextValueOpt(ctxOpts, "request_id", "12345")
	ctx, cancel := ctxOpts.Build()
	defer cancel()
	fmt.Printf("Context created with timeout: 5s\n")
	fmt.Printf("Context value request_id: %v\n", ctx.Value("request_id"))

	fmt.Println("\n=== All Examples Complete ===")
	fmt.Println("Note: Servers are not started in this example.")
	fmt.Println("Use server.ListenAndServe() to start the server.")
}
