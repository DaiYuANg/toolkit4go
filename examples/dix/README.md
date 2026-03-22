# examples/dix

Canonical example documentation now lives in the Hugo docs:

- [dix examples](../../docs/content/docs/dix/examples.md)

Run examples from this module:

```bash
go run ./basic
go run ./runtime_scope
go run ./inspect
go run ./backend   # configx + logx + eventx + httpx + dix + dbx(SQLite) 真实后端示例
```

## backend

模拟真实后端进程，集成：

- **configx**：配置加载（默认 dotenv → file → env）
- **logx**：日志
- **eventx**：事件总线（UserCreatedEvent）
- **httpx**：HTTP API（chi + Huma）
- **dix**：依赖注入
- **dbx**：SQLite CRUD

包结构：`config`、`domain`、`repo`（dbx）、`service`、`api`、`event`、`db`。

```bash
go run ./backend
# 可选: APP_SERVER_PORT=3000 APP_DB_DSN=file:app.db
# 访问 http://localhost:8080/docs
```
