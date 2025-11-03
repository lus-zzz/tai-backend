# Python 编译脚本使用说明

## 概述

`build.py` 是一个跨平台的 Python 编译脚本，支持 Windows 和 Linux 平台，能够自动集成 Git 版本信息，并将所有资源文件嵌入到二进制中。

## 功能特性

- ✅ 跨平台支持（Windows、Linux、macOS）
- ✅ 自动获取 Git 提交信息（commit、branch、tag）
- ✅ 版本信息嵌入到二进制文件
- ✅ **资源文件嵌入**（配置文件、静态文件全部打包到二进制中）
- ✅ **单文件部署**（只需一个二进制文件即可运行）
- ✅ 彩色终端输出
- ✅ 编译优化（去除调试符号，减小体积）
- ✅ 应用程序启动时自动记录版本信息到日志
- ✅ 提供 `--version` 命令行参数查看版本
- ✅ 提供 `/api/v1/version` API 接口获取版本信息

## 系统要求

### 必需
- Python 3.6+
- Go 1.21+

### 可选
- Git（用于获取版本信息）

## 安装

无需安装额外的 Python 包，脚本使用 Python 标准库。

## 使用方法

### 基本用法

```bash
# Windows
python build.py

# Linux/Mac
python3 build.py
```

### 命令行参数

```bash
python build.py [选项]

选项:
  -h, --help            显示帮助信息
  -v VERSION, --version VERSION
                        指定版本号（默认: 1.0.0，如果有 Git 标签则使用标签）
  -o OUTPUT, --output OUTPUT
                        指定输出目录（默认: release）
  -p {all,windows,linux,darwin,current}, --platform {all,windows,linux,darwin,current}
                        指定编译平台（默认: current）
  -c, --clean           清理之前的构建
  --skip-deps           跳过依赖下载
```

### 使用示例

#### 1. 快速编译当前平台

```bash
python build.py
```

这将编译当前平台的版本并输出到 `release` 目录。

#### 2. 编译所有平台

```bash
python build.py -p all -c
```

这将清理之前的构建，然后编译 Windows、Linux、macOS（Intel 和 ARM）所有版本。

#### 3. 指定版本号

```bash
python build.py -v 2.1.0
```

#### 4. 编译特定平台

```bash
# 仅编译 Windows 版本
python build.py -p windows

# 仅编译 Linux 版本
python build.py -p linux

# 仅编译 macOS 版本（包含 Intel 和 ARM）
python build.py -p darwin
```

#### 5. 自定义输出目录

```bash
python build.py -o dist -c
```

#### 6. 跳过依赖下载（加快编译速度）

```bash
python build.py --skip-deps
```

#### 7. 完整示例：发布版本

```bash
python build.py -v 1.0.0 -p all -o release -c
```

## 版本信息获取

### Git 版本信息

脚本会自动从 Git 获取以下信息：

- **Git Commit**: 当前提交的短 hash
- **Git Branch**: 当前分支名称
- **Git Tag**: 最近的标签（如果存在）
- **Git Status**: 是否有未提交的更改（clean/dirty）

### 版本号优先级

1. 如果指定了 `-v` 参数，使用指定的版本号
2. 如果存在 Git 标签且未指定版本号，使用 Git 标签
3. 否则使用默认版本号 `1.0.0`

### 版本标记

如果 Git 工作区有未提交的更改，版本号会自动添加 `-dirty` 后缀。

例如：`1.0.0-dirty`

## 查看应用程序版本信息

编译后的应用程序提供多种方式查看版本信息：

### 1. 命令行参数

```bash
# Windows
chat-backend-windows-amd64.exe --version

# Linux
./chat-backend-linux-amd64 --version

# 输出示例：
# Chat Backend
# Version:    1.0.0
# Build Time: 2024-01-15 10:30:00
# Git Commit: abc1234
# Git Branch: main
# Git Tag:    v1.0.0
```

### 2. 启动日志

应用程序启动时会自动将版本信息写入日志文件（`logs/` 目录）：

