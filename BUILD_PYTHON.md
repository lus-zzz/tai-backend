# Python 编译脚本使用说明

## 概述

`build.py` 是一个跨平台的 Python 编译脚本，支持 Windows 和 Linux 平台，能够自动集成 Git 版本信息，并将所有资源文件嵌入到二进制中。

## 功能特性

- ✅ 跨平台支持（Windows、Linux、macOS）
- ✅ 自动获取 Git 提交信息（commit、branch、tag）
- ✅ **智能版本管理**（自动版本号生成 + Git 标签管理）
- ✅ 版本信息嵌入到二进制文件
- ✅ **资源文件嵌入**（配置文件、静态文件全部打包到二进制中）
- ✅ **单文件部署**（只需一个二进制文件即可运行）
- ✅ **Swagger 文档自动生成**
- ✅ 彩色终端输出
- ✅ 编译优化（去除调试符号，减小体积）
- ✅ **Git 标签自动创建和推送**
- ✅ 应用程序启动时自动记录版本信息到日志
- ✅ 提供 `--version` 命令行参数查看版本
- ✅ 提供 `/api/v1/version` API 接口获取版本信息

## 系统要求

### 必需
- Python 3.6+
- Go 1.21+

### 可选
- Git（用于获取版本信息和标签管理）
- Swagger 工具（用于生成 API 文档，脚本会自动安装）

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
                        指定版本号（默认: auto - 自动生成版本号）
  -o OUTPUT, --output OUTPUT
                        指定输出目录（默认: release）
  -p {all,windows,linux,darwin,current}, --platform {all,windows,linux,darwin,current}
                        指定编译平台（默认: current）
  -c, --clean           清理之前的构建
  --skip-deps           跳过依赖下载
  --skip-swagger        跳过 Swagger 文档生成
  --no-tag              跳过自动创建和推送 Git 标签
  --no-push             创建标签但不推送到远程仓库
  --force-tag           强制在当前分支创建标签，忽略分支检查
```

### 使用示例

#### 1. 快速编译当前平台

```bash
python build.py
```

这将使用自动版本号编译当前平台的版本并输出到 `release` 目录，同时会：
- 自动生成版本号（基于提交计数）
- 生成 Swagger API 文档
- 创建并推送 Git 标签（如果在正确分支且无未提交更改）

#### 2. 编译所有平台

```bash
python build.py -p all -c
```

这将清理之前的构建，然后编译 Windows、Linux、macOS（Intel 和 ARM）所有版本。

#### 3. 指定版本号

```bash
python build.py -v 2.1.0
```

#### 4. 自动版本号模式

```bash
python build.py -v auto
```

脚本会根据 Git 提交计数自动生成版本号，格式为 `major.minor.patch`。

#### 5. 编译特定平台

```bash
# 仅编译 Windows 版本
python build.py -p windows

# 仅编译 Linux 版本
python build.py -p linux

# 仅编译 macOS 版本（包含 Intel 和 ARM）
python build.py -p darwin
```

#### 6. 自定义输出目录

```bash
python build.py -o dist -c
```

#### 7. 跳过依赖下载和 Swagger 生成（加快编译速度）

```bash
python build.py --skip-deps --skip-swagger
```

#### 8. 禁用 Git 标签管理

```bash
# 编译但不创建和推送标签
python build.py --no-tag

# 创建标签但不推送到远程
python build.py --no-push

