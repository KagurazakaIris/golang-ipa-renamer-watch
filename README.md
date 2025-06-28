# golang-ipa-renamer-watch

一个跨平台 Go 工具，自动监听目录下新增或变更的 IPA 文件并按自定义规则重命名。支持 Docker、GitHub Actions，可作为库或独立 watcher 使用。

## 功能特性
- 目录监听，自动检测新增/变更 IPA 文件（基于 [fsnotify](https://github.com/fsnotify/fsnotify)）
- 支持自定义重命名模板（如 `$raw@$CFBundleIdentifier`）
- 自动解包 Info.plist 并用 [howett.net/plist](https://github.com/DHowett/go-plist) 解析
- 跨平台：macOS、Linux、Windows、ARM/AMD64
- 支持 Docker、Compose 一键部署
- GitHub Actions 多平台镜像自动构建推送 GHCR
- 附带 Rust 原型（`rs/`，已弃用），原理来源 [xhofe/ipa-renamer](https://github.com/xhofe/ipa-renamer/blob/main/src/lib.rs)

## 快速开始

### 1. 本地构建运行
```sh
git clone https://github.com/KagurazakaIris/golang-ipa-renamer-watch.git
cd golang-ipa-renamer-watch
go build -o ipa-renamer main.go
./ipa-renamer
# 默认环境变量如下
# watchDir = ./watched
# outputDir = ./output
# 以下无特殊需要不要修改
# template = "$raw@$CFBundleIdentifier" 
# tempDir = "./temp"
```

### 2. Docker 运行
```sh
docker build -t ipa-renamer .
docker run -v $(pwd)/watched:/app/watched -v $(pwd)/output:/app/output ipa-renamer
```

### 3. Docker Compose 一键部署
```sh
docker-compose up -d
```

### 4. GHCR 拉取镜像
```sh
docker pull ghcr.io/KagurazakaIris/golang-ipa-renamer-watch:latest
docker run -v $(pwd)/watched:/app/watched -v $(pwd)/output:/app/output ghcr.io/KagurazakaIris/golang-ipa-renamer-watch:latest
```

## 配置参数
- `WATCH_DIR`：监听目录（默认 `/app/watched`）
- `OUTPUT_DIR`：重命名输出目录（默认 `/app/output`）
- `TEMPLATE`：重命名模板（默认 `$raw@$CFBundleIdentifier`）
- `TEMP_DIR`：解包 Info.plist 临时目录（默认 `/app/temp`）

所有参数均可通过环境变量或 Docker Compose `environment` 字段配置。

## 目录结构
```
├── main.go              # 目录监听与自动重命名主程序
├── ipa_renamer.go       # IPA重命名核心逻辑，可复用为库
├── Dockerfile           # 多平台构建支持的Dockerfile
├── docker-compose.yml   # 一键部署示例
├── .github/workflows/ghcr-docker.yml # GHCR自动构建流水线
├── rs/                  # Rust原型实现（已弃用）
...
```

## CI/CD
- `.github/workflows/ghcr-docker.yml`：推送到 `rewrite`/`master` 分支自动多平台构建并推送到 GHCR
- 支持 `linux/amd64` 和 `linux/arm64` 并行构建

## License

MIT License © 2025 神楽坂アイリス

