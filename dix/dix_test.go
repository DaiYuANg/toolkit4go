package dix_test

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"testing"

	"github.com/DaiYuANg/arcgo/dix"
	do "github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Example providers

type Config struct {
	DSN  string
	Port int
}

func ProvideConfig() Config {
	return Config{
		DSN:  "sqlite://test.db",
		Port: 8080,
	}
}

type Database struct {
	dsn string
}

func NewDatabase(dsn string) *Database {
	return &Database{dsn: dsn}
}

func (d *Database) Connect() error {
	fmt.Println("Database connected:", d.dsn)
	return nil
}

func (d *Database) Close() error {
	fmt.Println("Database closed")
	return nil
}

func ProvideDatabase(cfg Config) *Database {
	return NewDatabase(cfg.DSN)
}

type Server struct {
	addr string
}

func ProvideServer(cfg Config) *Server {
	return &Server{
		addr: fmt.Sprintf(":%d", cfg.Port),
	}
}

func (s *Server) Start() error {
	fmt.Println("Server starting on", s.addr)
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	fmt.Println("Server stopped")
	return nil
}

// Example modules

var DatabaseModule = dix.NewModule("database",
	dix.WithModuleProviders(
		dix.Provider0(ProvideConfig),
		dix.Provider1(ProvideDatabase),
	),
	dix.WithModuleSetup(func(c *dix.Container, lc dix.Lifecycle) error {
		dix.OnStartHook[*Database](c, func(ctx context.Context, db *Database) error {
			return db.Connect()
		})(lc)

		dix.OnStopHook[*Database](c, func(ctx context.Context, db *Database) error {
			return db.Close()
		})(lc)

		return nil
	}),
)

var ServerModule = dix.NewModule("server",
	dix.WithModuleProviders(
		dix.Provider1(ProvideServer),
	),
	dix.WithModuleImports(DatabaseModule),
	dix.WithModuleSetup(func(c *dix.Container, lc dix.Lifecycle) error {
		dix.OnStartHook[*Server](c, func(ctx context.Context, s *Server) error {
			return s.Start()
		})(lc)

		dix.OnStopHook[*Server](c, func(ctx context.Context, s *Server) error {
			return s.Stop(ctx)
		})(lc)

		return nil
	}),
)

// Tests

func TestApp_Build(t *testing.T) {
	app := dix.NewApp("testapp", ServerModule)

	err := app.Build()
	require.NoError(t, err)

	db, err := dix.ResolveAs[*Database](app.Container())
	require.NoError(t, err)
	assert.NotNil(t, db)
	assert.Equal(t, "sqlite://test.db", db.dsn)

	cfg, err := dix.ResolveAs[Config](app.Container())
	require.NoError(t, err)
	assert.Equal(t, 8080, cfg.Port)
}

func TestApp_StartStop(t *testing.T) {
	app := dix.NewApp("testapp", ServerModule)
	ctx := context.Background()

	err := app.Build()
	require.NoError(t, err)

	err = app.Start(ctx)
	require.NoError(t, err)

	err = app.Stop(ctx)
	require.NoError(t, err)
}

func TestApp_Run(t *testing.T) {
	app := dix.NewApp("testapp", ServerModule)

	err := app.Build()
	require.NoError(t, err)
	assert.Equal(t, dix.AppStateBuilt, app.State())

	ctx := context.Background()
	err = app.Start(ctx)
	require.NoError(t, err)
	assert.Equal(t, dix.AppStateStarted, app.State())

	err = app.Stop(ctx)
	require.NoError(t, err)
	assert.Equal(t, dix.AppStateStopped, app.State())
}

func TestModule_WithProfiles(t *testing.T) {
	devOnlyCalled := false
	prodOnlyCalled := false

	devModule := dix.NewModule("dev-only",
		dix.WithModuleProfiles(dix.ProfileDev),
		dix.WithModuleProviders(
			dix.Provider0(func() string {
				devOnlyCalled = true
				return "dev"
			}),
		),
	)

	prodModule := dix.NewModule("prod-only",
		dix.WithModuleProfiles(dix.ProfileProd),
		dix.WithModuleProviders(
			dix.Provider0(func() string {
				prodOnlyCalled = true
				return "prod"
			}),
		),
	)

	appDev := dix.New("test",
		dix.WithProfile(dix.ProfileDev),
		dix.WithModules(devModule, prodModule),
	)

	err := appDev.Build()
	require.NoError(t, err)

	devStr, err := dix.ResolveAs[string](appDev.Container())
	require.NoError(t, err)
	assert.Equal(t, "dev", devStr)
	assert.True(t, devOnlyCalled)
	assert.False(t, prodOnlyCalled)

	appProd := dix.New("test",
		dix.WithProfile(dix.ProfileProd),
		dix.WithModules(devModule, prodModule),
	)

	err = appProd.Build()
	require.NoError(t, err)

	prodStr, err := dix.ResolveAs[string](appProd.Container())
	require.NoError(t, err)
	assert.Equal(t, "prod", prodStr)
}

