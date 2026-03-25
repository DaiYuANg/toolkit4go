---
title: 'dix Getting Started'
linkTitle: 'getting-started'
description: 'Build, start, and stop a typed module graph'
weight: 2
---

## Getting Started

This page shows a **self-contained** `dix` program:

- define a couple of typed services
- compose them into modules
- `Build()` a runtime
- `Start()` and `Stop()` cleanly

## 1) Install

```bash
go get github.com/DaiYuANg/arcgo/dix@latest
go get github.com/DaiYuANg/arcgo/logx@latest
```

## 2) Create `main.go`

```go
package main

import (
	"context"
	"log/slog"

	"github.com/DaiYuANg/arcgo/dix"
	"github.com/DaiYuANg/arcgo/logx"
)

type Config struct {
	Port int
}

type Server struct {
	Logger *slog.Logger
	Config Config
}

func main() {
	configModule := dix.NewModule("config",
		dix.WithModuleProviders(
			dix.Provider0(func() Config { return Config{Port: 8080} }),
		),
	)

	serverModule := dix.NewModule("server",
		dix.WithModuleImports(configModule),
		dix.WithModuleProviders(
			dix.Provider2(func(logger *slog.Logger, cfg Config) *Server {
				return &Server{Logger: logger, Config: cfg}
			}),
		),
		dix.WithModuleHooks(
			dix.OnStart(func(ctx context.Context, srv *Server) error {
				srv.Logger.Info("server starting", "port", srv.Config.Port)
				return nil
			}),
			dix.OnStop(func(ctx context.Context, srv *Server) error {
				srv.Logger.Info("server stopping", "port", srv.Config.Port)
				return nil
			}),
		),
	)

	logger, err := logx.NewDevelopment()
	if err != nil {
		panic(err)
	}

	app := dix.New(
		"demo",
		dix.WithProfile(dix.ProfileDev),
		dix.WithLogger(logger),
		dix.WithModules(configModule, serverModule),
	)

	if err := app.Validate(); err != nil {
		panic(err)
	}

	rt, err := app.Build()
	if err != nil {
		panic(err)
	}

	if err := rt.Start(context.Background()); err != nil {
		panic(err)
	}
	defer func() {
		_, _ = rt.StopWithReport(context.Background())
	}()
}
```

## 3) Run

```bash
go mod init example.com/dix-hello
go get github.com/DaiYuANg/arcgo/dix@latest
go get github.com/DaiYuANg/arcgo/logx@latest
go run .
```

## Next

- Health checks and `net/http` handlers: [Health and lifecycle](./health-and-lifecycle)
- Advanced features (named/alias/scope/override): see [dix examples](./examples) and `dix/advanced`

