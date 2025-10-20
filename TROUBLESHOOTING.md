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

### 文档参考

- [Wails v3 文档](https://v3alpha.wails.io/)
- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [GoReleaser 文档](https://goreleaser.com/)
- [NFPM 文档](https://nfpm.goreleaser.com/)

---

**文档版本**: 1.0
**最后更新**: 2025-10-21
**维护者**: 浮浮酱 🐱
**适用项目**: Wails v3 多平台桌面应用

ฅ'ω'ฅ 祝构建顺利喵～
