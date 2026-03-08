---
title: 'logx'
linkTitle: 'logx'
description: '结构化日志与 slog 互通'
weight: 6
---

## logx

`logx` 是基于 `zerolog` 构建的结构化日志包，具有基于选项的配置和可选的 `slog` 互操作性。

## 路线图

- 模块路线图：[logx roadmap](./roadmap)
- 全局路线图：[ArcGo roadmap](../roadmap)

## 功能

- 强类型级别（`TraceLevel`、`DebugLevel`、`InfoLevel`...）
- 控制台和文件输出
- 通过 `lumberjack` 进行文件轮转
- 可选调用者和全局日志记录器设置
- `slog` 桥接辅助函数
- Oops 集成辅助函数

## 快速开始

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

## 常见场景

### 1) 开发配置

```go
logger, err := logx.NewDevelopment()
if err != nil { panic(err) }
defer logger.Close()
```

### 2) 生产配置

```go
logger, err := logx.NewProduction()
if err != nil { panic(err) }
defer logger.Close()
```

### 3) 文件输出 + 轮转

```go
logger, err := logx.New(
    logx.WithConsole(false),
    logx.WithFile("./logs/app.log"),
    logx.WithFileRotation(100, 7, 20), // 100MB, 7 天，20 个备份
    logx.WithCompress(true),
)
```

### 4) 结构化字段

```go
logger.WithField("request_id", reqID).Info("request accepted")
logger.WithFields(map[string]any{
    "order_id": orderID,
    "user_id":  userID,
}).Info("order placed")
```

### 5) Context 和 slog 桥接

```go
slogLogger := logx.NewSlog(logger)
slogLogger.Info("hello", "module", "billing")
```

### 6) 从 context 附加 trace/span ID

```go
ctx := trace.ContextWithSpanContext(context.Background(), spanContext)

logx.WithFieldT(logger, "tenant", "acme").
    WithTraceContext(ctx).
    Info("request accepted")
```

## 级别辅助函数

- 从字符串解析：`ParseLevel("debug")`
- 在无效级别时 panic：`MustParseLevel("info")`
- 便利构造函数：`Trace()`、`Debug()`、`Info()`...

## 错误/Oops 辅助函数

- `logger.WithError(err).Error("operation failed")`
- `logger.LogOops(err)`
- `logger.Oops()`、`logger.Oopsf(...)`、`logger.OopsWith(ctx)`

## 操作说明

- 启用文件输出时始终调用 `Close()`。
- `Sync()` 对于 `zerolog` 是无操作的（已经是同步写入路径）。
- 仅在希望进程共享一个日志记录器实例时使用 `WithGlobalLogger()`。

## 测试技巧

- 在单元测试中使用 `WithConsole(true)`。
- 对文件轮转测试使用临时文件。
- 对配置级别测试断言 `GetLevel()` / `GetLevelString()`。

## 常见问题

### 我应该直接使用 `logx` 日志记录器还是通过 `slog`？

两者都支持：

- 使用 `logx` 方法获得直接 `zerolog` 风格。
- 如果你的应用标准化在 `slog` API 表面，使用 `NewSlog`。

### 进程退出前需要 `Sync()` 吗？

`Sync()` 对于此实现目前是无操作的。
当启用文件输出时，`Close()` 是重要的生命周期调用。

### 我可以设置全局日志记录器吗？

可以，通过 `WithGlobalLogger()` 或 `SetGlobalLogger()`。
仅在进程有意共享一个日志记录器实例时使用它。

## 故障排除

### 没有创建日志文件

检查：

- 设置了 `WithFile(path)`。
- 进程对目标目录有写入权限。
- `WithConsole(false)` 没有隐藏文件设置失败的问题。

### 预期的 debug 日志缺失

验证日志级别（`WithLevel(DebugLevel)` 或等效）。
更高级别过滤更低严重性日志。

### 日志轮转行为不符合预期

检查 `WithFileRotation(maxSizeMB, maxAgeDays, maxBackups)` 值和单位。
`maxSize` 是 MB，不是字节。

## 反模式

- 每次请求创建短生命期日志记录器实例。
- 使用文件输出时忘记 `Close()`。
- 记录高基数、无界字段而没有采样控制。
- 在可恢复业务错误路径中使用 panic/fatal 级别。