# 强制在非主分支创建标签
python build.py --force-tag
```

#### 9. 完整示例：发布版本

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

### 版本号生成策略

脚本支持多种版本号生成策略：

#### 1. 自动版本号（默认）

```bash
python build.py -v auto
# 或
python build.py  # 默认使用 auto
```

基于 Git 提交计数自动生成版本号：
- 格式：`major.minor.patch`
- 计算方式：
  - major = 提交总数 ÷ 1000
  - minor = (提交总数 % 1000) ÷ 100
  - patch = 提交总数 % 100

示例：
- 提交计数 1: `0.0.1`
- 提交计数 50: `0.0.50`
- 提交计数 150: `0.1.50`
- 提交计数 1250: `1.2.50`

#### 2. 手动指定版本号

```bash
python build.py -v 1.2.3
```

使用指定的版本号。

### Git 标签自动管理

脚本会自动处理 Git 标签的创建和推送：

#### 自动标签策略

1. **分支检查**: 仅在推荐分支（main, master, release, develop）创建标签
2. **状态检查**: 仅在工作区干净（无未提交更改）时创建标签
3. **重复检查**: 如果标签已存在，跳过创建
4. **自动推送**: 默认推送标签到远程仓库

#### 标签格式

- 格式：`v{version}`
- 示例：`v1.0.0`, `v0.1.50`

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
- **Swagger 文档**: `docs/swagger.json`（自动生成）

这意味着：
- ✅ 无需随二进制分发任何额外文件
- ✅ 单个可执行文件即可完整运行
- ✅ 包含完整的 API 文档
- ✅ 简化部署流程
- ✅ 避免配置文件丢失问题

如需自定义配置，可在程序运行目录创建 `config/default_settings.json`，程序会优先使用外部配置文件。

## Swagger 文档生成

脚本会自动生成 Swagger API 文档：

### 自动安装 Swagger 工具

如果系统中没有 `swagger` 工具，脚本会自动安装：

```bash
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
```

### 文档生成

- 生成位置：`chat-backend/docs/swagger.json`
- 嵌入到二进制：无需单独分发文档文件
- 访问地址：`http://localhost:9090/swagger/index.html`

### 跳过文档生成

如果不需要 Swagger 文档或遇到问题，可以跳过：

```bash
python build.py --skip-swagger
```

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

### 6. Swagger 工具安装失败

**问题**: `swagger 安装失败` 或 `swagger 工具不可用`

**解决**:
- 确保 Go 环境正确配置
- 确保 `$GOPATH/bin` 或 `$GOBIN` 在 PATH 中
- 手动安装：`go install github.com/go-swagger/go-swagger/cmd/swagger@latest`
- 或使用 `--skip-swagger` 跳过文档生成

### 7. Git 标签相关问题

**问题**: `标签创建失败` 或 `推送标签失败`

**解决**:
- 确保有 Git 推送权限
- 检查远程仓库配置：`git remote -v`
- 使用 `--no-tag` 跳过标签管理
- 使用 `--no-push` 仅创建本地标签
- 使用 `--force-tag` 强制在当前分支创建标签

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
# 1. 确保代码已提交（避免 -dirty 版本）
git add .
git commit -m "Release v1.0.0"

# 2. 编译所有平台（会自动创建和推送标签）
python build.py -v 1.0.0 -p all -c

# 3. 检查输出（只有二进制文件）
ls release/

# 标签会自动创建并推送到远程仓库
```

### 自动版本发布模式

使用自动版本号进行发布：

```bash
# 1. 提交代码
git add .
git commit -m "Add new features"

# 2. 自动编译和发布
python build.py -p all -c

# 版本号会根据提交计数自动生成
# 标签会自动创建和推送
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
# GitHub Actions 示例 - 手动版本号
python build.py -v ${{ github.ref_name }} -p all -o artifacts -c --no-tag

# GitHub Actions 示例 - 自动版本号
python build.py -v auto -p all -o artifacts -c --no-tag

