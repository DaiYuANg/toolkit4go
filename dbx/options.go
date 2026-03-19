package dbx

import (
	"log/slog"

	"github.com/samber/lo"
)

type Option func(*options)

type options struct {
	logger *slog.Logger
	hooks  []Hook
	debug  bool
}

func defaultOptions() options {
	return options{
		logger: slog.Default(),
		hooks:  make([]Hook, 0, 4),
		debug:  false,
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(opts *options) {
		if logger != nil {
			opts.logger = logger
		}
	}
}

func WithHooks(hooks ...Hook) Option {
	return func(opts *options) {
		opts.hooks = append(opts.hooks, lo.Filter(hooks, func(hook Hook, _ int) bool {
			return hook != nil
		})...)
	}
}

func WithDebug(enabled bool) Option {
	return func(opts *options) {
		opts.debug = enabled
	}
}

func applyOptions(opts ...Option) options {
	config := defaultOptions()
	lo.ForEach(lo.Filter(opts, func(opt Option, _ int) bool {
		return opt != nil
	}), func(opt Option, _ int) {
		opt(&config)
	})
	return config
}
