package configx

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/knadh/koanf/providers/file"
)

// ChangeHandler is the signature for callbacks registered with [Watcher.OnChange].
// cfg holds the freshly loaded config; err is non-nil when the reload failed.
// When err is non-nil, cfg is nil and the previous config remains active.
type ChangeHandler func(cfg *Config, err error)

type changeHandlers []ChangeHandler

// ChangeHandlerT is the callback signature for typed hot-reload handlers.
// cfg is the newly decoded typed config value when err is nil.
type ChangeHandlerT[T any] func(cfg T, err error)

// Watcher manages a live-reloading *Config.
//
// It sets up an fsnotify watcher for every file listed in the original option
// set. Whenever any of those files is written or recreated, the Watcher
// performs a *full* reload (defaults → files → env) so that every source is
// always in sync. Multiple rapid saves are collapsed into a single reload via
// a configurable debounce window (default 100 ms).
//
// Typical usage:
//
//	w, err := configx.NewWatcher(
//	    configx.WithFiles("config.yaml"),
//	    configx.WithEnvPrefix("APP"),
//	    configx.WithWatchDebounce(200*time.Millisecond),
//	    configx.WithWatchErrHandler(func(err error) {
//	        slog.Error("config watch error", "err", err)
//	    }),
//	)
//
//	w.OnChange(func(cfg *configx.Config, err error) {
//	    if err == nil {
//	        slog.Info("config reloaded", "port", cfg.GetInt("server.port"))
//	    }
//	})
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	go w.Start(ctx)
//
//	// Always use w.Config() to get the latest snapshot.
//	port := w.Config().GetInt("server.port")
type Watcher struct {
	// cfg is replaced atomically after each successful reload.
	cfg atomic.Pointer[Config]

	opts *Options

	// subsMu serializes subscriber registration; notify reads an immutable
	// snapshot through subs without taking a lock.
	subsMu sync.Mutex
	subs   atomic.Pointer[changeHandlers]

	// providers are used *only* for change detection – actual loading is
	// always done by a fresh call to loadConfigFromOptions. They are immutable
	// after construction, so a plain slice remains the cheapest representation.
	providers []*file.File

	// stopCh is closed by Close to signal the Start loop to exit.
	stopCh   chan struct{}
	stopOnce sync.Once
}

// NewWatcher builds a Watcher from the supplied options, performs the initial
// config load, and prepares fsnotify watchers for every supported config file.
//
// Call [Watcher.Start] (typically in a goroutine) to begin watching.
func NewWatcher(opts ...Option) (*Watcher, error) {
	options := NewOptions()
	for _, opt := range opts {
		if opt != nil {
			opt(options)
		}
	}
	return newWatcherFromOptions(context.Background(), options)
}

// newWatcherFromOptions is the internal constructor shared by NewWatcher and
// Loader.NewWatcher so that the options pointer is reused without re-applying
// functional options a second time.
func newWatcherFromOptions(ctx context.Context, opts *Options) (*Watcher, error) {
	cfg, err := loadConfigFromOptions(ctx, opts)
	if err != nil {
		logError(opts, "configx watcher initial load failed", "error", err)
		return nil, fmt.Errorf("configx: watcher initial load: %w", err)
	}

	w := &Watcher{
		opts:      opts,
		providers: buildWatchProviders(opts.files),
		stopCh:    make(chan struct{}),
	}
	w.cfg.Store(cfg)
	logDebug(opts, "configx watcher created", "providers", len(w.providers))
	return w, nil
}

// Config returns the most recently successfully loaded config snapshot.
// It is safe to call from multiple goroutines.
func (w *Watcher) Config() *Config {
	return w.cfg.Load()
}

// OnChange registers fn to be called after every reload attempt.
//
//   - On success: cfg is the new config, err is nil.
//   - On failure: cfg is nil, err describes what went wrong; the previous
//     config remains active (w.Config() is unchanged).
//
// Handlers are invoked in registration order from a single goroutine, so they
// do not need to be goroutine-safe relative to each other.  Heavy work should
// be dispatched to a separate goroutine to avoid blocking the reload loop.
func (w *Watcher) OnChange(fn ChangeHandler) {
	if fn == nil {
		return
	}
	w.subsMu.Lock()
	defer w.subsMu.Unlock()

	current := w.loadSubscribers()
	next := make(changeHandlers, len(current), len(current)+1)
	copy(next, current)
	next = append(next, fn)
	w.subs.Store(&next)
}

