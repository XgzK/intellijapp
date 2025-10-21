# 多平台构建问题排查与解决方案

本文档记录了在为 Wails v3 项目配置 GitHub Actions 多平台自动构建过程中遇到的所有问题及其解决方案。

---

## 目录

1. [前端构建问题](#1-前端构建问题)
2. [TypeScript 绑定缺失](#2-typescript-绑定缺失)
3. [Linux 依赖库缺失](#3-linux-依赖库缺失)
4. [Linux CGO 链接错误](#4-linux-cgo-链接错误)
5. [Taskfile 命令参数错误](#5-taskfile-命令参数错误)
6. [WebView2 下载失败](#6-webview2-下载失败)
7. [Windows NSIS 路径问题](#7-windows-nsis-路径问题)
8. [Linux 打包文件扩展名错误](#8-linux-打包文件扩展名错误)
9. [MSIX 打包命令不支持](#9-msix-打包命令不支持)
10. [依赖安装顺序问题](#10-依赖安装顺序问题)
11. [GoReleaser 配置错误](#11-goreleaser-配置错误)
12. [GoReleaser Before Hooks 重复构建](#12-goreleaser-before-hooks-重复构建)
13. [GoReleaser Builds 触发 CGO 编译](#13-goreleaser-builds-触发-cgo-编译)
14. [Test Job 缺少 Linux 依赖](#14-test-job-缺少-linux-依赖)

---

## 1. 前端构建问题

### 问题描述
```
pattern all:frontend/dist: no matching files found
```

### 原因分析
- Go 的 `embed` 指令要求 `frontend/dist` 目录必须存在
- CI 环境中未构建前端即尝试编译 Go 代码

### 解决方案

**方案 A：完整构建前端**
```yaml
- name: Install frontend dependencies
  run: npm install --prefix frontend

- name: Build frontend
  run: npm run build --prefix frontend
```

**方案 B：使用占位符文件（测试环境推荐）**
```yaml
- name: Create placeholder frontend dist
  run: |
    mkdir -p frontend/dist
    echo "<!DOCTYPE html><html><body>Test</body></html>" > frontend/dist/index.html
```

### 适用场景
- 方案 A：生产构建
- 方案 B：单元测试、快速验证

---

## 2. TypeScript 绑定缺失

### 问题描述
```
Cannot find module '../../bindings/github.com/XgzK/intellijapp/internal/service/configservice'
```

### 原因分析
- Wails v3 需要生成 TypeScript 绑定供前端调用后端服务
- 绑定文件由 `wails3 generate bindings` 命令生成
- CI 环境中未执行绑定生成步骤

### 解决方案

**生产环境：生成真实绑定**
```yaml
- name: Generate Wails Bindings
  run: wails3 generate bindings
```

**测试环境：创建占位符**
```yaml
- name: Create placeholder bindings
  run: |
    mkdir -p frontend/bindings/github.com/XgzK/intellijapp/internal/service
    echo "// Placeholder for CI testing" > frontend/bindings/github.com/XgzK/intellijapp/internal/service/index.ts
```

### 最佳实践
1. 绑定生成应在前端构建**之前**执行
2. 使用 Wails Task 系统自动化：`task common:generate:bindings`

---

## 3. Linux 依赖库缺失

### 问题描述
```
Package gtk+-3.0 was not found in the pkg-config search path
Package webkit2gtk-4.1 was not found
```

### 原因分析
- Wails v3 在 Linux 上依赖 GTK3 和 WebKit2GTK
- Ubuntu runner 默认未安装这些开发库

### 解决方案

**完整依赖安装命令**
```yaml
- name: Install Linux Dependencies
  if: runner.os == 'Linux'
  run: |
    sudo apt-get update
    sudo apt-get install -y \
      build-essential \
      pkg-config \
      libgtk-3-dev \
      libwebkit2gtk-4.1-dev \
      libjavascriptcoregtk-4.1-dev \
      libglib2.0-dev \
      libpango1.0-dev \
      libcairo2-dev \
      libgdk-pixbuf-2.0-dev \
      libsoup-3.0-dev \
      libharfbuzz-dev \
      libatk1.0-dev \
      nsis \
      rpm
```

### 依赖库说明
| 库名 | 用途 |
|------|------|
| libgtk-3-dev | GTK3 UI 框架 |
| libwebkit2gtk-4.1-dev | WebView 渲染引擎 |
| libjavascriptcoregtk-4.1-dev | JavaScript 引擎 |
| libsoup-3.0-dev | HTTP 客户端库 |
| nsis, rpm | 打包工具 |

### 注意事项
- 使用 WebKit2GTK **4.1** 版本（不是 4.0）
- GTK4 目前不被 Wails v3 支持

---

## 4. Linux CGO 链接错误

### 问题描述
```
/usr/bin/ld: cannot find -lwebkit2gtk-4.1: No such file or directory
/usr/bin/ld: cannot find -lgtk-3: No such file or directory
/usr/bin/ld: cannot find -lpangocairo-1.0: No such file or directory
```

### 原因分析
- Wails3 CLI 本身是 CGO 程序，需要链接 C 库
- 即使安装了 `-dev` 包，链接器也可能找不到库文件
- `PKG_CONFIG_PATH` 未正确配置

### 解决方案

**为 Linux 单独配置 Wails CLI 安装**
```yaml
- name: Install Wails CLI (Linux)
  if: runner.os == 'Linux'
  run: |
    # 设置 CGO 标志帮助链接器找到库
    export CGO_ENABLED=1
    export PKG_CONFIG_PATH=/usr/lib/x86_64-linux-gnu/pkgconfig:$PKG_CONFIG_PATH
    go install github.com/wailsapp/wails/v3/cmd/wails3@latest

- name: Install Wails CLI (Non-Linux)
  if: runner.os != 'Linux'
  run: go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

### 技术要点
1. **CGO_ENABLED=1**：显式启用 CGO 编译
2. **PKG_CONFIG_PATH**：指向库的 `.pc` 配置文件路径
3. **分平台处理**：Linux 需要特殊配置，其他平台无需

### 为什么 Wails3 CLI 需要 CGO？
- Windows: 使用 WebView2，不需要 CGO
- macOS: 使用系统 WebKit，需要 CGO
- **Linux: 使用 WebKit2GTK，强制需要 CGO**

---

## 5. Taskfile 命令参数错误

### 问题描述
```
flag provided but not defined: -ts
Usage: wails3 generate icons [flags]
```

### 原因分析
- `-ts` 标志仅适用于 `wails3 generate bindings` 命令
- 错误地将其用于 `wails3 generate icons` 命令

### 解决方案

**错误示例**
```yaml
# ❌ 错误
cmds:
  - wails3 generate -ts icons
```

**正确示例**
```yaml
# ✅ 正确
generate:icons:
  cmds:
    - wails3 generate icons -input appicon.png

generate:bindings:
  cmds:
    - wails3 generate bindings -ts  # -ts 只能用于 bindings
```

### 常用 Wails3 生成命令
```bash
wails3 generate bindings        # 生成 Go ↔ TS 绑定
wails3 generate bindings -ts    # 同时生成 TypeScript 声明
wails3 generate icons           # 生成应用图标
wails3 generate syso            # 生成 Windows 资源文件
```

---

## 6. WebView2 下载失败

### 问题描述
```
Response status code does not indicate success: 404 (Not Found)
https://go.microsoft.com/fwlink/p/?LinkId=2124703&platform=x64
```

### 原因分析
- URL 中包含不必要的 `&platform=x64` 参数
- Microsoft 下载链接已更新

### 解决方案

**修正后的 WebView2 安装脚本**
```yaml
- name: Install WebView2 Runtime
  if: runner.os == 'Windows'
  shell: pwsh
  run: |
    # Windows runners 通常预装了 WebView2
    # 如果未安装，下载 x64 安装器
    $installer = 'MicrosoftEdgeWebView2RuntimeInstallerX64.exe'
    $downloadUrl = "https://go.microsoft.com/fwlink/p/?LinkId=2124703"  # 移除 platform 参数
    $installerPath = Join-Path $env:TEMP $installer
    try {
      Invoke-WebRequest -Uri $downloadUrl -OutFile $installerPath -UseBasicParsing -ErrorAction Stop
      Start-Process -FilePath $installerPath -ArgumentList '/silent','/install' -NoNewWindow -Wait
      Remove-Item $installerPath -Force
    } catch {
      Write-Host "WebView2 可能已安装或下载失败: $_"
    }
```

### 最佳实践
1. 添加错误处理避免阻塞构建
2. GitHub Windows runner 通常已预装 WebView2
3. 使用官方稳定下载链接

---

## 7. Windows NSIS 路径问题

### 问题描述
```
File: "D:\a\intellijapp\intellijapp/bin/intellijapp.exe" -> no files found
```

### 原因分析
- Task 的 `{{.ROOT_DIR}}` 变量在 Windows 上使用反斜杠 `\`
- 与手动添加的正斜杠 `/` 混合，导致路径无效
- NSIS 对路径格式敏感

### 解决方案

**方案 A：使用相对路径（推荐）**
```yaml
# ✅ 推荐：使用相对路径和统一的反斜杠
cmds:
  - makensis -DARG_WAILS_AMD64_BINARY="..\..\..\bin\intellijapp.exe" project.nsi
```

**方案 B：使用绝对路径（不推荐）**
```yaml
# ⚠️ 不推荐：路径分隔符可能混合
cmds:
  - makensis -DARG_WAILS_AMD64_BINARY="{{.ROOT_DIR}}/bin/intellijapp.exe" project.nsi
```

### 路径规则
| 平台 | 分隔符 | 示例 |
|------|--------|------|
| Windows | `\` | `build\windows\bin\app.exe` |
| Linux/macOS | `/` | `build/linux/bin/app` |
| 相对路径 | `..\` (Win) | `..\..\..\bin\app.exe` |

### 调试技巧
```yaml
# 添加调试步骤查看实际路径
- name: Debug paths
  run: |
    echo "ROOT_DIR: {{.ROOT_DIR}}"
    echo "BIN_DIR: {{.BIN_DIR}}"
    dir "{{.ROOT_DIR}}\{{.BIN_DIR}}"
```

---

## 8. Linux 打包文件扩展名错误

### 问题描述
```
error creating package: matching "./bin/intellijapp.exe": file does not exist
```

### 原因分析
- nfpm 配置文件从 Windows 模板复制而来
- 错误地在 Linux 配置中使用了 `.exe` 扩展名
- Linux 可执行文件没有扩展名

### 解决方案

**修复 `build/linux/nfpm/nfpm.yaml`**
```yaml
# ❌ 错误
name: "intellijapp.exe"
contents:
  - src: "./bin/intellijapp.exe"
    dst: "/usr/local/bin/intellijapp.exe"

# ✅ 正确
name: "intellijapp"
contents:
  - src: "./bin/intellijapp"
    dst: "/usr/local/bin/intellijapp"
  - src: "./build/appicon.png"
    dst: "/usr/share/icons/hicolor/128x128/apps/intellijapp.png"
  - src: "./build/linux/intellijapp.desktop"
    dst: "/usr/share/applications/intellijapp.desktop"
```

**修复 Desktop 文件 `build/linux/intellijapp.desktop`**
```ini
# ❌ 错误
[Desktop Entry]
Exec=/usr/local/bin/intellijapp.exe %u
Icon=intellijapp.exe
StartupWMClass=intellijapp.exe

# ✅ 正确
[Desktop Entry]
Exec=/usr/local/bin/intellijapp %u
Icon=intellijapp
StartupWMClass=intellijapp
```

**文件重命名**
```bash
# Desktop 文件也需要重命名
mv build/linux/desktop build/linux/intellijapp.desktop
```

### 跨平台二进制命名规范
| 平台 | 扩展名 | 示例 |
|------|--------|------|
| Windows | `.exe` | `intellijapp.exe` |
| Linux | 无 | `intellijapp` |
| macOS | 无 | `intellijapp` |

---

## 9. MSIX 打包命令不支持

### 问题描述
```
Error: flag provided but not defined: -name
task: [windows:create:msix:package] wails3 tool msix --name "intellijapp" ...
```

### 原因分析
- Wails v3 当前版本不支持完整的 `wails3 tool msix` 命令
- MSIX 打包功能可能未完全实现或参数格式不同
- CI 环境中无法测试 MSIX 工具的正确参数

### 解决方案

**方案 A：禁用 MSIX 打包（推荐）**
```yaml
# CI Workflow
matrix:
  include:
    - os: windows-latest
      platform: windows
      formats: "nsis"  # 仅使用 NSIS，移除 msix

# 注释掉 MSIX 步骤
# - name: Package Windows MSIX
#   if: runner.os == 'Windows' && contains(matrix.formats, 'msix')
#   run: task windows:create:msix:package
```

**方案 B：等待 Wails v3 正式版**
- MSIX 支持可能在未来的 Wails v3 版本中完善
- 当前使用 NSIS 安装程序已足够

### NSIS vs MSIX 对比
| 特性 | NSIS | MSIX |
|------|------|------|
| 支持版本 | Windows 7+ | Windows 10+ |
| 签名要求 | 可选 | 必须 |
| Microsoft Store | ❌ | ✅ |
| 当前 Wails v3 | ✅ 支持 | ⚠️ 部分支持 |

### 最佳实践
1. 生产环境优先使用 NSIS
2. MSIX 可用于 Microsoft Store 分发（需签名）
3. 等待 Wails v3 稳定版后重新评估

---

## 10. 依赖安装顺序问题

### 问题描述
即使安装了所有 Linux 依赖，Wails CLI 安装仍然失败

### 原因分析
- Wails CLI 安装步骤在**所有平台**上统一执行
- Linux 平台上，Wails CLI 本身需要编译（CGO）
- 如果依赖未提前安装，编译会失败

### 解决方案

**正确的步骤顺序**
```yaml
# ✅ 正确顺序
steps:
  # 1. 安装系统级依赖
  - name: Install Linux Dependencies
    if: runner.os == 'Linux'
    run: sudo apt-get install -y libgtk-3-dev ...

  # 2. 安装 Wails CLI（依赖已就绪）
  - name: Install Wails CLI (Linux)
    if: runner.os == 'Linux'
    run: |
      export CGO_ENABLED=1
      export PKG_CONFIG_PATH=/usr/lib/x86_64-linux-gnu/pkgconfig:$PKG_CONFIG_PATH
      go install github.com/wailsapp/wails/v3/cmd/wails3@latest

  # 3. 构建应用
  - name: Build Application
    run: task build
```

**错误的顺序示例**
```yaml
# ❌ 错误：Wails CLI 先于依赖安装
steps:
  - name: Install Wails CLI  # 此时依赖未安装，Linux 上会失败
    run: go install github.com/wailsapp/wails/v3/cmd/wails3@latest

  - name: Install Linux Dependencies
    if: runner.os == 'Linux'
    run: sudo apt-get install ...
```

### 依赖关系图
```
Linux Dependencies (GTK, WebKit)
        ↓
    Wails CLI (CGO 编译)
        ↓
  Application Build
```

---

## 总结与最佳实践

### CI/CD 配置要点

1. **测试与打包分离**
   ```yaml
   jobs:
     test:      # 快速反馈（使用占位符）
     package:   # 完整构建（真实资源）
     release:   # 发布到 GitHub
   ```

2. **平台特定处理**
   - Linux: 需要安装大量依赖，配置 CGO
   - Windows: 需要 WebView2，路径使用反斜杠
   - macOS: 相对简单，系统自带 WebKit

3. **依赖安装顺序**
   ```
   系统依赖 → Wails CLI → Task CLI → 构建应用
   ```

4. **错误处理**
   - 使用 `continue-on-error` 处理可选步骤
   - 添加调试输出便于排查
   - 保留详细的日志信息

### 常用调试命令

```bash
# 检查依赖库
pkg-config --list-all | grep gtk
pkg-config --modversion gtk+-3.0
pkg-config --cflags --libs webkit2gtk-4.1

# 检查链接器路径
echo $PKG_CONFIG_PATH
echo $LD_LIBRARY_PATH

# 手动测试 Wails CLI 安装
CGO_ENABLED=1 go install github.com/wailsapp/wails/v3/cmd/wails3@latest

# 检查生成的文件
ls -la bin/
ls -la frontend/dist/
ls -la frontend/bindings/
```

---

## 11. GoReleaser 配置错误

### 问题描述
```
yaml: unmarshal errors:
  line 75: field extra_files not found in type config.Project
```

### 原因分析
- GoReleaser 配置文件结构不正确
- `extra_files` 字段被放在了顶层，而非 `release` 节点下
- GoReleaser v1.26+ 要求 `extra_files` 必须是 `release` 配置的子项

### 解决方案

**错误的配置：**
```yaml
archives:
  - format: binary

# ❌ 错误：extra_files 在顶层
extra_files:
  - glob: ./bin/*-installer.exe

release:
  draft: true
```

**正确的配置：**
```yaml
archives:
  - format: binary

release:
  draft: true
  # ✅ 正确：extra_files 在 release 节点下
  extra_files:
    - glob: ./bin/*-installer.exe
    - glob: ./bin/*.AppImage
    - glob: ./bin/*.deb
    - glob: ./bin/*.rpm
    - glob: ./bin/*.dmg
    - glob: ./bin/*.pkg
```

### 最佳实践
1. 参考最新的 GoReleaser 官方文档：https://goreleaser.com/customization/release/
2. 使用 `goreleaser check` 命令验证配置文件语法
3. 注意 GoReleaser 版本升级可能带来的配置变更

---

## 12. GoReleaser Before Hooks 重复构建

### 问题描述
```
Run goreleaser/goreleaser-action@v5
  building       binaries=0 builds=2
  running        before hooks
error=hook failed: shell: 'npm run build --prefix frontend': exit status 127:
sh: 1: vue-tsc: not found
Error: The process '/opt/hostedtoolcache/goreleaser-action/1.26.2/x64/goreleaser' failed with exit code 1
```

### 原因分析
1. **GoReleaser 运行在 release job 中**
   - release job 只负责从 package job 下载已构建的包
   - release job 不需要重新构建前端或生成绑定

2. **before hooks 设计用于本地开发**
   - 本地运行 `goreleaser release` 时需要从源代码构建
   - CI/CD 中 package job 已经完成了所有构建工作

3. **Node 依赖未安装**
   - release job 只安装了 Go 和 Node 环境
   - 未执行 `npm install`，导致 `vue-tsc` 找不到

### 解决方案

**修改 `.goreleaser.yaml`，禁用 before hooks：**

```yaml
project_name: intellijapp

# IMPORTANT: In CI/CD, all packages are pre-built by the package job.
# GoReleaser only needs to create the GitHub release and upload artifacts.
# Therefore, before hooks are disabled to avoid rebuilding from scratch.
#
# If running GoReleaser locally for development, uncomment these hooks:
# before:
#   hooks:
#     - npm install --prefix frontend
#     - npm run build --prefix frontend
#     - wails3 generate bindings -f "-tags production" -clean=true -ts
#     - wails3 generate -ts icons -input build/appicon.png -macfilename build/darwin/icons.icns -windowsfilename build/windows/icon.ico

builds:
  - id: unix
    # ... rest of config
```

### CI/CD 流程说明

**完整的 CI/CD 流程：**

```
┌─────────────────┐
│   test job      │ ← 运行单元测试（使用占位符）
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  package job    │ ← 完整构建（前端 + 绑定 + 打包）
│  (3 platforms)  │   - npm install & build
└────────┬────────┘   - wails3 generate bindings
         │            - task package
         │            生成所有平台的安装包
         ▼
┌─────────────────┐
│  release job    │ ← 仅创建 GitHub Release
│  (downloads)    │   - 下载所有 artifacts
└─────────────────┘   - 运行 GoReleaser（无需构建）
                      - 上传到 GitHub Releases
```

### 关键理解
- **package job 的输出：** bin/ 目录中的所有安装包（.exe, .AppImage, .deb, .rpm, .dmg, .pkg）
- **release job 的职责：** 仅收集和发布，不重新构建
- **GoReleaser 的角色：** 创建 release、生成 changelog、上传 artifacts

### 最佳实践

1. **分离构建和发布逻辑**
   ```yaml
   # package job: 负责构建
   - name: Build Application
     run: task build

   # release job: 负责发布
   - name: Release with GoReleaser
     uses: goreleaser/goreleaser-action@v5
   ```

2. **本地开发时的用法**
   - 取消注释 `.goreleaser.yaml` 中的 before hooks
   - 运行 `goreleaser release --snapshot --clean` 进行本地测试

3. **CI/CD 中的用法**
   - 保持 before hooks 禁用
   - 确保 package job 生成所有需要的文件
   - 使用 `extra_files` 上传 package job 的输出

### 相关原则
- **KISS (简单至上)：** 每个 job 只做一件事
- **DRY (杜绝重复)：** 不在 release job 中重复 package job 的工作
- **YAGNI (精益求精)：** GoReleaser 只负责发布，不负责构建

---

## 13. GoReleaser Builds 触发 CGO 编译

### 问题描述
```
# github.com/wailsapp/wails/v3/internal/operatingsystem
# [pkg-config --cflags  -- gtk+-3.0 webkit2gtk-4.1]
Package gtk+-3.0 was not found in the pkg-config search path.
Perhaps you should add the directory containing `gtk+-3.0.pc'
to the PKG_CONFIG_PATH environment variable
Package 'gtk+-3.0', required by 'virtual:world', not found
Package 'webkit2gtk-4.1', required by 'virtual:world', not found
```

### 原因分析

1. **GoReleaser 的 builds 配置会编译二进制文件**
   - 即使 package job 已经构建了所有平台的包
   - GoReleaser 仍会根据 `builds` 配置重新编译

2. **Wails 项目依赖 CGO**
   - 即使设置 `CGO_ENABLED=0`，Wails 在某些平台仍需要 CGO
   - Linux 构建需要 GTK3 和 WebKit2GTK 依赖

3. **Release job 环境缺少依赖**
   - Release job 运行在 ubuntu-latest
   - 没有安装 Linux 构建所需的系统依赖库
   - 导致编译失败

### 错误理解 vs 正确理解

**❌ 错误理解（之前的假设）：**
```
GoReleaser 的作用：
1. 下载 package job 的 artifacts
2. 创建 GitHub Release
3. 上传预构建的文件
4. 不会重新编译任何代码
```

**✅ 正确理解（实际行为）：**
```
GoReleaser 的作用：
1. 如果有 builds 配置，会重新编译所有二进制文件
2. 根据 archives 配置打包二进制文件
3. 创建 GitHub Release
4. 上传编译/打包的文件 + extra_files
```

### 解决方案

**完全禁用 GoReleaser 的构建功能，仅用于发布管理：**

```yaml
# GoReleaser configuration for intellijapp.
#
# IMPORTANT: This configuration is optimized for CI/CD pipelines where
# all platform-specific packages are pre-built by the package job.
# GoReleaser's role is ONLY to create the GitHub release and upload artifacts.

project_name: intellijapp

# Disable builds entirely - all binaries are pre-built by package job
builds:
  - skip: true

# Disable archives since we're uploading installer packages directly
archives:
  - id: skip-archives
    format: binary

# Skip checksum generation (optional)
checksum:
  disable: true

changelog:
  use: git
  filters:
    exclude:
      - '^docs?:'
      - '^ci:'

release:
  draft: true
  # Upload all pre-built packages from package job
  extra_files:
    - glob: ./bin/*-installer.exe  # Windows NSIS installers
    - glob: ./bin/*.AppImage       # Linux AppImage
    - glob: ./bin/*.deb            # Debian packages
    - glob: ./bin/*.rpm            # RedHat packages
    - glob: ./bin/*.dmg            # macOS disk images
    - glob: ./bin/*.pkg            # macOS installer packages
```

### 架构设计说明

**CI/CD 职责分离：**

```
┌──────────────────────────────────────────────┐
│           Package Job (3 platforms)           │
│  ┌────────────┐ ┌──────────┐ ┌─────────────┐ │
│  │  Windows   │ │  macOS   │ │    Linux    │ │
│  │  runner    │ │  runner  │ │   runner    │ │
│  └─────┬──────┘ └────┬─────┘ └──────┬──────┘ │
│        │             │               │        │
│   Build + Package    │          Build + Package
│        │             │               │        │
│  ┌─────▼──────┐ ┌───▼──────┐ ┌──────▼──────┐ │
│  │ .exe       │ │ .app     │ │ .AppImage   │ │
│  │ -installer │ │ .dmg     │ │ .deb        │ │
│  │            │ │ .pkg     │ │ .rpm        │ │
│  └────────────┘ └──────────┘ └─────────────┘ │
│        │             │               │        │
│        └─────────────┼───────────────┘        │
│                      │                        │
│                 Upload to                     │
│              GitHub Artifacts                 │
└──────────────────────┼───────────────────────┘
                       │
                       ▼
┌──────────────────────────────────────────────┐
│              Release Job (Linux)              │
│                                               │
│  1. Download all artifacts from package job  │
│  2. Run GoReleaser (NO building, NO hooks)   │
│     - Create GitHub Release                  │
│     - Generate changelog                     │
│     - Upload files via extra_files           │
│                                               │
└──────────────────────────────────────────────┘
```

### 关键要点

1. **GoReleaser 仅用于 Release 管理**
   - 创建 GitHub Release
   - 生成 changelog
   - 上传文件（通过 `extra_files`）

2. **所有构建工作在 Package Job 完成**
   - 每个平台使用原生 runner
   - 使用 Task/Wails 工具链构建
   - 生成平台特定的安装包

3. **避免重复构建**
   - Package job 已经构建了所有内容
   - Release job 不应该重新编译
   - 符合 DRY 原则

### 最佳实践

1. **构建与发布分离**
   ```yaml
   # Package job: 专注构建
   jobs:
     package:
       strategy:
         matrix:
           os: [windows-latest, macos-latest, ubuntu-latest]
       steps:
         - name: Build and Package
           run: task build && task package

   # Release job: 专注发布
   jobs:
     release:
       steps:
         - name: Download artifacts
           uses: actions/download-artifact@v4
         - name: Create release
           uses: goreleaser/goreleaser-action@v5
   ```

2. **GoReleaser 配置最小化**
   - 禁用不需要的功能（builds, archives, checksums）
   - 只保留必要的配置（changelog, release, extra_files）
   - 添加清晰的注释说明设计意图

3. **优势总结**
   - ✅ **性能：** 避免重复编译，节省 CI/CD 时间
   - ✅ **可靠：** 使用各平台原生环境构建，兼容性更好
   - ✅ **简单：** 职责清晰，易于理解和维护
   - ✅ **灵活：** 可以使用 Wails 专用工具链（Task, wails3 package）

### 相关原则
- **KISS (简单至上)：** GoReleaser 只做发布，不做构建
- **DRY (杜绝重复)：** 构建一次，发布一次
- **单一职责：** Package job 负责构建，Release job 负责发布
- **关注点分离：** 构建逻辑与发布逻辑完全隔离

---

## 14. Test Job 缺少 Linux 依赖

### 问题描述
```
# github.com/wailsapp/wails/v3/internal/operatingsystem
# [pkg-config --cflags  -- gtk+-3.0 webkit2gtk-4.1]
Package gtk+-3.0 was not found in the pkg-config search path.
Perhaps you should add the directory containing `gtk+-3.0.pc'
to the PKG_CONFIG_PATH environment variable
Package 'gtk+-3.0', required by 'virtual:world', not found
Package 'webkit2gtk-4.1', required by 'virtual:world', not found
```

**错误发生在：** Unit Tests job 的 `Run go test` 步骤

### 原因分析

1. **Test job 运行在 ubuntu-latest**
   - 需要测试导入了 Wails 包的代码
   - Wails 依赖 CGO 和 Linux 系统库（GTK3, WebKit2GTK）

2. **Test job 没有安装 Linux 依赖**
   - 原始设计使用占位符文件来避免真实构建
   - 但是 `go test ./...` 仍然会编译和导入 Wails 的真实代码

3. **编译测试代码时触发 CGO**
   - `go test` 需要编译测试文件
   - 测试文件导入了 Wails 包
   - Wails 包需要链接 GTK/WebKit 库
   - 缺少依赖导致编译失败

### 错误理解 vs 正确理解

**❌ 错误理解：**
```
Test job 的设计思路：
- 创建占位符 bindings 和 frontend/dist
- 这样 go test 就不会尝试编译真实的 Wails 代码
- 不需要安装 Linux 依赖
```

**✅ 正确理解：**
```
实际情况：
- 占位符只是避免 embed 指令找不到文件
- go test 仍然会导入和编译 main.go 等文件
- main.go 导入了 wails3 包
- wails3 包依赖 CGO 和系统库
- 必须安装 Linux 依赖才能编译成功
```

### 解决方案

**在 test job 中添加 Linux 依赖安装步骤：**

```yaml
jobs:
  test:
    name: Unit Tests
    runs-on: ubuntu-latest

    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: stable
          check-latest: true

      - name: Set up Node
        uses: actions/setup-node@v4
        with:
          node-version: "lts/*"
          check-latest: true
          cache: "npm"
          cache-dependency-path: frontend/package-lock.json

      # 添加这一步 - 安装 Linux 依赖
      - name: Install Linux Dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y \
            build-essential \
            pkg-config \
            libgtk-3-dev \
            libwebkit2gtk-4.1-dev \
            libjavascriptcoregtk-4.1-dev \
            libglib2.0-dev \
            libpango1.0-dev \
            libcairo2-dev \
            libgdk-pixbuf-2.0-dev \
            libsoup-3.0-dev \
            libharfbuzz-dev \
            libatk1.0-dev

      - name: Install frontend dependencies
        run: npm install --prefix frontend

      - name: Create placeholder bindings (for compilation)
        run: |
          mkdir -p frontend/bindings/github.com/XgzK/intellijapp/internal/service
          echo "// Placeholder for CI testing" > frontend/bindings/github.com/XgzK/intellijapp/internal/service/index.ts

      - name: Create placeholder frontend dist (for embed)
        run: |
          mkdir -p frontend/dist
          echo "<!DOCTYPE html><html><body>Test</body></html>" > frontend/dist/index.html

      - name: Run go test
        run: go test ./...
```

### 关键理解

1. **占位符的作用有限**
   - 占位符只能避免 `embed` 指令报错
   - 无法避免 `import` 语句触发的依赖编译

2. **go test 的编译行为**
   - `go test ./...` 会编译所有测试包
   - 编译过程会递归处理所有 import
   - Wails 包的 import 会触发 CGO 编译

3. **Test 环境需要与 Build 环境一致**
   - Test job 需要安装与 Package job 相同的依赖
   - 这样才能确保测试环境的真实性

### 最佳实践

1. **统一依赖安装**
   ```yaml
   # 可以考虑创建可复用的 action
   # .github/actions/setup-linux-deps/action.yml
   - name: Install Linux Dependencies
     uses: ./.github/actions/setup-linux-deps
   ```

2. **环境一致性**
   - Test job 应该使用与实际构建相同的环境
   - 避免"测试通过但构建失败"的情况

3. **依赖清单文档**
   - 在 README 或文档中列出所有系统依赖
   - 便于本地开发环境配置

### 替代方案

如果不想在 test job 中安装完整的 Linux 依赖，可以考虑：

**方案 B：跳过需要 CGO 的测试**
```yaml
- name: Run go test
  run: go test -tags=!cgo ./...
  env:
    CGO_ENABLED: 0
```

**方案 C：只在 package job 中运行测试**
```yaml
# 移除独立的 test job
# 在 package job 的构建前运行测试
- name: Run tests
  run: go test ./...
- name: Build Application
  run: task build
```

但这些方案都有缺点：
- 方案 B：无法测试完整功能
- 方案 C：测试失败会浪费多平台构建资源

因此推荐 **方案 A（当前方案）**：在 test job 中安装依赖

### 相关原则
- **环境一致性：** 测试环境应该尽可能接近生产环境
- **快速失败：** Test job 先运行，尽早发现问题
- **DRY：** 可以考虑提取依赖安装步骤为可复用 action

---

### 文档参考

- [Wails v3 文档](https://v3alpha.wails.io/)
- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [GoReleaser 文档](https://goreleaser.com/)
- [NFPM 文档](https://nfpm.goreleaser.com/)

---

**文档版本**: 1.3
**最后更新**: 2025-10-21
**问题总数**: 14 个
**维护者**: 浮浮酱 🐱
**适用项目**: Wails v3 多平台桌面应用

ฅ'ω'ฅ 祝构建顺利喵～
