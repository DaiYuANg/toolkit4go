# ArcGo 版本文档管理

本文档说明如何管理 ArcGo 的多版本文档。

## 📁 目录结构

```
docs/
├── content/
│   ├── docs/                # 当前版本文档（最新版本）
│   └── versioned/           # 历史版本文档
│       ├── v0.2.2/          # v0.2.2 版本文档
│       ├── v0.2.1/          # v0.2.1 版本文档
│       └── v0.0.1/          # v0.0.1 版本文档
├── data/
│   └── versions.yaml        # 版本配置文件
├── layouts/
│   └── partials/
│       ├── navbar/
│       │   └── end.html     # 导航栏版本切换器
│       └── version-banner.html  # 版本提示横幅
├── scripts/
│   └── sync-versions.go     # Go 版本同步脚本（跨平台）
└── VERSIONING.md            # 本文档
```

## 🚀 快速开始

### 添加新版本

当你发布新的 git tag 后，执行以下步骤：

#### 方式一：使用 Taskfile（推荐）

```bash
# 同步版本文档
go tool task docs:version:sync

# 启动开发服务器
go tool task docs:serve
```

#### 方式二：直接运行 Go 脚本

```bash
cd docs
go run scripts/sync-versions.go
```

脚本会自动：
1. 读取所有 git tags
2. 更新 `data/versions.yaml` 配置文件
3. 创建对应的版本文档目录

### 手动管理版本

#### 1. 更新 versions.yaml

编辑 `docs/data/versions.yaml`：

```yaml
versions:
  - name: "v0.3.0"          # 新版本号
    release: "v0.3.0"
    path: ""                # 当前版本路径为空
    current: true           # 标记为当前版本
    
  - name: "v0.2.2"          # 之前的版本
    release: "v0.2.2"
    path: "/versioned/v0.2.2"
    current: false
```

#### 2. 创建版本文档目录

```bash
# 复制当前文档到新版本目录
cp -r docs/content/docs docs/content/versioned/v0.2.2
```

#### 3. 更新当前文档

更新 `docs/content/docs` 下的文档内容，作为下一个版本的文档。

## 🎨 版本切换器

版本切换器位于导航栏右侧，显示为一个文档图标 📄。

- 点击图标会显示所有可用版本
- 选择版本后会跳转到对应版本的文档
- 当前版本会显示 ✓ 标记

## 📊 版本配置说明

### versions.yaml 字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `name` | string | 版本显示名称 |
| `release` | string | Git tag 名称 |
| `path` | string | 文档路径（当前版本为空） |
| `current` | boolean | 是否为当前版本 |

### 版本顺序

版本按在 `versions.yaml` 中的顺序排列，**第一个版本为当前版本**。

## 🔧 开发服务器

启动开发服务器预览版本文档：

```bash
cd docs
go tool hugo server -D
```

访问：
- 当前版本：http://localhost:1313/
- 历史版本：http://localhost:1313/versioned/v0.2.2/

## 📝 最佳实践

### 1. 发布新版本时

1. 创建 git tag：`git tag v0.3.0`
2. 运行同步脚本：`./scripts/sync-versions.sh`
3. 验证版本切换器显示正确
4. 提交更改

### 2. 文档更新

- **当前版本文档**：直接编辑 `docs/content/docs` 下的文件
- **历史版本文档**：编辑 `docs/content/versioned/<version>` 下的文件

### 3. 版本数量

建议只保留最近的 **3-5 个主要版本** 的文档，避免站点过大。

## 🐛 故障排除

### 版本切换器不显示

1. 检查 `docs/data/versions.yaml` 是否存在
2. 检查 `docs/layouts/partials/navbar/end.html` 是否正确
3. 重启 Hugo 服务器

### 版本文档内容不正确

1. 确认版本目录结构正确
2. 检查 `_index.md` 文件是否存在
3. 清除缓存：`hugo --gc`

## 📖 示例

### 查看当前版本配置

```bash
cat docs/data/versions.yaml
```

### 列出版本文档目录

```bash
ls -la docs/content/versioned/
```

## 🔗 相关资源

- [Hugo 多版本文档](https://gohugo.io/content-management/multilingual/)
- [Hextra 主题文档](https://imfing.github.io/hextra/)
