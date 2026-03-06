package httpx

import "github.com/samber/lo"

// Endpoint 定义路由注册接口，用于组织类似 Controller 的代码
type Endpoint interface {
	RegisterRoutes(server *Server)
}

// BaseEndpoint 提供空实现的基类，用户可嵌入后只覆盖需要的方法
type BaseEndpoint struct{}

// RegisterRoutes 空实现，可被嵌入的结构体重写
func (e *BaseEndpoint) RegisterRoutes(server *Server) {}

type EndpointHookFunc func(server *Server, endpoint Endpoint)

type EndpointHooks struct {
	Before EndpointHookFunc
	After  EndpointHookFunc
}

func (s *Server) Register(endpoint Endpoint, hooks ...EndpointHooks) {
	if endpoint == nil {
		return
	}

	// Before hooks
	lo.ForEach(hooks, func(h EndpointHooks, index int) {
		if h.Before != nil {
			h.Before(s, endpoint)
		}
	})

	endpoint.RegisterRoutes(s)

	// After hooks
	lo.ForEach(hooks, func(h EndpointHooks, index int) {
		if h.After != nil {
			h.After(s, endpoint)
		}
	})
}

func (s *Server) RegisterOnly(endpoints ...Endpoint) {
	lo.ForEach(endpoints, func(e Endpoint, _ int) {
		if e == nil {
			if s.logger != nil {
				s.logger.Warn("skipping nil endpoint")
			}
			return
		}
		e.RegisterRoutes(s)
	})
}
