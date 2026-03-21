# Hugo 文档部署到 GitHub Pages

## 前置条件

- 已安装 Go（含 `go tool hugo`，由 `docs/go.mod` 的 tool 提供）
- 本仓库已配置好 `origin` 远程

## 一、修改 baseURL（部署前必做）

当前 `hugo.yaml` 中的 `baseURL: https://arcgo.example.org/` 需改为 GitHub Pages 地址。

**项目页（Project Site）** 格式：`https://<username>.github.io/<repo>/`

编辑 `docs/hugo.yaml`：

```yaml
baseURL: https://DaiYuANg.github.io/arcgo/
```

> 若使用自定义域名，将 `baseURL` 改为你的域名即可。

## 二、部署方式

### 方式 A：本地手动部署

```bash
# 1. 同步版本信息（从 git tags）
go run ./scripts/deploy-docs sync

# 2. 构建并部署到 gh-pages 分支
go tool task docs:deploy
```

或分步执行：

```bash
cd docs
go tool hugo --gc --minify    # 构建到 docs/public
# 然后由 deploy 脚本将 public/ 推送到 gh-pages
go run ./scripts/deploy-docs deploy
```

环境变量（可选）：

- `DOCS_REMOTE`：远程名称，默认 `origin`
- `DOCS_BRANCH`：部署分支，默认 `gh-pages`

### 方式 B：GitHub Actions 自动部署

在仓库根目录添加 `.github/workflows/deploy-docs.yml`（见下），每次推送到 `main` 时自动构建并部署到 `gh-pages`。

## 三、GitHub 仓库设置

1. 打开仓库：https://github.com/DaiYuANg/arcgo
2. **Settings** → **Pages**
3. **Source** 选择 **Deploy from a branch**
4. **Branch** 选择 `gh-pages`，路径选择 `/ (root)`
5. 保存后等待数分钟

部署完成后访问：**https://DaiYuANg.github.io/arcgo/**

## 四、常用命令速查

| 命令 | 说明 |
|------|------|
| `go tool task docs:serve` | 本地预览（http://127.0.0.1:1313） |
| `go tool task docs:build` | 仅构建，不部署 |
| `go tool task docs:version:sync` | 同步版本元数据 |
| `go tool task docs:deploy` | 构建并部署到 gh-pages |
