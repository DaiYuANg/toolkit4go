package dix

import (
	"fmt"
	"strings"

	do "github.com/samber/do/v2"
)

// ProviderFunc is a function that registers providers in a container.
type ProviderFunc func(c *Container)

// InvokeFunc is a function that executes an invoke with a container.
type InvokeFunc func(c *Container) error

// Module is the core building block for composing applications.
type Module struct {
	Name            string
	Description     string
	Providers       []ProviderFunc
	Invokes         []InvokeFunc
	Imports         []Module
	Profiles        []Profile
	ExcludeProfiles []Profile
	Disabled        bool
	Setup           SetupFunc
	DoSetup         DoSetupFunc
	Tags            []string
}

// SetupFunc is called during container build.
type SetupFunc func(c *Container, lc Lifecycle) error

// DoSetupFunc is a narrow escape hatch for do-specific integration work.
// Keep this at the module/framework boundary, not in normal business code.
type DoSetupFunc func(raw do.Injector) error

// NewModule creates a new Module with the given name and options.
func NewModule(name string, opts ...ModuleOption) Module {
	m := Module{
		Name:      name,
		Providers: make([]ProviderFunc, 0),
		Invokes:   make([]InvokeFunc, 0),
		Imports:   make([]Module, 0),
		Profiles:  make([]Profile, 0),
		Tags:      make([]string, 0),
		Disabled:  false,
	}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}

// ModuleOption configures a Module
type ModuleOption func(*Module)

func WithModuleProviders(providers ...ProviderFunc) ModuleOption {
	return func(m *Module) { m.Providers = append(m.Providers, providers...) }
}
func WithModuleInvokes(invokes ...InvokeFunc) ModuleOption {
	return func(m *Module) { m.Invokes = append(m.Invokes, invokes...) }
}
func WithModuleImports(modules ...Module) ModuleOption {
	return func(m *Module) { m.Imports = append(m.Imports, modules...) }
}
func WithModuleProfiles(profiles ...Profile) ModuleOption {
	return func(m *Module) { m.Profiles = append(m.Profiles, profiles...) }
}
func WithModuleExcludeProfiles(profiles ...Profile) ModuleOption {
	return func(m *Module) { m.ExcludeProfiles = append(m.ExcludeProfiles, profiles...) }
}
func WithModuleDescription(desc string) ModuleOption {
	return func(m *Module) { m.Description = desc }
}
func WithModuleTags(tags ...string) ModuleOption {
	return func(m *Module) { m.Tags = append(m.Tags, tags...) }
}
func WithModuleSetup(fn SetupFunc) ModuleOption {
	return func(m *Module) { m.Setup = fn }
}
func WithModuleDoSetup(fn DoSetupFunc) ModuleOption {
	return func(m *Module) { m.DoSetup = fn }
}
func WithModuleDisabled(disabled bool) ModuleOption {
	return func(m *Module) { m.Disabled = disabled }
}

func Provider0[T any](fn func() T) ProviderFunc       { return func(c *Container) { ProvideT(c, fn) } }
func Provider1[T, D1 any](fn func(D1) T) ProviderFunc { return func(c *Container) { Provide1T(c, fn) } }
func Provider2[T, D1, D2 any](fn func(D1, D2) T) ProviderFunc {
	return func(c *Container) { Provide2T(c, fn) }
}
func Provider3[T, D1, D2, D3 any](fn func(D1, D2, D3) T) ProviderFunc {
	return func(c *Container) { Provide3T(c, fn) }
}
func Provider4[T, D1, D2, D3, D4 any](fn func(D1, D2, D3, D4) T) ProviderFunc {
	return func(c *Container) { Provide4T(c, fn) }
}
func Provider5[T, D1, D2, D3, D4, D5 any](fn func(D1, D2, D3, D4, D5) T) ProviderFunc {
	return func(c *Container) { Provide5T(c, fn) }
}
func Provider6[T, D1, D2, D3, D4, D5, D6 any](fn func(D1, D2, D3, D4, D5, D6) T) ProviderFunc {
	return func(c *Container) { Provide6T(c, fn) }
}

func Invoke0(fn func()) InvokeFunc { return func(c *Container) error { dixInvoke0(c, fn); return nil } }
func Invoke1[T any](fn func(T)) InvokeFunc {
	return func(c *Container) error { return dixInvoke1(c, fn) }
}
func Invoke2[T1, T2 any](fn func(T1, T2)) InvokeFunc {
	return func(c *Container) error { return dixInvoke2(c, fn) }
}
func Invoke3[T1, T2, T3 any](fn func(T1, T2, T3)) InvokeFunc {
	return func(c *Container) error { return dixInvoke3(c, fn) }
}
func Invoke4[T1, T2, T3, T4 any](fn func(T1, T2, T3, T4)) InvokeFunc {
	return func(c *Container) error { return dixInvoke4(c, fn) }
}
func Invoke5[T1, T2, T3, T4, T5 any](fn func(T1, T2, T3, T4, T5)) InvokeFunc {
	return func(c *Container) error { return dixInvoke5(c, fn) }
}
func Invoke6[T1, T2, T3, T4, T5, T6 any](fn func(T1, T2, T3, T4, T5, T6)) InvokeFunc {
	return func(c *Container) error { return dixInvoke6(c, fn) }
}

