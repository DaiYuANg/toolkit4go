package httpx

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/samber/lo"
)

// Route registers related handlers.
func Route[I, O any](s *Server, method, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return registerTypedWithPrefix(s, "", method, path, handler, operationOptions...)
}

// Get registers related handlers.
func Get[I, O any](s *Server, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return Route(s, MethodGet, path, handler, operationOptions...)
}

// Post registers related handlers.
func Post[I, O any](s *Server, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return Route(s, MethodPost, path, handler, operationOptions...)
}

// Put registers related handlers.
func Put[I, O any](s *Server, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return Route(s, MethodPut, path, handler, operationOptions...)
}

// Patch registers related handlers.
func Patch[I, O any](s *Server, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return Route(s, MethodPatch, path, handler, operationOptions...)
}

// Delete registers related handlers.
func Delete[I, O any](s *Server, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return Route(s, MethodDelete, path, handler, operationOptions...)
}

// Head registers related handlers.
func Head[I, O any](s *Server, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return Route(s, MethodHead, path, handler, operationOptions...)
}

// Options registers related handlers.
func Options[I, O any](s *Server, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return Route(s, MethodOptions, path, handler, operationOptions...)
}

// GroupRoute registers related handlers.
func GroupRoute[I, O any](g *Group, method, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	if g == nil || g.server == nil {
		return fmt.Errorf("%w: route group is nil", ErrRouteNotRegistered)
	}
	return registerTypedWithPrefix(g.server, g.prefix, method, path, handler, operationOptions...)
}

// GroupGet registers related handlers.
func GroupGet[I, O any](g *Group, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return GroupRoute(g, MethodGet, path, handler, operationOptions...)
}

// GroupPost registers related handlers.
func GroupPost[I, O any](g *Group, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return GroupRoute(g, MethodPost, path, handler, operationOptions...)
}

// GroupPut registers related handlers.
func GroupPut[I, O any](g *Group, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return GroupRoute(g, MethodPut, path, handler, operationOptions...)
}

// GroupPatch registers related handlers.
func GroupPatch[I, O any](g *Group, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return GroupRoute(g, MethodPatch, path, handler, operationOptions...)
}

// GroupDelete registers related handlers.
func GroupDelete[I, O any](g *Group, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return GroupRoute(g, MethodDelete, path, handler, operationOptions...)
}

// MustRoute registers related handlers and panics on error.
func MustRoute[I, O any](s *Server, method, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	lo.Must0(Route(s, method, path, handler, operationOptions...))
}

// MustGet registers related handlers and panics on error.
func MustGet[I, O any](s *Server, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	lo.Must0(Get(s, path, handler, operationOptions...))
}

// MustPost registers related handlers and panics on error.
func MustPost[I, O any](s *Server, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	lo.Must0(Post(s, path, handler, operationOptions...))
}

// MustPut registers related handlers and panics on error.
func MustPut[I, O any](s *Server, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	lo.Must0(Put(s, path, handler, operationOptions...))
}

// MustPatch registers related handlers and panics on error.
func MustPatch[I, O any](s *Server, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	lo.Must0(Patch(s, path, handler, operationOptions...))
}

// MustDelete registers related handlers and panics on error.
func MustDelete[I, O any](s *Server, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	lo.Must0(Delete(s, path, handler, operationOptions...))
}

// MustHead registers related handlers and panics on error.
func MustHead[I, O any](s *Server, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	lo.Must0(Head(s, path, handler, operationOptions...))
}

// MustOptions registers related handlers and panics on error.
func MustOptions[I, O any](s *Server, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	lo.Must0(Options(s, path, handler, operationOptions...))
}

// MustGroupRoute registers related handlers and panics on error.
func MustGroupRoute[I, O any](g *Group, method, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	lo.Must0(GroupRoute(g, method, path, handler, operationOptions...))
}

// MustGroupGet registers related handlers and panics on error.
func MustGroupGet[I, O any](g *Group, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	lo.Must0(GroupGet(g, path, handler, operationOptions...))
}

