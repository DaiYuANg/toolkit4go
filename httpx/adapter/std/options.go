package std

import (
	"github.com/DaiYuANg/arcgo/httpx/adapter"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

// New constructs a std adapter backed by a chi router and Huma API.
func New(router *chi.Mux, opts ...adapter.HumaOptions) *Adapter {
	if router == nil {
		router = chi.NewMux()
	}

	humaOpts := adapter.MergeHumaOptions(opts...)
	cfg := huma.DefaultConfig(humaOpts.Title, humaOpts.Version)
	adapter.ApplyHumaConfig(&cfg, humaOpts)
	api := humachi.New(router, cfg)

	return &Adapter{
		router:    router,
		huma:      api,
		lifecycle: &lifecycleState{},
	}
}