func dixInvoke0(c *Container, fn func()) { fn() }
func dixInvoke1[T any](c *Container, fn func(T)) error {
	t, err := ResolveAs[T](c)
	if err != nil {
		return err
	}
	fn(t)
	return nil
}
func dixInvoke2[T1, T2 any](c *Container, fn func(T1, T2)) error {
	t1, err := ResolveAs[T1](c)
	if err != nil {
		return err
	}
	t2, err := ResolveAs[T2](c)
	if err != nil {
		return err
	}
	fn(t1, t2)
	return nil
}
func dixInvoke3[T1, T2, T3 any](c *Container, fn func(T1, T2, T3)) error {
	t1, err := ResolveAs[T1](c)
	if err != nil {
		return err
	}
	t2, err := ResolveAs[T2](c)
	if err != nil {
		return err
	}
	t3, err := ResolveAs[T3](c)
	if err != nil {
		return err
	}
	fn(t1, t2, t3)
	return nil
}
func dixInvoke4[T1, T2, T3, T4 any](c *Container, fn func(T1, T2, T3, T4)) error {
	t1, err := ResolveAs[T1](c)
	if err != nil {
		return err
	}
	t2, err := ResolveAs[T2](c)
	if err != nil {
		return err
	}
	t3, err := ResolveAs[T3](c)
	if err != nil {
		return err
	}
	t4, err := ResolveAs[T4](c)
	if err != nil {
		return err
	}
	fn(t1, t2, t3, t4)
	return nil
}
func dixInvoke5[T1, T2, T3, T4, T5 any](c *Container, fn func(T1, T2, T3, T4, T5)) error {
	t1, err := ResolveAs[T1](c)
	if err != nil {
		return err
	}
	t2, err := ResolveAs[T2](c)
	if err != nil {
		return err
	}
	t3, err := ResolveAs[T3](c)
	if err != nil {
		return err
	}
	t4, err := ResolveAs[T4](c)
	if err != nil {
		return err
	}
	t5, err := ResolveAs[T5](c)
	if err != nil {
		return err
	}
	fn(t1, t2, t3, t4, t5)
	return nil
}
func dixInvoke6[T1, T2, T3, T4, T5, T6 any](c *Container, fn func(T1, T2, T3, T4, T5, T6)) error {
	t1, err := ResolveAs[T1](c)
	if err != nil {
		return err
	}
	t2, err := ResolveAs[T2](c)
	if err != nil {
		return err
	}
	t3, err := ResolveAs[T3](c)
	if err != nil {
		return err
	}
	t4, err := ResolveAs[T4](c)
	if err != nil {
		return err
	}
	t5, err := ResolveAs[T5](c)
	if err != nil {
		return err
	}
	t6, err := ResolveAs[T6](c)
	if err != nil {
		return err
	}
	fn(t1, t2, t3, t4, t5, t6)
	return nil
}

// flattenModules recursively flattens all imported modules.
func flattenModules(modules []Module, profile Profile) ([]Module, error) {
	result := make([]Module, 0)
	visited := make(map[string]struct{})
	visiting := make(map[string]struct{})

	var walk func(mod Module, path []string) error
	walk = func(mod Module, path []string) error {
		if mod.Disabled || !isActiveForProfile(mod, profile) {
			return nil
		}
		key := moduleKey(mod)
		if _, ok := visited[key]; ok {
			return nil
		}
		if _, ok := visiting[key]; ok {
			return fmt.Errorf("module import cycle detected: %s -> %s", formatModulePath(path), key)
		}
		visiting[key] = struct{}{}
		path = append(path, key)
		for _, imported := range mod.Imports {
			if err := walk(imported, path); err != nil {
				return err
			}
		}
		delete(visiting, key)
		visited[key] = struct{}{}
		result = append(result, mod)
		return nil
	}

	for _, mod := range modules {
		if err := walk(mod, nil); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func moduleKey(mod Module) string {
	if mod.Name != "" {
		return mod.Name
	}
	return fmt.Sprintf("<anonymous:%p>", &mod)
}
func formatModulePath(path []string) string {
	if len(path) == 0 {
		return "<root>"
	}
	return strings.Join(path, " -> ")
}
func isActiveForProfile(mod Module, profile Profile) bool {
	if mod.Disabled {
		return false
	}
	for _, p := range mod.ExcludeProfiles {
		if p == profile {
			return false
		}
	}
	if len(mod.Profiles) > 0 {
		for _, p := range mod.Profiles {
			if p == profile {
				return true
			}
		}
		return false
	}
	return true
}
