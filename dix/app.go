package dix

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// Profile represents an application profile (environment).
type Profile string

const (
	// ProfileDefault is the default profile.
	ProfileDefault Profile = "default"
	// ProfileDev is for development environment.
	ProfileDev Profile = "dev"
	// ProfileTest is for testing environment.
	ProfileTest Profile = "test"
	// ProfileProd is for production environment.
	ProfileProd Profile = "prod"
)

// AppMeta contains application metadata.
type AppMeta struct {
	Name        string
	Version     string
	Description string
}

// AppState represents the current state of the application.
type AppState int32

const (
	// AppStateCreated indicates the app has been created but not built.
	AppStateCreated AppState = iota
	// AppStateBuilt indicates the app has been built and is ready to start.
	AppStateBuilt
	// AppStateStarted indicates the app is running.
	AppStateStarted
	// AppStateStopped indicates the app has been stopped.
	AppStateStopped
)

// App is the main application container.
// It orchestrates module loading, dependency injection, and lifecycle management.
//
// App provides a simple, typed-first approach to dependency injection
// without the complexity of traditional DI frameworks.
//
// Example:
//
//	app := NewApp("myapp",
//	    DatabaseModule,
//	    HTTPModule,
//	)
//	app.Run()
type App struct {
	meta      AppMeta
	profile   Profile
	modules   []Module
	container *Container
	lifecycle *lifecycleImpl
	state     AppState
	built     bool
}

// AppOption configures an App.
type AppOption func(*App)

// WithProfile sets the application profile.
func WithProfile(profile Profile) AppOption {
	return func(a *App) {
		a.profile = profile
	}
}

// NewApp creates a new application with the given name and modules.
//
// The app is created but not built. Call Build(), Start(), or Run()
// to initialize the application.
//
// Example:
//
//	app := NewApp("myapp",
//	    WithProviders(ProvideConfig, ProvideDatabase),
//	    WithSetup(func(c Container, lc Lifecycle) error {
//	        lc.OnStart(func(db *Database) error {
//	            return db.Connect()
//	        })
//	        return nil
//	    }),
//	)
func NewApp(name string, modules ...Module) *App {
	return &App{
		meta: AppMeta{
			Name: name,
		},
		profile:   ProfileDefault,
		modules:   modules,
		container: newContainer(),
		lifecycle: newLifecycle(),
		state:     AppStateCreated,
		built:     false,
	}
}

// NewAppWithOptions creates a new application with additional options.
func NewAppWithOptions(name string, opts []AppOption, modules ...Module) *App {
	app := NewApp(name, modules...)
	for _, opt := range opts {
		opt(app)
	}
	return app
}

// Name returns the application name.
func (a *App) Name() string {
	return a.meta.Name
}

// Profile returns the current application profile.
func (a *App) Profile() Profile {
	return a.profile
}

// Container returns the underlying container for direct access.
// This is primarily useful for testing.
func (a *App) Container() *Container {
	return a.container
}

// Run builds, starts, waits for shutdown signal, and then stops the application.
// This is the main entry point for running the application.
//
// Run performs the following steps:
// 1. Build: Loads all modules and registers providers
// 2. Start: Executes all OnStart hooks in order
// 3. Wait: Blocks until a shutdown signal is received
// 4. Stop: Executes all OnStop hooks in reverse order
//
// Example:
//
//	if err := app.Run(); err != nil {
//	    log.Fatal(err)
//	}
func (a *App) Run() error {
	// Build
	if err := a.Build(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	// Start
	ctx := context.Background()
	if err := a.Start(ctx); err != nil {
		return fmt.Errorf("start failed: %w", err)
	}

	// Wait for shutdown signal
	a.waitForShutdown()

	// Stop
	if err := a.Stop(ctx); err != nil {
		return fmt.Errorf("stop failed: %w", err)
	}

	return nil
}

// Build initializes the application by loading all modules and registering providers.
// After Build() is called, the container is ready to resolve dependencies.
//
// Build performs the following steps:
// 1. Flatten all modules (including imports)
// 2. Filter modules by profile
// 3. Register all providers in the container
// 4. Execute module Setup functions
// 5. Execute all Invokes
//
// Build is idempotent - calling it multiple times has no effect after the first call.
func (a *App) Build() error {
	if a.built {
		return nil
	}

	// Flatten modules and filter by profile
	flatModules := flattenModules(a.modules, a.profile)

	// Register all providers
	for _, mod := range flatModules {
		for _, provider := range mod.Providers {
			provider(a.container)
		}
	}

	// Execute Setup functions
	for _, mod := range flatModules {
		if mod.Setup != nil {
			if err := mod.Setup(a.container, a.lifecycle); err != nil {
				return fmt.Errorf("setup failed for module %s: %w", mod.Name, err)
			}
		}
	}

	// Execute Invokes
	for _, mod := range flatModules {
		for _, invoke := range mod.Invokes {
			if err := invoke(a.container); err != nil {
				return fmt.Errorf("invoke failed in module %s: %w", mod.Name, err)
			}
		}
	}

	a.built = true
	a.state = AppStateBuilt
	return nil
}

// Start starts the application by executing all OnStart hooks.
// Hooks are executed in the order they were registered.
//
// Start requires Build() to be called first.
func (a *App) Start(ctx context.Context) error {
	if !a.built {
		return fmt.Errorf("app must be built before starting, call Build() first")
	}

	if err := a.lifecycle.executeStartHooks(ctx, a.container); err != nil {
		return err
	}

	a.state = AppStateStarted
	return nil
}

// Stop stops the application by executing all OnStop hooks.
// Hooks are executed in reverse order of registration.
//
// Stop requires Start() to have been called.
func (a *App) Stop(ctx context.Context) error {
	if a.state != AppStateStarted {
		return fmt.Errorf("app must be started before stopping")
	}

	if err := a.lifecycle.executeStopHooks(ctx, a.container); err != nil {
		return err
	}

	// Shutdown the container
	if err := a.container.Shutdown(ctx); err != nil {
		return fmt.Errorf("container shutdown failed: %w", err)
	}

	a.state = AppStateStopped
	return nil
}

// waitForShutdown blocks until a shutdown signal is received.
func (a *App) waitForShutdown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}

// State returns the current application state.
func (a *App) State() AppState {
	return a.state
}
