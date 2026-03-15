package dix

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
	Tags            []string
}

// SetupFunc is called during container build.
type SetupFunc func(c *Container, lc Lifecycle) error

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

// WithProviders adds provider registration functions.
func WithProviders(providers ...ProviderFunc) ModuleOption {
	return func(m *Module) {
		m.Providers = append(m.Providers, providers...)
	}
}

// WithInvokes adds invoke functions.
func WithInvokes(invokes ...InvokeFunc) ModuleOption {
	return func(m *Module) {
		m.Invokes = append(m.Invokes, invokes...)
	}
}

// WithImports adds module dependencies.
func WithImports(modules ...Module) ModuleOption {
	return func(m *Module) {
		m.Imports = append(m.Imports, modules...)
	}
}

// WithProfiles sets profiles for which this module is active.
func WithProfiles(profiles ...Profile) ModuleOption {
	return func(m *Module) {
		m.Profiles = append(m.Profiles, profiles...)
	}
}

// WithExcludeProfiles sets profiles from which this module is excluded.
func WithExcludeProfiles(profiles ...Profile) ModuleOption {
	return func(m *Module) {
		m.ExcludeProfiles = append(m.ExcludeProfiles, profiles...)
	}
}

// WithDescription sets the module description.
func WithDescription(desc string) ModuleOption {
	return func(m *Module) {
		m.Description = desc
	}
}

// WithTags adds tags for module identification.
func WithTags(tags ...string) ModuleOption {
	return func(m *Module) {
		m.Tags = append(m.Tags, tags...)
	}
}

// WithSetup sets the setup function.
func WithSetup(fn SetupFunc) ModuleOption {
	return func(m *Module) {
		m.Setup = fn
	}
}

// WithDisabled marks the module as disabled.
func WithDisabled(disabled bool) ModuleOption {
	return func(m *Module) {
		m.Disabled = disabled
	}
}

// Provider helpers

func Provider0[T any](fn func() T) ProviderFunc {
	return func(c *Container) { ProvideT(c, fn) }
}

func Provider1[T, D1 any](fn func(D1) T) ProviderFunc {
	return func(c *Container) { Provide1T(c, fn) }
}

func Provider2[T, D1, D2 any](fn func(D1, D2) T) ProviderFunc {
	return func(c *Container) { Provide2T(c, fn) }
}

func Provider3[T, D1, D2, D3 any](fn func(D1, D2, D3) T) ProviderFunc {
	return func(c *Container) { Provide3T(c, fn) }
}

// Invoke helpers

func Invoke0(fn func()) InvokeFunc {
	return func(c *Container) error { dixInvoke0(c, fn); return nil }
}

func Invoke1[T any](fn func(T)) InvokeFunc {
	return func(c *Container) error { return dixInvoke1(c, fn) }
}

func Invoke2[T1, T2 any](fn func(T1, T2)) InvokeFunc {
	return func(c *Container) error { return dixInvoke2(c, fn) }
}

func Invoke3[T1, T2, T3 any](fn func(T1, T2, T3)) InvokeFunc {
	return func(c *Container) error { return dixInvoke3(c, fn) }
}

// Internal invoke helpers
func dixInvoke0(c *Container, fn func()) {
	fn()
}

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

// flattenModules recursively flattens all imported modules.
func flattenModules(modules []Module, profile Profile) []Module {
	result := make([]Module, 0)

	for _, mod := range modules {
		if mod.Disabled || !isActiveForProfile(mod, profile) {
			continue
		}

		if len(mod.Imports) > 0 {
			imported := flattenModules(mod.Imports, profile)
			result = append(result, imported...)
		}

		result = append(result, mod)
	}

	return result
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