```
[INFO] ========================================
[INFO] Flowy 聊天后端服务启动
[INFO] 版本: 1.0.0
[INFO] 构建时间: 2024-01-15 10:30:00
[INFO] Git Commit: abc1234
[INFO] Git Branch: main
[INFO] Git Tag: v1.0.0
[INFO] 工作目录: /app
[INFO] Flowy API: http://10.18.13.10:8888/api/v1
[INFO] 服务端口: 9090
[INFO] API 文档: http://localhost:9090/swagger/index.html
[INFO] 版本信息: http://localhost:9090/api/v1/version
[INFO] 健康检查: http://localhost:9090/health
[INFO] 日志查看器: http://localhost:9090/static/log-viewer.html
[INFO] 服务器监听端口: 9090
[INFO] ========================================
```

### 3. API 接口

```bash
# 获取版本信息
curl http://localhost:9090/api/v1/version

# 响应示例：
{
  "success": true,
  "message": "获取版本信息成功",
  "data": {
    "version": "1.0.0",
    "build_time": "2024-01-15 10:30:00",
    "git_commit": "abc1234",
    "git_branch": "main",
    "git_tag": "v1.0.0"
  }
}
```

## 输出文件结构

编译完成后，输出目录结构非常简洁：

```
release/
└── chat-backend-windows-amd64.exe    # 单个二进制文件（所有资源已嵌入）
```

或多平台编译：

```
release/
├── chat-backend-windows-amd64.exe    # Windows 64位版本
├── chat-backend-linux-amd64          # Linux 64位版本
├── chat-backend-darwin-amd64         # macOS Intel版本
└── chat-backend-darwin-arm64         # macOS Apple Silicon版本
```

**注意：**
- ✅ **无需任何配置文件或静态资源目录**
- ✅ **单个二进制文件包含所有内容**
- ✅ 配置文件（`config/default_settings.json`）已嵌入二进制
- ✅ 静态文件（`static/log-viewer.html`）已嵌入二进制
- 📁 `logs` 目录会在程序首次运行时自动创建

## 编译优化

脚本使用以下编译优化选项：

- `-s`: 去除符号表，减小文件体积
- `-w`: 去除 DWARF 调试信息
- `-trimpath`: 去除文件路径信息，增强安全性
- `CGO_ENABLED=0`: 禁用 CGO，生成静态链接的二进制文件

## 资源嵌入

使用 Go 1.16+ 的 `embed` 功能，将以下资源嵌入到二进制文件中：

- **配置文件**: `config/default_settings.json`
- **静态文件**: `static/log-viewer.html`

这意味着：
- ✅ 无需随二进制分发任何额外文件
- ✅ 单个可执行文件即可完整运行
- ✅ 简化部署流程
- ✅ 避免配置文件丢失问题

如需自定义配置，可在程序运行目录创建 `config/default_settings.json`，程序会优先使用外部配置文件。

## 平台支持

| 平台 | 架构 | 输出文件名 |
|------|------|-----------|
| Windows | amd64 | chat-backend-windows-amd64.exe |
| Linux | amd64 | chat-backend-linux-amd64 |
| macOS | amd64 (Intel) | chat-backend-darwin-amd64 |
| macOS | arm64 (Apple Silicon) | chat-backend-darwin-arm64 |

## 常见问题

### 1. Python 版本错误

**问题**: `SyntaxError` 或 `f-string` 错误

**解决**: 确保使用 Python 3.6 或更高版本：

```bash
python --version
# 或
python3 --version
```

### 2. Go 环境未找到

**问题**: `未找到 Go 环境，请先安装 Go`

**解决**: 安装 Go 并确保已添加到 PATH：

```bash
go version
```

### 3. Git 命令失败

**问题**: Git 信息显示为 "unknown"

**解决**: 这不会影响编译，但如果需要 Git 信息：
- 确保已安装 Git
- 确保在 Git 仓库中运行脚本

### 4. 权限错误（Linux/Mac）

**问题**: `Permission denied`