// Start begins watching config files for changes and blocks until ctx is
// canceled or [Watcher.Close] is called.
//
// If no files are configured Start simply waits for the context to be done, so
// it is always safe to run in a goroutine regardless of the option set.
//
// Errors from individual file watchers are forwarded to the handler registered
// with [WithWatchErrHandler]; Start itself only returns a non-nil error when
// it cannot set up an fsnotify watcher for a file.
func (w *Watcher) Start(ctx context.Context) error {
	ctx = normalizeWatcherContext(ctx)

	// Nothing to watch - block until signaled.
	if len(w.providers) == 0 {
		logDebug(w.opts, "configx watcher started without providers")
		return w.waitForStop(ctx)
	}

	debounce := normalizeWatchDebounce(w.opts.watchDebounce)
	reloadCh := make(chan struct{}, 1)
	if err := w.startProviders(func() {
		queueWatcherReload(reloadCh)
	}); err != nil {
		return err
	}

	logDebug(w.opts, "configx watcher started", "providers", len(w.providers), "debounce_ms", debounce.Milliseconds())
	return w.run(ctx, debounce, reloadCh)
}

func normalizeWatcherContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}

	return ctx
}

func normalizeWatchDebounce(debounce time.Duration) time.Duration {
	if debounce <= 0 {
		return 100 * time.Millisecond
	}

	return debounce
}

func queueWatcherReload(reloadCh chan<- struct{}) {
	select {
	case reloadCh <- struct{}{}:
	default:
	}
}

func (w *Watcher) startProviders(trigger func()) error {
	for i, fp := range w.providers {
		if err := fp.Watch(w.watchProvider(i, trigger)); err != nil {
			w.cleanupStartedProviders(i)
			logError(w.opts, "configx watcher start failed", "index", i, "error", err)
			return fmt.Errorf("configx: start file watcher: %w", err)
		}
	}

	return nil
}

func (w *Watcher) watchProvider(index int, trigger func()) func(_ any, err error) {
	return func(_ any, err error) {
		if err != nil {
			logError(w.opts, "configx watcher provider error", "index", index, "error", err)
			w.handleErr(fmt.Errorf("configx: fsnotify error on file %d: %w", index, err))
			return
		}

		logDebug(w.opts, "configx watcher change detected", "index", index)
		trigger()
	}
}

func (w *Watcher) cleanupStartedProviders(count int) {
	for i := range count {
		if err := w.providers[i].Unwatch(); err != nil {
			w.handleErr(fmt.Errorf("configx: cleanup file watcher %d: %w", i, err))
		}
	}
}

func (w *Watcher) run(ctx context.Context, debounce time.Duration, reloadCh <-chan struct{}) error {
	resetTimer, stopTimer := newDebounceTimer(debounce, func() {
		w.reload(ctx)
	})
	defer stopTimer()
	for {
		select {
		case <-ctx.Done():
			if err := w.Close(); err != nil {
				w.handleErr(fmt.Errorf("configx: close watcher: %w", err))
			}
			return nil

		case <-w.stopCh:
			return nil

		case <-reloadCh:
			logDebug(w.opts, "configx watcher reload queued")
			resetTimer()
		}
	}
}

func newDebounceTimer(debounce time.Duration, fn func()) (reset, stop func()) {
	var (
		timer   *time.Timer
		timerMu sync.Mutex
	)

	reset = func() {
		timerMu.Lock()
		defer timerMu.Unlock()
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(debounce, fn)
	}

	stop = func() {
		timerMu.Lock()
		defer timerMu.Unlock()
		if timer != nil {
			timer.Stop()
		}
	}

	return reset, stop
}

func (w *Watcher) waitForStop(ctx context.Context) error {
	select {
	case <-ctx.Done():
	case <-w.stopCh:
	}

	return nil
}

