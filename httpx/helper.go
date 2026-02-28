package httpx

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/samber/lo"
)

// BaseEndpoint 基础 endpoint 结构，提供常用辅助方法
type BaseEndpoint struct{}

// JSON 返回 JSON 响应
func (BaseEndpoint) JSON(w http.ResponseWriter, data interface{}, status ...int) {
	code := lo.FirstOr(status, http.StatusOK)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

// Error 返回错误响应
func (BaseEndpoint) Error(w http.ResponseWriter, message string, code ...int) {
	status := lo.FirstOr(code, http.StatusInternalServerError)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// Success 返回成功响应
func (BaseEndpoint) Success(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   data,
	})
}

// GetHeader 获取请求头
func (BaseEndpoint) GetHeader(r *http.Request, key string) string {
	return r.Header.Get(key)
}

// GetQuery 获取查询参数
func (BaseEndpoint) GetQuery(r *http.Request, key string, defaultValue ...string) string {
	value := r.URL.Query().Get(key)
	return lo.Ternary(value == "" && len(defaultValue) > 0, defaultValue[0], value)
}

// GetQueryOrDefault 获取查询参数，带默认值
func (BaseEndpoint) GetQueryOrDefault(r *http.Request, key string, defaultValue string) string {
	value := r.URL.Query().Get(key)
	return lo.Ternary(value == "", defaultValue, value)
}

// HandlerEndpoint 实际的处理 endpoint 示例
type HandlerEndpoint struct {
	BaseEndpoint
}

// GetUserList 获取用户列表
// 通过方法名自动生成路由：GET /user/list
func (e *HandlerEndpoint) GetUserList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	e.Success(w, map[string]interface{}{
		"users": []string{"user1", "user2", "user3"},
	})
	return nil
}

// GetUserByID 获取单个用户
// 通过方法名自动生成路由：GET /user/by/id
func (e *HandlerEndpoint) GetUserByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := e.GetQuery(r, "id", "1")
	e.Success(w, map[string]interface{}{
		"user": "user" + id,
	})
	return nil
}

// CreateNewUser 创建用户
// 通过方法名自动生成路由：POST /new/user
func (e *HandlerEndpoint) CreateNewUser(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	e.Success(w, map[string]interface{}{
		"message": "user created",
	})
	return nil
}

// UpdateUserInfo 更新用户
// 通过方法名自动生成路由：PUT /user/info
func (e *HandlerEndpoint) UpdateUserInfo(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	e.Success(w, map[string]interface{}{
		"message": "user updated",
	})
	return nil
}

// DeleteUserByID 删除用户
// 通过方法名自动生成路由：DELETE /user/by/id
func (e *HandlerEndpoint) DeleteUserByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	e.Success(w, map[string]interface{}{
		"message": "user deleted",
	})
	return nil
}
