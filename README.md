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

### API 调用

```bash
elements api /api/users/this                                 # GET 请求
elements api /api/devices -q page=0 -q limit=10              # 带 query params
elements api /api/devices -X POST -f name=test               # POST body fields
echo '{}' | elements api /api/devices -X POST --input -      # 从 stdin 读取 JSON body
elements api /api/users/this -H "Sudo: user@example.com"     # 自定义 header
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

# 设备流量
elements device traffic monthly 202604 --device <device-id>   # 查询月流量
elements device traffic daily 202604 <device-id>              # 查询日流量

# 设备客户端
elements device clients list <device-id>                      # 列设备接入的客户端
elements device clients batch <device-id>...                  # 批量查询客户端

# 设备告警
elements device alert                                         # 列告警
elements device alert --device-name router --state unconfirmed # 按条件过滤

# 设备配置
elements device config get <device-id>                        # 获取运行配置
elements device config set <device-id> --content "..."        # 下发配置
```

### 设备分组 (`devicegroup`, `dg`)

```bash
elements devicegroup list                                     # 列分组
elements devicegroup list --parent <parent-id>                # 按父分组过滤
elements devicegroup get <group-id>                           # 分组详情
elements devicegroup create --name "Factory A"                # 创建分组
elements devicegroup create --name "Line 1" --parent <id>     # 创建子分组
elements devicegroup update <group-id> --name "New Name"      # 更新分组名
elements devicegroup delete <group-id>                        # 删除分组

# 分组内设备管理
elements devicegroup devices <group-id> list                  # 列分组内设备
elements devicegroup devices <group-id> list --recursive      # 包含子分组设备
elements devicegroup devices <group-id> add <device-id>...    # 添加设备到分组
elements devicegroup devices <group-id> remove <device-id>... # 从分组移除设备
elements devicegroup devices <group-id> available             # 可添加到分组的设备
```

### 远程隧道 (`tunnel`)

```bash
elements tunnel list                                          # 列隧道
elements tunnel list --device-id <id>                         # 按设备过滤
elements tunnel create \
  --name ssh-tunnel \
  --device-id <id> \
  --proto tcp \
  --local-address 127.0.0.1 \
  --local-port 22               # 创建隧道
elements tunnel update <tunnel-id> --name "new-name"          # 更新隧道
elements tunnel delete <tunnel-id>                            # 删除隧道
elements tunnel connect <tunnel-id>                           # 连接隧道
elements tunnel disconnect <tunnel-id>                        # 断开隧道
```

### DRC 配置模板 (`drc`)

```bash
elements drc list                                             # 列配置模板
elements drc list --model IR615                               # 按设备型号过滤
elements drc get <template-id>                                # 模板详情
elements drc create \
  --name "IR615-default" \
  --model IR615 \
  --content "..."               # 创建模板
elements drc delete <template-id>                             # 删除模板

# 模板设备管理
elements drc devices <template-id> list                       # 列已分配设备
elements drc devices <template-id> list --status running      # 按状态过滤
elements drc devices <template-id> add <device-id>...         # 分配设备
elements drc devices <template-id> add <device-id> --group <group-id>  # 分配设备组
elements drc devices <template-id> remove <device-id>         # 移除设备
elements drc devices <template-id> restart <device-id>        # 重启设备任务
```

### 边缘计算 (`edge`)

#### 边缘引擎 (`edge agent`)

```bash
elements edge agent list                                # 列引擎
elements edge agent list --version v1.0                 # 按版本过滤
elements edge agent get <agent-id>                       # 引擎详情
elements edge agent upload <file-path> --description "IR615 engine"  # 上传引擎
elements edge agent update <agent-id> --description "new desc"       # 更新引擎
elements edge agent delete <agent-id>                    # 删除引擎
elements edge agent devices <agent-id>                   # 已部署设备列表
elements edge agent devices <agent-id> --status READY    # 按状态过滤
```

#### 边缘应用 (`edge app`)

```bash
elements edge app list                                   # 列应用
elements edge app get <app-id>                           # 应用详情
elements edge app create --name "my-app" --description "..."          # 创建应用
elements edge app update <app-id> --description "new desc"           # 更新应用
elements edge app delete <app-id>                        # 删除应用
```

#### 应用版本 (`edge version`)

```bash
elements edge version list <app-id>                      # 列版本
elements edge version upload <file-path> --app <app-id>  # 上传版本
elements edge version update <app-id> <version> --notes "Release notes"  # 更新日志
elements edge version delete <app-id> <version>          # 删除版本
elements edge version deploy <app-id> <version> --device <id> --group <id>  # 部署版本
```

#### 应用配置 (`edge config`)

```bash
elements edge config list <app-id>                       # 列配置
elements edge config list <app-id> --version v1.0        # 按版本过滤
elements edge config get <app-id> <config-id>            # 配置详情
elements edge config create <app-id> --version v1.0 --content "..."    # 创建配置
elements edge config update <app-id> <config-id> --description "..."   # 更新配置
elements edge config delete <app-id> <config-id>         # 删除配置
elements edge config deploy <app-id> <version> --device <id> --group <id>  # 部署配置
```

#### 远程控制 (`edge control`)

```bash
elements edge control start <device-id> <app-id>         # 启动应用
elements edge control stop <device-id> <app-id>          # 停止应用
elements edge control restart <device-id> <app-id>       # 重启应用
```

### 固件管理 (`firmware`)

```bash
elements firmware list                                      # 列固件
elements firmware list --model IR615                        # 按型号过滤
elements firmware upload <file-path>                        # 上传固件文件
elements firmware create \
  --fid <file-id> \
  --name "IR615-v2.0" \
  --version 2.0.0 \
  --model IR615             # 创建固件记录
elements firmware upgrade <device-id> --firmware-id <id>    # 单台设备升级

# 批量升级管理
elements firmware devices <firmware-id> list                # 列升级任务中的设备
elements firmware devices <firmware-id> add <device-id>...  # 批量添加设备升级
elements firmware devices <firmware-id> add --group <group-id>...  # 按组升级
elements firmware devices <firmware-id> remove <device-id>  # 取消设备升级
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
| `json`（默认） | 彩色 pretty JSON | 紧凑 JSON |
| `table` | 对齐表格 | TSV |
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

路径：`~/.config/elements/config.yaml`（权限 `0600`）

配置文件存储所有 context 信息（host、token 等），通过 `elements config` 子命令管理。

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
  cmd/              # 各子命令实现
    auth/           # 登录、登出、认证状态
    config/         # Context 管理
    device/         # 设备管理
    devicegroup/    # 设备分组管理
    tunnel/         # 远程隧道管理
    drc/            # DRC 配置模板管理
    edge/           # 边缘计算（引擎/应用/版本/配置/控制）
    firmware/       # 固件管理与升级
    version/        # 版本信息
  cmdutil/          # 通用 list flag（cursor/limit/verbose）、query 构建
  config/           # 配置文件读写、Context 模型
  debug/            # 调试输出（--debug / ELEMENTS_DEBUG）
  factory/          # 依赖注入工厂
  iostreams/        # 终端输出、格式化（JSON/Table/YAML/jq）
```