// Close stops all file watchers and unblocks [Watcher.Start].
// It is idempotent and safe to call from multiple goroutines.
func (w *Watcher) Close() error {
	w.stopOnce.Do(func() { close(w.stopCh) })
	logDebug(w.opts, "configx watcher closing", "providers", len(w.providers))

	var errs []error
	for _, fp := range w.providers {
		if err := fp.Unwatch(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		logError(w.opts, "configx watcher close completed with errors", "errors", len(errs))
	} else {
		logDebug(w.opts, "configx watcher closed")
	}
	return errors.Join(errs...)
}

// ─── internal helpers ─────────────────────────────────────────────────────────

// reload performs a full config reload and notifies subscribers.
func (w *Watcher) reload(ctx context.Context) {
	logDebug(w.opts, "configx watcher reload started")
	newCfg, err := loadConfigFromOptions(ctx, w.opts)
	if err != nil {
		wrapped := fmt.Errorf("configx: reload failed: %w", err)
		logError(w.opts, "configx watcher reload failed", "error", wrapped)
		w.handleErr(wrapped)
		w.notify(nil, wrapped)
		return
	}

	w.cfg.Store(newCfg)
	logDebug(w.opts, "configx watcher reload completed")

	w.notify(newCfg, nil)
}

// notify calls every registered ChangeHandler in order.
func (w *Watcher) notify(cfg *Config, err error) {
	logDebug(w.opts, "configx watcher notifying subscribers", "subscribers", len(w.loadSubscribers()), "has_error", err != nil)
	for _, fn := range w.loadSubscribers() {
		fn(cfg, err)
	}
}

// handleErr forwards err to the watchErrHandler when one is configured.
func (w *Watcher) handleErr(err error) {
	if err == nil || w.opts.watchErrHandler == nil {
		return
	}
	w.opts.watchErrHandler(err)
}

func (w *Watcher) loadSubscribers() []ChangeHandler {
	subs := w.subs.Load()
	if subs == nil {
		return nil
	}
	return *subs
}

// buildWatchProviders creates one *file.File provider per supported config
// file path. These providers are used exclusively for change detection;
// loadConfigFromOptions handles the actual reading and parsing.
func buildWatchProviders(paths []string) []*file.File {
	out := make([]*file.File, 0, len(paths))
	for _, p := range paths {
		switch filepath.Ext(p) {
		case ".yaml", ".yml", ".json", ".toml":
			out = append(out, file.Provider(p))
		}
	}
	return out
}

// WatcherT provides typed hot-reload snapshots on top of Watcher.
type WatcherT[T any] struct {
	base    *Watcher
	current atomic.Pointer[T]
}

func newWatcherTFromOptions[T any](ctx context.Context, opts *Options) (*WatcherT[T], error) {
	base, err := newWatcherFromOptions(ctx, opts)
	if err != nil {
		return nil, err
	}

	var initial T
	if err := base.Config().Unmarshal("", &initial); err != nil {
		return nil, fmt.Errorf("initial typed watcher unmarshal: %w", errors.Join(ErrUnmarshal, err))
	}
	if err := base.Config().validateStruct(initial); err != nil {
		return nil, fmt.Errorf("initial typed watcher value: %w", errors.Join(ErrValidate, err))
	}

	w := &WatcherT[T]{base: base}
	w.current.Store(&initial)
	return w, nil
}

// Config returns the latest successfully decoded typed snapshot.
func (w *WatcherT[T]) Config() T {
	ptr := w.current.Load()
	if ptr == nil {
		var zero T
		return zero
	}
	return *ptr
}

// RawConfig returns the underlying dynamic config snapshot.
func (w *WatcherT[T]) RawConfig() *Config {
	return w.base.Config()
}

// OnChange registers a typed callback. Decode/validate failures are surfaced
// via err and do not replace the current typed snapshot.
func (w *WatcherT[T]) OnChange(fn ChangeHandlerT[T]) {
	if fn == nil {
		return
	}
	w.base.OnChange(func(cfg *Config, err error) {
		var zero T
		if err != nil {
			fn(zero, err)
			return
		}
		var out T
		if err := cfg.Unmarshal("", &out); err != nil {
			wrapped := fmt.Errorf("watcher typed callback decode: %w", errors.Join(ErrUnmarshal, err))
			w.base.handleErr(wrapped)
			fn(zero, wrapped)
			return
		}
		if err := cfg.validateStruct(out); err != nil {
			wrapped := fmt.Errorf("watcher typed callback value: %w", errors.Join(ErrValidate, err))
			w.base.handleErr(wrapped)
			fn(zero, wrapped)
			return
		}
		w.current.Store(&out)
		fn(out, nil)
	})
}

// Start starts the underlying watcher loop.
func (w *WatcherT[T]) Start(ctx context.Context) error {
	return w.base.Start(ctx)
}

// Close stops the underlying watcher.
func (w *WatcherT[T]) Close() error {
	return w.base.Close()
}
