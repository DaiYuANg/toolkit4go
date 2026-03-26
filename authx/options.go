package authx

import "log/slog"

type EngineOption func(*Engine)

func WithAuthenticationManager(manager AuthenticationManager) EngineOption {
	return func(engine *Engine) {
		engine.SetAuthenticationManager(manager)
	}
}

func WithAuthorizer(authorizer Authorizer) EngineOption {
	return func(engine *Engine) {
		engine.SetAuthorizer(authorizer)
	}
}

func WithHook(hook Hook) EngineOption {
	return func(engine *Engine) {
		engine.AddHook(hook)
	}
}

func WithLogger(logger *slog.Logger) EngineOption {
	return func(engine *Engine) {
		if engine != nil && logger != nil {
			engine.logger = logger
		}
	}
}

func WithDebug(enabled bool) EngineOption {
	return func(engine *Engine) {
		if engine != nil {
			engine.debug = enabled
		}
	}
}
