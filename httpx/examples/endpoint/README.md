# Endpoint Pattern Example

This example demonstrates how to organize HTTP handlers using the **Endpoint pattern** (similar to Controller-based design).

## Run the Example

```bash
go run ./httpx/examples/endpoint
```

## Key Concepts

### 1. Define an Endpoint

Embed `httpx.BaseEndpoint` and implement `RegisterRoutes`:

```go
type UserEndpoint struct {
    httpx.BaseEndpoint
}

func (e *UserEndpoint) RegisterRoutes(server *httpx.Server) {
    api := server.Group("/api/v1/users")
    
    _ = httpx.GroupGet(api, "/{id}", getUserHandler)
    _ = httpx.GroupPost(api, "", createUserHandler)
}
```

### 2. Register Endpoints

**Simple registration (no hooks):**

```go
server.RegisterOnly(
    &HealthEndpoint{},
    &UserEndpoint{},
    &OrderEndpoint{},
)
```

**Registration with hooks:**

```go
server.Register(&OrderEndpoint{},
    httpx.EndpointHooks{
        Before: func(s *httpx.Server, e httpx.Endpoint) {
            // Register middleware, setup, etc.
        },
        After: func(s *httpx.Server, e httpx.Endpoint) {
            // Logging, cleanup, etc.
        },
    },
)
```

### 3. Benefits

- **Modular code**: Each endpoint encapsulates related routes
- **Testable**: Test each endpoint independently
- **Organized**: Group routes by business domain (User, Order, Health, etc.)
- **Flexible**: Use hooks for cross-cutting concerns

## API Endpoints

After running the example:

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/api/v1/users` | List users |
| GET | `/api/v1/users/{id}` | Get user by ID |
| POST | `/api/v1/users` | Create user |
| POST | `/api/v1/orders` | Create order |

## Documentation

- OpenAPI JSON: http://localhost:8080/openapi.json
- Swagger UI: http://localhost:8080/docs
