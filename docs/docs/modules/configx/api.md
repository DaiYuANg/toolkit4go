---
sidebar_position: 4
---

# API 参考

本文档提供 `configx` 的完整 API 参考。

## 函数

### Load

```go
func Load(cfg any, opts ...Option) error
```

加载配置到指定的结构体中。

**参数：**
- `cfg` - 目标配置结构体指针
- `opts` - 配置选项

**示例：**

```go
var cfg Config
err := configx.Load(&cfg,
    configx.WithFiles("config.yaml"),
    configx.WithEnvPrefix("APP"),
)
```

---

### LoadConfig

```go
func LoadConfig(opts ...Option) (*Config, error)
```

加载配置并返回 Config 对象。

**返回：**
- `*Config` - 配置对象
- `error` - 错误信息

**示例：**

```go
cfg, err := configx.LoadConfig(
    configx.WithFiles("config.yaml"),
    configx.WithEnvPrefix("APP"),
)
```

---

### NewConfig

```go
func NewConfig(opts ...Option) (*Config, error)
```

创建并加载配置实例，是 `LoadConfig` 的别名。

---

## Option 选项

### WithDotenv

```go
func WithDotenv(files ...string) Option
```

启用 .env 文件加载。

**参数：**
- `files` - .env 文件路径列表，默认为 `.env`

**示例：**

```go
configx.WithDotenv()                              // 加载 .env
configx.WithDotenv(".env.local", ".env")          // 加载多个文件
```

---

### WithFiles

```go
func WithFiles(files ...string) Option
```

设置配置文件路径。

**参数：**
- `files` - 配置文件路径列表（支持 YAML、JSON、TOML）

**示例：**

```go
configx.WithFiles("config.yaml")
configx.WithFiles("config.default.yaml", "config.yaml")
```

---

### WithEnvPrefix

```go
func WithEnvPrefix(prefix string) Option
```

设置环境变量前缀。

**参数：**
- `prefix` - 环境变量前缀（如 "APP"）

**示例：**

```go
configx.WithEnvPrefix("APP")  // 加载 APP_ 开头的环境变量
```

---

### WithEnvPrefixs

```go
func WithEnvPrefixs(prefixes ...string) Option
```

设置多个环境变量前缀。

**参数：**
- `prefixes` - 环境变量前缀列表

**示例：**

```go
configx.WithEnvPrefixs("APP", "CONFIG")
```

---

### WithPriority

```go
func WithPriority(p ...Source) Option
```

设置配置源优先级。

**参数：**
- `p` - 配置源列表，后者优先级高于前者

**示例：**

```go
configx.WithPriority(
    configx.SourceDotenv,
    configx.SourceFile,
    configx.SourceEnv,
)
```

**配置源类型：**
- `SourceDotenv` - .env 文件
- `SourceFile` - 配置文件
- `SourceEnv` - 环境变量
- `SourceDefault` - 默认值

---

### WithDefaults

```go
func WithDefaults(m map[string]any) Option
```

设置默认值。

**参数：**
- `m` - 默认值 map

**示例：**

```go
configx.WithDefaults(map[string]any{
    "app.name":  "my-app",
    "app.port":  8080,
    "app.debug": false,
})
```

---

### WithValidateLevel

```go
func WithValidateLevel(level ValidateLevel) Option
```

设置验证级别。

**参数：**
- `level` - 验证级别

**验证级别：**
- `ValidateLevelNone` - 不验证（默认）
- `ValidateLevelStruct` - 验证结构体标签
- `ValidateLevelRequired` - 验证 required 标签

**示例：**

```go
configx.WithValidateLevel(configx.ValidateLevelRequired)
```

---

### WithValidator

```go
func WithValidator(v *validator.Validate) Option
```

设置自定义 validator。

**参数：**
- `v` - validator 实例

**示例：**

```go
validate := validator.New()
validate.RegisterValidation("custom", customValidator)
configx.WithValidator(validate)
```

---

## Config 方法

### GetString

```go
func (c *Config) GetString(path string) string
```

获取字符串值。

---

### GetInt

```go
func (c *Config) GetInt(path string) int
```

获取整数值。

---

### GetInt64

```go
func (c *Config) GetInt64(path string) int64
```

获取 64 位整数值。

---

### GetFloat64

```go
func (c *Config) GetFloat64(path string) float64
```

获取浮点数值。

---

### GetBool

```go
func (c *Config) GetBool(path string) bool
```

获取布尔值。

---

### GetDuration

```go
func (c *Config) GetDuration(path string) time.Duration
```

获取时长值（支持 "30s", "1m", "1h" 等格式）。

---

### GetStringSlice

```go
func (c *Config) GetStringSlice(path string) []string
```

获取字符串切片。

---

### GetIntSlice

```go
func (c *Config) GetIntSlice(path string) []int
```

获取整数切片。

---

### Get

```go
func (c *Config) Get(path string) any
```

获取任意类型的值。

---

### Exists

```go
func (c *Config) Exists(path string) bool
```

检查键是否存在。

---

### All

```go
func (c *Config) All() map[string]any
```

获取所有配置。

---

### Unmarshal

```go
func (c *Config) Unmarshal(path string, out any) error
```

解构到结构体。

**参数：**
- `path` - 配置路径，空字符串表示根路径
- `out` - 目标结构体指针

**示例：**

```go
var cfg Config
err := c.Unmarshal("", &cfg)
```

---

### Cut

```go
func (c *Config) Cut(path string) *Config
```

获取配置子树。

**参数：**
- `path` - 子树路径

**示例：**

```go
dbCfg := cfg.Cut("database")
host := dbCfg.GetString("host")
```

---

### MarshalJSON

```go
func (c *Config) MarshalJSON() ([]byte, error)
```

将配置导出为 JSON。

---

## 错误处理

### 错误类型

配置加载可能返回以下错误：

- `os.ErrNotExist` - 配置文件不存在
- `validator.ValidationErrors` - 验证失败
- 其他解析错误

### 错误示例

```go
var cfg Config
err := configx.Load(&cfg,
    configx.WithFiles("config.yaml"),
    configx.WithValidateLevel(configx.ValidateLevelRequired),
)

if err != nil {
    if errors.Is(err, os.ErrNotExist) {
        // 配置文件不存在
    }
    
    var ve validator.ValidationErrors
    if errors.As(err, &ve) {
        // 验证失败
        for _, e := range ve {
            fmt.Printf("Field: %s, Error: %s\n", e.Field(), e.Tag())
        }
    }
}
```
