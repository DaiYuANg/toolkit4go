package service

import (
	"log/slog"

	"github.com/DaiYuANg/arcgo/dix"
	"github.com/DaiYuANg/arcgo/eventx"
	"github.com/DaiYuANg/arcgo/examples/dix/backend/event"
	"github.com/DaiYuANg/arcgo/examples/dix/backend/repo"
)

var Module = dix.NewModule("service",
	dix.WithModuleImports(repo.Module, event.Module),
	dix.WithModuleProviders(
		dix.Provider3(func(r repo.UserRepository, bus eventx.BusRuntime, log *slog.Logger) UserService {
			return NewUserService(r, bus, log)
		}),
	),
)
