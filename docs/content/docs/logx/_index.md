---
title: 'logx'
linkTitle: 'logx'
description: 'Structured Logging with slog Interoperability'
weight: 6
---

## logx

`logx` is a structured logging package built on `zerolog` with option-based configuration and optional `slog` interoperability.

## Roadmap

- Module roadmap: [logx roadmap](./roadmap)
- Global roadmap: [ArcGo roadmap](../roadmap)

## Features

- Strongly typed levels (`TraceLevel`, `DebugLevel`, `InfoLevel`...)
- Console and file output
- File rotation via `lumberjack`
- Optional caller and global logger setting
- `slog` bridge helpers
- Oops integration helpers

## Quick Start

```go
logger, err := logx.New(
    logx.WithConsole(true),
    logx.WithLevel(logx.InfoLevel),
)
if err != nil {
    panic(err)
}
defer logger.Close()

logger.Info("service started", "service", "user-api")
```

## Common Scenarios

### 1) Development Configuration

```go
logger, err := logx.NewDevelopment()
if err != nil { panic(err) }
defer logger.Close()
```

### 2) Production Configuration

```go
logger, err := logx.NewProduction()
if err != nil { panic(err) }
defer logger.Close()
```

### 3) File Output + Rotation

```go
logger, err := logx.New(
    logx.WithConsole(false),
    logx.WithFile("./logs/app.log"),
    logx.WithFileRotation(100, 7, 20), // 100MB, 7 days, 20 backups
    logx.WithCompress(true),
)
```

### 4) Structured Fields

```go
logger.WithField("request_id", reqID).Info("request accepted")
logger.WithFields(map[string]any{
    "order_id": orderID,
    "user_id":  userID,
}).Info("order placed")
```

### 5) Context and slog Bridge

```go
slogLogger := logx.NewSlog(logger)
slogLogger.Info("hello", "module", "billing")
```

### 6) Attach trace/span ID from context

```go
ctx := trace.ContextWithSpanContext(context.Background(), spanContext)

logx.WithFieldT(logger, "tenant", "acme").
    WithTraceContext(ctx).
    Info("request accepted")
```

## Level Helpers

- Parse from string: `ParseLevel("debug")`
- Panic on invalid level: `MustParseLevel("info")`
- Convenience constructors: `Trace()`, `Debug()`, `Info()`...

## Error/Oops Helpers

- `logger.WithError(err).Error("operation failed")`
- `logger.LogOops(err)`
- `logger.Oops()`, `logger.Oopsf(...)`, `logger.OopsWith(ctx)`

## Operational Notes

- Always call `Close()` when file output is enabled.
- Only use `WithGlobalLogger()` when you intend to share a single logger instance across the process.

## Testing Tips

- Use `WithConsole(true)` in unit tests.
- Use temp files for file rotation tests.
- Assert `GetLevel()` / `GetLevelString()` for configuration level tests.

## FAQ

### Should I use the `logx` logger directly or via `slog`?

Both are supported:

- Use `logx` methods for direct `zerolog` style.
- Use `NewSlog` if your application standardizes on the `slog` API surface.

### Do I need `Close()` before process exit?

`Close()` is the critical lifecycle call when file output is enabled. If you're only logging to console, `Close()` is not required (but calling it is still safe).

### Can I set a global logger?

Yes, via `WithGlobalLogger()` or `SetGlobalLogger()`.
Only use this when the process intentionally shares a single logger instance.

## Troubleshooting

### Log file not created

Check:

- `WithFile(path)` is set.
- Process has write permission to target directory.
- `WithConsole(false)` isn't hiding file setup failures.

### Expected debug logs missing

Verify log level (`WithLevel(DebugLevel)` or equivalent).
Higher levels filter lower severity logs.

### Log rotation behavior not as expected

Check `WithFileRotation(maxSizeMB, maxAgeDays, maxBackups)` values and units.
`maxSize` is in MB, not bytes.

## Anti-Patterns

- Creating short-lived logger instances per request.
- Forgetting `Close()` when using file output.
- Logging high-cardinality, unbounded fields without sampling control.
- Using panic/fatal levels in recoverable business error paths.