# 构建产物：
# artifacts/chat-backend-windows-amd64.exe
# artifacts/chat-backend-linux-amd64
# artifacts/chat-backend-darwin-amd64
# artifacts/chat-backend-darwin-arm64
```

### CI/CD 中的标签管理策略

**为什么在 CI/CD 中使用 `--no-tag`？**

CI/CD 环境中通常使用 `--no-tag` 参数，因为标签管理由 CI/CD 系统统一处理，避免冲突和重复操作：

#### 1. **GitHub Actions 标签管理**

```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - 'v*'  # 当推送 v* 标签时触发

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Python
        uses: actions/setup-python@v4
        with:
          python-version: '3.9'
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Build binaries
        run: |
          # 从标签名提取版本号 (v1.0.0 -> 1.0.0)
          VERSION=${GITHUB_REF#refs/tags/v}
          python build.py -v $VERSION -p all -o artifacts -c --no-tag
      
      - name: Create Release
        uses: actions/create-release@v1
        with:
          tag_name: ${{ github.ref_name }}
          release_name: Release ${{ github.ref_name }}
          files: artifacts/*
```

#### 2. **GitLab CI 标签管理**

```yaml
# .gitlab-ci.yml
stages:
  - build
  - release

build:
  stage: build
  script:
    - python build.py -v $CI_COMMIT_TAG -p all -o artifacts -c --no-tag
  artifacts:
    paths:
      - artifacts/
  only:
    - tags

release:
  stage: release
  script:
    - echo "Creating release for $CI_COMMIT_TAG"
  only:
    - tags
```

#### 3. **标签创建工作流**

**传统方式（手动）：**
```bash
# 开发者手动创建和推送标签
git tag v1.0.0
git push origin v1.0.0
# 触发 CI/CD 构建和发布
```

**自动化方式（推荐）：**
```yaml
# GitHub Actions - 自动标签和发布
name: Auto Release

on:
  push:
    branches: [ main ]

jobs:
  auto-release:
    if: "contains(github.event.head_commit.message, '[release]')"
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0  # 获取完整历史用于计算版本号
      
      - name: Generate version
        id: version
        run: |
          # 使用脚本生成版本号
          VERSION=$(python -c "
          import subprocess
          result = subprocess.run(['git', 'rev-list', '--count', 'HEAD'], 
                                capture_output=True, text=True)
          count = int(result.stdout.strip())
          major = count // 1000
          minor = (count % 1000) // 100  
          patch = count % 100
          print(f'{major}.{minor}.{patch}')
          ")
          echo "version=v$VERSION" >> $GITHUB_OUTPUT
          echo "version_number=$VERSION" >> $GITHUB_OUTPUT
      
      - name: Create and push tag
        run: |
          git config user.name github-actions
          git config user.email github-actions@github.com
          git tag ${{ steps.version.outputs.version }}
          git push origin ${{ steps.version.outputs.version }}
      
      - name: Build binaries
        run: |
          python build.py -v ${{ steps.version.outputs.version_number }} -p all -o artifacts -c --no-tag
      
      - name: Create Release
        uses: actions/create-release@v1
        with:
          tag_name: ${{ steps.version.outputs.version }}
          release_name: Release ${{ steps.version.outputs.version }}
          files: artifacts/*
```

#### 4. **分离关注点的好处**

1. **避免权限问题**：CI/CD 系统有统一的 Git 操作权限
2. **确保原子性**：标签创建和构建发布在同一个事务中
3. **支持回滚**：如果构建失败，可以删除标签重新发布
4. **审计追踪**：所有标签操作都有完整的 CI/CD 日志
5. **多环境支持**：不同环境可以有不同的标签策略

#### 5. **混合策略示例**

```bash
# 开发环境：允许本地标签管理
python build.py -v 1.0.0-dev -p current

# 测试环境：使用自动版本，无标签
python build.py -v auto -p all --no-tag

# 生产环境：CI/CD 管理标签，编译使用指定版本
python build.py -v $CI_COMMIT_TAG -p all --no-tag
```

这种分离策略确保了版本管理的一致性和可控性，避免了本地脚本和 CI/CD 系统之间的冲突。

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
python build.py                          # 编译当前平台（自动版本号）
python build.py -p all                   # 编译所有平台（自动版本号）
python build.py -v 1.0.0 -p all -c      # 发布版本编译

# 版本管理
python build.py -v auto                  # 自动版本号
python build.py -v 1.2.3                # 指定版本号
python build.py --no-tag                # 跳过标签管理
python build.py --no-push               # 创建标签但不推送

# 快速编译
python build.py --skip-deps --skip-swagger  # 跳过依赖和文档生成

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

**最后更新**: 2025-11-04