func TestModule_WithExcludeProfiles(t *testing.T) {
	called := false

	module := dix.NewModule("not-test",
		dix.WithModuleExcludeProfiles(dix.ProfileTest),
		dix.WithModuleProviders(
			dix.Provider0(func() string {
				called = true
				return "value"
			}),
		),
	)

	appTest := dix.New("test",
		dix.WithProfile(dix.ProfileTest),
		dix.WithModule(module),
	)

	err := appTest.Build()
	require.NoError(t, err)

	_, err = dix.ResolveAs[string](appTest.Container())
	assert.Error(t, err)
	assert.False(t, called)
}

func TestResolveAs(t *testing.T) {
	app := dix.NewApp("test",
		dix.NewModule("test",
			dix.WithModuleProviders(
				dix.Provider0(ProvideConfig),
			),
		),
	)

	err := app.Build()
	require.NoError(t, err)

	cfg, err := dix.ResolveAs[Config](app.Container())
	require.NoError(t, err)
	assert.Equal(t, 8080, cfg.Port)
}

func TestMustResolveAs(t *testing.T) {
	app := dix.NewApp("test",
		dix.NewModule("test",
			dix.WithModuleProviders(
				dix.Provider0(ProvideConfig),
			),
		),
	)

	err := app.Build()
	require.NoError(t, err)

	cfg := dix.MustResolveAs[Config](app.Container())
	assert.Equal(t, 8080, cfg.Port)
}

func TestInvoke(t *testing.T) {
	invoked := false

	app := dix.NewApp("test",
		dix.NewModule("test",
			dix.WithModuleProviders(
				dix.Provider0(ProvideConfig),
			),
			dix.WithModuleInvokes(
				dix.Invoke1(func(cfg Config) {
					invoked = true
					assert.Equal(t, 8080, cfg.Port)
				}),
			),
		),
	)

	err := app.Build()
	require.NoError(t, err)
	assert.True(t, invoked)
}

func TestModule_Imports(t *testing.T) {
	baseModule := dix.NewModule("base",
		dix.WithModuleProviders(
			dix.Provider0(ProvideConfig),
		),
	)

	dependentModule := dix.NewModule("dependent",
		dix.WithModuleImports(baseModule),
		dix.WithModuleProviders(
			dix.Provider1(ProvideDatabase),
		),
	)

	app := dix.NewApp("test", dependentModule)

	err := app.Build()
	require.NoError(t, err)

	cfg, err := dix.ResolveAs[Config](app.Container())
	require.NoError(t, err)

	db, err := dix.ResolveAs[*Database](app.Container())
	require.NoError(t, err)

	assert.NotNil(t, cfg)
	assert.NotNil(t, db)
}

func TestModule_ImportDeduplicatesSharedDependency(t *testing.T) {
	called := 0

	shared := dix.NewModule("shared",
		dix.WithModuleProviders(
			dix.Provider0(func() string {
				called++
				return "shared"
			}),
		),
	)

	left := dix.NewModule("left", dix.WithModuleImports(shared))
	right := dix.NewModule("right", dix.WithModuleImports(shared))

	app := dix.NewApp("test", left, right)
	require.NoError(t, app.Build())

	value, err := dix.ResolveAs[string](app.Container())
	require.NoError(t, err)
	assert.Equal(t, "shared", value)
	assert.Equal(t, 1, called)
}

func TestApp_ValidateGraph(t *testing.T) {
	app := dix.NewApp("test", DatabaseModule, ServerModule)
	require.NoError(t, app.Validate())
}

func TestProvider4AndAggregateDependencyStruct(t *testing.T) {
	type Params struct {
		Config Config
		DB     *Database
		Server *Server
		Label  string
	}

	module := dix.NewModule("test",
		dix.WithModuleProviders(
			dix.Provider0(ProvideConfig),
			dix.Provider1(ProvideDatabase),
			dix.Provider1(ProvideServer),
			dix.Provider0(func() string { return "ok" }),
			dix.Provider4(func(cfg Config, db *Database, srv *Server, label string) Params {
				return Params{Config: cfg, DB: db, Server: srv, Label: label}
			}),
		),
	)

	app := dix.NewApp("test", module)
	require.NoError(t, app.Build())

	params, err := dix.ResolveAs[Params](app.Container())
	require.NoError(t, err)
	assert.Equal(t, 8080, params.Config.Port)
	assert.Equal(t, "sqlite://test.db", params.DB.dsn)
	assert.Equal(t, ":8080", params.Server.addr)
	assert.Equal(t, "ok", params.Label)
}

