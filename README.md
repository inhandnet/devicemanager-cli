# elements CLI

InHand Device Manager (DM) 平台的命令行工具，支持认证、多环境 context 切换、设备管理及多种输出格式。

## 安装

### 从源码构建

```bash
# 需要 Go 1.24+
make build    # 输出到 bin/elements
make install  # 安装到 $GOPATH/bin
```

> macOS 下必须 `CGO_ENABLED=0` 构建（Makefile 已默认设置），否则可能遇到 dyld LC_UUID 错误。

### 跨平台构建

CI 会自动构建以下平台的二进制文件：

- `linux/amd64`、`linux/arm64`
- `darwin/amd64`、`darwin/arm64`
- `windows/amd64`

## 快速开始

### 1. 登录

```bash
elements auth login                          # 默认登录中国区 (iot.inhand.com.cn)
elements auth login --host global            # 登录全球区 (iot.inhandnetworks.com)
elements auth login --host iot.example.com   # 自定义域名
elements auth login --context prod           # 创建/更新指定 context
```

登录使用 OAuth 2.0 Authorization Code 流程，会自动打开浏览器完成授权。

CLI 复用平台前端的 SPA OAuth client，登录回调由本地启动的回调服务（默认 `http://localhost:18920/callback`）接收授权码，并自动换取 Token。

### 2. 验证

```bash
elements auth status
elements device list
```

## 命令速查

### 认证

```bash
elements auth login                    # Browser-based OAuth login
elements auth status                   # View current auth status
elements auth logout                   # Log out
```

### Context 管理

Context 在登录时通过 `--context` 创建/更新。其他子命令用于切换、查看、删除：

```bash
elements config use-context <name>
elements config current-context
elements config list-contexts
elements config delete-context <name>
```

### 设备管理

```bash
elements device list                                          # 列设备（默认 limit 20）
elements device list --online 1 --model IR615                 # 按状态/型号过滤
elements device list --name router-01 --serial-number GL5022  # 按名称/SN 过滤
elements device list --cursor 20 --limit 50                   # 分页：跳过 20 条，取 50 条
elements device list --verbose 100 -o json                    # 完整字段 + JSON 输出

elements device get <device-id> --verbose 100                 # 设备详情
elements device create --name <name> --serial-number <sn>     # 添加设备
elements device signal <device-id> --after <ISO> --before <ISO>  # 历史信号质量
elements device kick <device-id>                              # 强制断开
elements device reboot <device-id> --timeout 15000            # 重启（毫秒）
```

### 调试

```bash
elements device list --debug                              # 输出 config/auth/HTTP 调试信息到 stderr
ELEMENTS_DEBUG=1 elements device list                     # 通过环境变量开启
elements device list --debug -o json 2>/tmp/debug.log     # 调试信息写文件，不影响 stdout
```

### 全局 Flag

```bash
elements --context prod auth status            # 临时切换 context
elements --debug device list                   # 开启调试输出
elements --jq '.[].name' device list           # 用 jq 表达式过滤 JSON
elements version                                # 查看版本
```

## 输出格式

通过 `-o` 指定输出格式：

| 格式 | TTY 行为 | 管道行为 |
|------|---------|---------|
| `table`（TTY 默认） | 对齐表格 | TSV |
| `json`（非 TTY 默认） | 彩色 pretty JSON | 紧凑 JSON |
| `yaml` | YAML | YAML |

```bash
elements device list -o table                   # 表格输出
elements device list -o yaml                    # YAML 输出
elements device list --jq '.[] | .name'         # 通过 jq 表达式过滤
```

服务端返回的 `{"result": ...}` 信封会在 yaml/json/jq 模式下自动剥掉，只保留 `result` 内的内容。

### 字段筛选（`--verbose`）

DM API 用 `verbose` 参数控制返回字段的详尽程度（1-100，越高越详细）：

```bash
elements device list --verbose 10                # 默认列表字段
elements device list --verbose 100 -o json       # 完整字段
elements device get <id> --verbose 100           # 获取设备完整详情
```

### 分页

```bash
elements device list --cursor 0 --limit 20       # 第一页（默认）
elements device list --cursor 20 --limit 50      # 跳过 20 条，取 50 条
```

DM 平台使用 `cursor`（skip 偏移）+ `limit` 分页。`list` 命令也接受 `--page-size`/`--per-page` 作为 `--limit` 的隐藏别名。

## 环境变量

| 变量 | 作用 |
|------|------|
| `ELEMENTS_CONTEXT` | 覆盖当前 context |
| `ELEMENTS_HOST` | 覆盖 context 中的 host |
| `ELEMENTS_TOKEN` | 覆盖 context 中的 token |
| `ELEMENTS_DEBUG` | 设为任意非空值开启调试输出 |

## 配置文件

路径：`<UserConfigDir>/elements/config.yaml`（权限 `0600`）

- Linux/macOS：`~/.config/elements/config.yaml`
- Windows：`%AppData%\elements\config.yaml`

配置文件存储所有 context 信息（host、token、refresh_token、user、过期时间），通过 `elements auth login --context <name>` 创建/更新，通过 `elements config` 子命令切换和管理。Token 401 时会自动用 `refresh_token` 刷新并写回。

## 开发指南

### 前置依赖

- Go 1.24+
- [golangci-lint](https://golangci-lint.run/)
- [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports)

### 构建 & 测试

```bash
make build       # 构建到 bin/elements
make build-all   # 跨平台构建
make install     # 安装到 GOPATH
make test        # 运行测试
make fmt         # gofmt -w .
make lint        # 运行 golangci-lint
make clean       # 清理构建产物
```

### 项目结构

```
cmd/elements/       # CLI 入口
internal/
  api/              # OAuth 认证、Token 传输与自动刷新、REST 客户端、回调服务
  build/            # 注入 Version/Commit/Date
  cmd/              # 各子命令实现（auth、config、device、version）
  cmdutil/          # 通用 list flag（cursor/limit/verbose）、query 构建
  config/           # 配置文件读写、Context 模型
  debug/            # 调试输出（--debug / ELEMENTS_DEBUG）
  factory/          # 依赖注入工厂
  iostreams/        # 终端输出、格式化（JSON/Table/YAML/jq）
```