**解决**: 给脚本添加执行权限：

```bash
chmod +x build.py
./build.py
```

### 5. Windows 颜色输出问题

**问题**: 终端显示乱码或没有颜色

**解决**: 
- 使用 Windows 10 或更高版本
- 使用 Windows Terminal
- 脚本会自动处理 ANSI 颜色支持

## 环境变量

可以通过环境变量覆盖默认配置：

```bash
# 设置输出目录
export BUILD_OUTPUT=dist

# 设置版本号
export BUILD_VERSION=2.0.0
```

## 技术实现

### 版本信息注入

脚本通过 Go 的 `-ldflags` 参数在编译时注入版本信息：

```go
var (
    Version   = "dev"
    BuildTime = "unknown"
    GitCommit = "unknown"
    GitBranch = "unknown"
    GitTag    = ""
)
```

编译时使用：

```bash
go build -ldflags "-X 'main.Version=1.0.0' -X 'main.BuildTime=2024-01-15 10:30:00' ..."
```

## 开发建议

### 开发模式

在开发过程中，使用 `current` 平台快速编译：

```bash
python build.py -p current
```

这将在 `release/` 目录生成单个二进制文件，可直接运行测试。

### 发布模式

发布新版本时，编译所有平台：

```bash
# 1. 创建 Git 标签
git tag v1.0.0
git push origin v1.0.0

# 2. 编译所有平台
python build.py -p all -c

# 3. 检查输出（只有二进制文件）
ls release/
```

### 运行和测试

```bash
# 直接运行二进制文件
cd release
./chat-backend-windows-amd64.exe

# 程序会自动：
# - 创建 logs 目录
# - 从嵌入的资源加载配置
# - 启动 HTTP 服务
# - 记录启动信息到日志
```

### CI/CD 集成

在 CI/CD 流程中使用：

```bash
# GitHub Actions 示例
python build.py -v ${{ github.ref_name }} -p all -o artifacts -c

# 构建产物：
# artifacts/chat-backend-windows-amd64.exe
# artifacts/chat-backend-linux-amd64
# artifacts/chat-backend-darwin-amd64
# artifacts/chat-backend-darwin-arm64
```

## 与 PowerShell 脚本对比

| 特性 | Python 脚本 | PowerShell 脚本 |
|------|------------|----------------|
| 跨平台 | ✅ Windows/Linux/Mac | ❌ 仅 Windows |
| Git 集成 | ✅ 完整 | ✅ 完整 |
| 彩色输出 | ✅ | ✅ |
| 资源嵌入 | ✅ 单文件部署 | ❌ 需要复制文件 |
| 输出清洁 | ✅ 只生成二进制 | ❌ 生成多个文件 |
| 依赖 | Python 3.6+ | PowerShell 5.1+ |
| 执行 | `python build.py` | `.\build-release.ps1` |

**推荐使用 Python 脚本**，因为它支持跨平台且输出更简洁。

## 许可证

与主项目相同。

## 技术支持

如有问题，请检查：
1. Python 版本 (≥ 3.6)
2. Go 版本 (≥ 1.21，建议 1.16+ 以支持 embed)
3. Git 是否可用（可选，用于版本信息）
4. 项目目录结构是否正确

编译脚本输出会以不同颜色标识：
- 🔵 蓝色 (INFO): 信息提示
- 🟢 绿色 (SUCCESS): 成功操作
- 🔴 红色 (ERROR): 错误信息

## 快速参考

```bash
# 最常用命令
python build.py                          # 编译当前平台
python build.py -p all                   # 编译所有平台
python build.py -v 1.0.0 -p all -c      # 发布版本编译

# 查看版本
./release/chat-backend.exe --version

# 运行程序
./release/chat-backend.exe

# 访问服务
# http://localhost:9090/swagger/index.html  - API 文档
# http://localhost:9090/api/v1/version     - 版本信息
# http://localhost:9090/static/log-viewer.html - 日志查看器
```

---

**最后更新**: 2025-10-30