// MustGroupPost registers related handlers and panics on error.
func MustGroupPost[I, O any](g *Group, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	lo.Must0(GroupPost(g, path, handler, operationOptions...))
}

// MustGroupPut registers related handlers and panics on error.
func MustGroupPut[I, O any](g *Group, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	lo.Must0(GroupPut(g, path, handler, operationOptions...))
}

// MustGroupPatch registers related handlers and panics on error.
func MustGroupPatch[I, O any](g *Group, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	lo.Must0(GroupPatch(g, path, handler, operationOptions...))
}

// MustGroupDelete registers related handlers and panics on error.
func MustGroupDelete[I, O any](g *Group, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	lo.Must0(GroupDelete(g, path, handler, operationOptions...))
}

func registerTypedWithPrefix[I, O any](s *Server, prefix, method, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	if s == nil {
		return fmt.Errorf("%w: server is nil", ErrRouteNotRegistered)
	}

	api := lo.Ternary(s.adapter == nil, nil, s.adapter.HumaAPI())
	if api == nil {
		return ErrAdapterNotFound
	}

	// Combine server base path with group prefix
	fullPath := joinRoutePath(s.basePath, joinRoutePath(prefix, path))

	// Always wrap handler with error handling and panic recovery
	wrappedHandler := withInputValidation(s, handler)

	opID := defaultOperationID(method, fullPath)
	handlerName := handlerName(handler)

	// Create operation and apply options
	op := huma.Operation{
		OperationID: opID,
		Method:      method,
		Path:        fullPath,
	}

	// Apply operation options
	lo.ForEach(operationOptions, func(opt OperationOption, _ int) {
		if opt != nil {
			opt(&op)
		}
	})

	huma.Register(api, op, func(ctx context.Context, input *I) (*O, error) {
		return wrappedHandler(ctx, input)
	})

	s.addRoute(RouteInfo{
		Method:      method,
		Path:        fullPath,
		HandlerName: handlerName,
	})

	return nil
}

func withInputValidation[I, O any](s *Server, handler TypedHandler[I, O]) TypedHandler[I, O] {
	if handler == nil || s == nil {
		return handler
	}

	return func(ctx context.Context, input *I) (out *O, err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				out = nil
				err = huma.Error500InternalServerError(fmt.Sprintf("panic in handler: %v", recovered))
			}
		}()

		if err = s.validateInput(input); err != nil {
			message := validationErrorMessage(err)
			return nil, huma.Error400BadRequest(message, err)
		}

		out, err = handler(ctx, input)
		if err != nil {
			// Check if it's a httpx Error and convert to huma StatusError
			var httpxErr *Error
			if errors.As(err, &httpxErr) {
				return nil, lo.Ternary(
					httpxErr.Err != nil,
					huma.NewError(httpxErr.Code, httpxErr.Message, httpxErr.Err),
					huma.NewError(httpxErr.Code, httpxErr.Message),
				)
			}

			// Check if it's already a huma StatusError
			var se huma.StatusError
			if errors.As(err, &se) {
				return nil, err
			}

			// Wrap unknown errors as 500
			return nil, huma.Error500InternalServerError(err.Error(), err)
		}
		return out, nil
	}
}

func handlerName(fn any) string {
	v := reflect.ValueOf(fn)
	if !v.IsValid() || v.Kind() != reflect.Func {
		return "unknown"
	}

	runtimeFn := runtime.FuncForPC(v.Pointer())
	return lo.Ternary(runtimeFn != nil, lo.LastOr(strings.Split(runtimeFn.Name(), "/"), "unknown"), "unknown")
}

func defaultOperationID(method, path string) string {
	cleanPath := strings.Trim(path, "/")
	cleanPath = lo.Ternary(cleanPath == "", "root", cleanPath)
	cleanPath = strings.NewReplacer("/", "-", "{", "", "}", "", ":", "").Replace(cleanPath)
	return strings.ToLower(method) + "-" + cleanPath
}
