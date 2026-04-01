// Package adapter provides package-level APIs.
package adapter

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

// Host is the minimal runtime surface required by httpx core.
type Host interface {
	// Name returns the host name, such as `std`, `gin`, `echo`, or `fiber`.
	Name() string

	// HumaAPI exposes the underlying Huma API.
	HumaAPI() huma.API
}

// ListenableAdapter is implemented by adapters that can listen directly.
type ListenableAdapter interface {
	Listen(addr string) error
}

// ContextListenableAdapter is implemented by adapters that support context-aware shutdown.
type ContextListenableAdapter interface {
	ListenContext(ctx context.Context, addr string) error
}

// ShutdownAdapter is implemented by adapters that can stop an active listener.
type ShutdownAdapter interface {
	Shutdown() error
}

// HumaOptions configures Huma-backed metadata and docs exposure for an adapter.
type HumaOptions struct {
	// Title sets the OpenAPI info title.
	Title string
	// Version sets the OpenAPI info version.
	Version string
	// Description sets the OpenAPI info description.
	Description string
	// DocsPath sets the docs UI route.
	DocsPath string
	// OpenAPIPath sets the OpenAPI spec route prefix, without extension.
	OpenAPIPath string
	// SchemasPath sets the schema route prefix.
	SchemasPath string
	// DocsRenderer selects the built-in docs renderer.
	DocsRenderer string
	// DisableDocsRoutes disables docs, OpenAPI, and schema routes.
	DisableDocsRoutes bool
	// Transformers modifies response bodies before serialization.
	Transformers []huma.Transformer
}

// DefaultHumaOptions provides default behavior.
func DefaultHumaOptions() HumaOptions {
	return HumaOptions{
		Title:       "My API",
		Version:     "1.0.0",
		Description: "API Documentation",
		DocsPath:    "/docs",
		OpenAPIPath: "/openapi",
		SchemasPath: "/schemas",
	}
}

// MergeHumaOptions merges multiple HumaOptions into one.
// Later options override earlier options for non-empty fields.
func MergeHumaOptions(opts ...HumaOptions) HumaOptions {
	result := DefaultHumaOptions()
	for i := range opts {
		mergeHumaOption(&result, opts[i])
	}
	return result
}

func mergeHumaOption(dst *HumaOptions, opt HumaOptions) {
	if dst == nil {
		return
	}
	mergeString := func(current *string, next string) {
		if next != "" {
			*current = next
		}
	}

	mergeString(&dst.Title, opt.Title)
	mergeString(&dst.Version, opt.Version)
	mergeString(&dst.Description, opt.Description)
	mergeString(&dst.DocsPath, opt.DocsPath)
	mergeString(&dst.OpenAPIPath, opt.OpenAPIPath)
	mergeString(&dst.SchemasPath, opt.SchemasPath)
	mergeString(&dst.DocsRenderer, opt.DocsRenderer)
	if opt.DisableDocsRoutes {
		dst.DisableDocsRoutes = true
	}
	if len(opt.Transformers) > 0 {
		dst.Transformers = append(dst.Transformers, opt.Transformers...)
	}
}

// ApplyHumaConfig copies adapter Huma options into a Huma config.
func ApplyHumaConfig(cfg *huma.Config, opts HumaOptions) {
	if cfg == nil {
		return
	}

	cfg.Info.Description = opts.Description

	if opts.DisableDocsRoutes {
		cfg.DocsPath = ""
		cfg.OpenAPIPath = ""
		cfg.SchemasPath = ""
		return
	}

	cfg.DocsPath = normalizeDocsPath(opts.DocsPath)
	cfg.OpenAPIPath = normalizeOpenAPIPath(opts.OpenAPIPath)
	cfg.SchemasPath = normalizeSchemasPath(opts.SchemasPath)
	if opts.DocsRenderer != "" {
		cfg.DocsRenderer = opts.DocsRenderer
	}
	if len(opts.Transformers) > 0 {
		cfg.Transformers = append(cfg.Transformers, opts.Transformers...)
	}
}

func normalizeDocsPath(path string) string {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "/docs"
	}
	return ensureLeadingSlash(trimmed)
}

func normalizeOpenAPIPath(path string) string {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "/openapi"
	}

	trimmed = strings.TrimSuffix(trimmed, ".json")
	trimmed = strings.TrimSuffix(trimmed, ".yaml")
	trimmed = strings.TrimSuffix(trimmed, "/")
	if trimmed == "" {
		return "/openapi"
	}
	return ensureLeadingSlash(trimmed)
}

func normalizeSchemasPath(path string) string {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "/schemas"
	}

	trimmed = strings.TrimSuffix(trimmed, "/")
	if trimmed == "" {
		return "/schemas"
	}
	return ensureLeadingSlash(trimmed)
}

func ensureLeadingSlash(path string) string {
	if strings.HasPrefix(path, "/") {
		return path
	}
	return "/" + path
}
