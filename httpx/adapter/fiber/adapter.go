package fiber

import (
	"github.com/DaiYuANg/arcgo/httpx/adapter"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
)

// Adapter implements the fiber runtime bridge for httpx.
type Adapter struct {
	app  *fiber.App
	huma huma.API
}

// New constructs a fiber adapter backed by a fiber app and Huma API.
func New(app *fiber.App, opts ...adapter.HumaOptions) *Adapter {
	var a *fiber.App
	if app != nil {
		a = app
	} else {
		a = fiber.New()
	}

	humaOpts := adapter.MergeHumaOptions(opts...)
	cfg := huma.DefaultConfig(humaOpts.Title, humaOpts.Version)
	adapter.ApplyHumaConfig(&cfg, humaOpts)
	api := humafiber.New(a, cfg)

	return &Adapter{
		app:  a,
		huma: api,
	}
}

// Name returns the adapter name.
func (a *Adapter) Name() string {
	return "fiber"
}

// Router exposes the underlying fiber app.
func (a *Adapter) Router() *fiber.App {
	return a.app
}

// HumaAPI exposes the underlying Huma API.
func (a *Adapter) HumaAPI() huma.API {
	return a.huma
}