func TestLifecycleHookReceivesContext(t *testing.T) {
	type ctxKey string
	const key ctxKey = "trace"

	received := ""
	module := dix.NewModule("ctx",
		dix.WithModuleProviders(dix.Provider0(func() string { return "value" })),
		dix.WithModuleSetup(func(c *dix.Container, lc dix.Lifecycle) error {
			dix.OnStartHook[string](c, func(ctx context.Context, value string) error {
				received = ctx.Value(key).(string) + ":" + value
				return nil
			})(lc)
			return nil
		}),
	)

	app := dix.NewApp("test", module)
	require.NoError(t, app.Build())
	require.NoError(t, app.Start(context.WithValue(context.Background(), key, "abc")))
	assert.Equal(t, "abc:value", received)
}

func TestContainerRegisterProviderDefinitionReturnsError(t *testing.T) {
	app := dix.NewApp("test")
	err := app.Container().Register(dix.Definition{Kind: dix.DefinitionProvider})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}

// Examples

func ExampleApp() {
	app := dix.NewApp("myapp",
		DatabaseModule,
		ServerModule,
	)
	_ = app
}

func ExampleNewModule() {
	module := dix.NewModule("example",
		dix.WithModuleProviders(
			dix.Provider0(func() string { return "hello" }),
			dix.Provider1(func(s string) int { return len(s) }),
		),
		dix.WithModuleSetup(func(c *dix.Container, lc dix.Lifecycle) error {
			lc.OnStart(func(ctx context.Context) error {
				s, _ := dix.ResolveAs[string](c)
				fmt.Println("Starting with:", s)
				return nil
			})
			return nil
		}),
	)

	app := dix.NewApp("test", module)
	_ = app.Build()
}

func ExampleHook() {
	module := dix.NewModule("example",
		dix.WithModuleProviders(
			dix.Provider0(func() *http.Server {
				return &http.Server{Addr: ":8080"}
			}),
		),
		dix.WithModuleSetup(func(c *dix.Container, lc dix.Lifecycle) error {
			dix.OnStartHook[*http.Server](c, func(ctx context.Context, s *http.Server) error {
				go s.ListenAndServe()
				return nil
			})(lc)

			dix.OnStopHook[*http.Server](c, func(ctx context.Context, s *http.Server) error {
				return s.Shutdown(ctx)
			})(lc)

			return nil
		}),
	)

	app := dix.NewApp("test", module)
	_ = app
}

func TestResolveOptionalAs(t *testing.T) {
	app := dix.NewApp("test")
	require.NoError(t, app.Build())
	_, ok := dix.ResolveOptionalAs[string](app.Container())
	assert.False(t, ok)
}

func TestWithDoSetup(t *testing.T) {
	called := false
	module := dix.NewModule("advanced",
		dix.WithModuleDoSetup(func(raw do.Injector) error {
			called = raw != nil
			return nil
		}),
	)
	app := dix.NewApp("test", module)
	require.NoError(t, app.Build())
	assert.True(t, called)
}

func TestHealthCheckReport(t *testing.T) {
	module := dix.NewModule("health",
		dix.WithModuleSetup(func(c *dix.Container, lc dix.Lifecycle) error {
			c.RegisterHealthCheck("db", func(ctx context.Context) error { return nil })
			c.RegisterHealthCheck("cache", func(ctx context.Context) error { return fmt.Errorf("down") })
			return nil
		}),
	)
	app := dix.NewApp("test", module)
	require.NoError(t, app.Build())
	report := app.CheckHealth(context.Background())
	assert.False(t, report.Healthy())
	require.Error(t, report.Error())
	assert.Contains(t, report.Error().Error(), "cache")
}

func TestNew_WithModulesOption(t *testing.T) {
	app := dix.New("test",
		dix.WithProfile(dix.ProfileDev),
		dix.WithModule(DatabaseModule),
	)

	err := app.Build()
	require.NoError(t, err)

	logger, err := dix.ResolveAs[*slog.Logger](app.Container())
	require.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestHealthKinds(t *testing.T) {
	mod := dix.NewModule("health",
		dix.WithModuleSetup(func(c *dix.Container, lc dix.Lifecycle) error {
			c.RegisterLivenessCheck("live", func(ctx context.Context) error { return nil })
			c.RegisterReadinessCheck("ready", func(ctx context.Context) error { return nil })
			return nil
		}),
	)
	app := dix.New("health-app", dix.WithModule(mod))
	require.NoError(t, app.Build())

	live := app.CheckLiveness(context.Background())
	ready := app.CheckReadiness(context.Background())

	assert.True(t, live.Healthy())
	assert.True(t, ready.Healthy())
	assert.Len(t, live.Checks, 1)
	assert.Len(t, ready.Checks, 1)
}

func TestNewDefault(t *testing.T) {
	app := dix.NewDefault()
	assert.Equal(t, dix.DefaultAppName, app.Name())
}
