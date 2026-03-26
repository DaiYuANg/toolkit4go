package dix

// WithModuleProviders appends provider registrations to a module.
func WithModuleProviders(providers ...ProviderFunc) ModuleOption {
	return func(spec *moduleSpec) { spec.providers.Add(providers...) }
}

// WithModuleSetups appends setup registrations to a module.
func WithModuleSetups(setups ...SetupFunc) ModuleOption {
	return func(spec *moduleSpec) { spec.setups.Add(setups...) }
}

// WithModuleInvokes appends invoke registrations to a module.
func WithModuleInvokes(invokes ...InvokeFunc) ModuleOption {
	return func(spec *moduleSpec) { spec.invokes.Add(invokes...) }
}

// WithModuleHooks appends lifecycle hook registrations to a module.
func WithModuleHooks(hooks ...HookFunc) ModuleOption {
	return func(spec *moduleSpec) { spec.hooks.Add(hooks...) }
}

// WithModuleImports appends imported modules to a module.
func WithModuleImports(modules ...Module) ModuleOption {
	return func(spec *moduleSpec) { spec.imports.Add(modules...) }
}

// WithModuleProfiles restricts a module to the listed profiles.
func WithModuleProfiles(profiles ...Profile) ModuleOption {
	return func(spec *moduleSpec) { spec.profiles.Add(profiles...) }
}

// WithModuleExcludeProfiles excludes a module from the listed profiles.
func WithModuleExcludeProfiles(profiles ...Profile) ModuleOption {
	return func(spec *moduleSpec) { spec.excludeProfiles.Add(profiles...) }
}

// WithModuleDescription sets the module description.
func WithModuleDescription(desc string) ModuleOption {
	return func(spec *moduleSpec) { spec.description = desc }
}

// WithModuleTags appends tags to a module.
func WithModuleTags(tags ...string) ModuleOption {
	return func(spec *moduleSpec) { spec.tags.Add(tags...) }
}

// WithModuleSetup appends a typed setup callback to a module.
func WithModuleSetup(fn func(*Container, Lifecycle) error) ModuleOption {
	return WithModuleSetups(Setup(fn))
}

// WithModuleDisabled sets whether the module is disabled.
func WithModuleDisabled(disabled bool) ModuleOption {
	return func(spec *moduleSpec) { spec.disabled = disabled }
}
