// Package huma 提供 Huma OpenAPI 文档集成功能
package huma

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// DefaultConfig 创建默认 Huma 配置
func DefaultConfig(title, version string) huma.Config {
	return huma.DefaultConfig(title, version)
}

// Register 注册路由到 Huma（泛型函数）
func Register[I, O any](api huma.API, method, path, operationID string, handler func(context.Context, *I) (*O, error)) {
	huma.Register(api, huma.Operation{
		OperationID: operationID,
		Method:      method,
		Path:        path,
		Summary:     operationID,
		Tags:        []string{"httpx"},
	}, handler)
}

// Service Huma OpenAPI 服务
type Service struct {
	api    huma.API
	config huma.Config
}

// NewService 创建 Huma 服务
func NewService(api huma.API, title, version, description string) *Service {
	config := huma.DefaultConfig(title, version)
	config.OpenAPI.Info.Description = description

	return &Service{
		api:    api,
		config: config,
	}
}

// API 返回 Huma API 实例
func (s *Service) API() huma.API {
	return s.api
}

// RegisterHandler 注册 OpenAPI 文档路由
func (s *Service) RegisterHandler(mux *http.ServeMux, docsPath, openAPIPath string) {
	mux.HandleFunc(openAPIPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.api.OpenAPI())
	})

	mux.HandleFunc(docsPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(s.swaggerUIHTML(openAPIPath)))
	})
}

// swaggerUIHTML 生成 Swagger UI HTML
func (s *Service) swaggerUIHTML(openAPIPath string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>%s</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
        SwaggerUIBundle({url: "%s", dom_id: '#swagger-ui'});
    </script>
</body>
</html>`, s.config.OpenAPI.Info.Title, openAPIPath)
}
