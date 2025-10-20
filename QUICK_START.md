# 🚀 IntellijApp 快速开始

> 快速参考卡片 - 开发、构建和发布的完整命令

---

## 📦 安装依赖

### 一键检查
```bash
wails3 doctor
```

### Windows
```bash
# 安装 NSIS（用于创建安装程序）
choco install nsis -y
```

### macOS
```bash
# 安装 Xcode Command Line Tools
xcode-select --install
```

### Linux (Ubuntu/Debian)
```bash
sudo apt-get update
sudo apt-get install -y build-essential libgtk-3-dev libwebkit2gtk-4.1-dev nsis rpm
```

---

## 🛠️ 开发命令

```bash
# 开发模式（热重载）
wails3 dev

# 生产构建
wails3 build

# 运行构建后的应用
./bin/intellijapp        # macOS/Linux
bin\intellijapp.exe      # Windows
```

---

## 📦 打包命令

### Windows

```bash
# NSIS 安装程序（推荐）
wails3 task windows:create:nsis:installer

# MSIX 包（Microsoft Store）
wails3 task windows:create:msix:package

# 快速打包（默认格式）
wails3 package
```

**输出：**
- `bin/intellijapp-amd64-installer.exe` (NSIS)
- `bin/intellijapp-amd64.msix` (MSIX)

---

### macOS

```bash
# .app 应用包
wails3 package
```

**输出：**
- `bin/intellijapp.app`

---

### Linux

```bash
# 创建所有格式
wails3 package

# 或单独创建
wails3 task linux:create:appimage  # AppImage（推荐）
wails3 task linux:create:deb       # Debian/Ubuntu
wails3 task linux:create:rpm       # Red Hat/Fedora
wails3 task linux:create:aur       # Arch Linux
```

**输出：**
- `bin/intellijapp-*.AppImage`
- `bin/intellijapp_*.deb`
- `bin/intellijapp-*.rpm`
- `bin/intellijapp-*.pkg.tar.zst`

---

## 🚀 发布流程

### 本地测试发布
```bash
# 测试 GoReleaser 配置（不上传）
goreleaser release --snapshot --clean

# 查看生成的文件
ls -lh bin/
```

### GitHub 自动发布
```bash
# 1. 创建并推送标签
git tag v1.0.0
git push --tags

# 2. GitHub Actions 自动执行：
#    ✅ 运行测试
#    ✅ 多平台打包
#    ✅ 创建 GitHub Release（草稿）

# 3. 在 GitHub 上审核并发布 Release
```

---

## 🎯 常用任务

### 清理构建
```bash
# 清理 bin 目录
rm -rf bin/

# 清理 Node 依赖
rm -rf frontend/node_modules/

# 重新安装依赖
npm install --prefix frontend
```

### 更新依赖
```bash
# 更新 Go 依赖
go get -u ./...
go mod tidy

# 更新前端依赖
npm update --prefix frontend

# 更新 Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

### 重新生成绑定
```bash
# 重新生成 TypeScript 绑定
wails3 generate bindings -f "-tags production" -clean=true -ts

# 重新生成图标
wails3 generate icons -input build/appicon.png \
  -macfilename build/darwin/icons.icns \
  -windowsfilename build/windows/icon.ico
```

---

## 🐛 故障排查

### 检查系统
```bash
# 完整诊断
wails3 doctor

# 检查 Go 环境
go version
go env

# 检查 Node 环境
node --version
npm --version
```

### 清理并重建
```bash
# 清理所有构建产物
rm -rf bin/ frontend/dist/ frontend/node_modules/

# 重新安装并构建
npm install --prefix frontend
wails3 build
```

### Windows 特定
```bash
# 检查 NSIS
where makensis

# 检查 WebView2
wails3 doctor
```

### macOS 特定
```bash
# 检查 Xcode
xcode-select -p

# 重新安装 Xcode Command Line Tools
xcode-select --install
```

### Linux 特定
```bash
# 检查依赖
pkg-config --modversion gtk+-3.0
pkg-config --modversion webkit2gtk-4.1

# 检查 ldd
ldd bin/intellijapp
```

---

## 📚 详细文档

- [PACKAGING.md](./PACKAGING.md) - 完整打包指南
- [README.md](./README.md) - 项目说明
- [Wails v3 官方文档](https://v3alpha.wails.io)

---

## 🎨 文件结构速查

```
intellijapp/
├── bin/                    # 构建输出目录
├── build/                  # 构建配置和资源
│   ├── appicon.png        # 应用图标
│   ├── config.yml         # Wails 配置
│   ├── darwin/            # macOS 特定文件
│   ├── linux/             # Linux 特定文件
│   └── windows/           # Windows 特定文件
├── frontend/              # 前端代码
│   ├── index.html
│   ├── main.js
│   └── package.json
├── .github/workflows/     # CI/CD 配置
│   └── release.yml        # 自动发布工作流
├── .goreleaser.yaml       # GoReleaser 配置
├── Taskfile.yml           # Task 构建任务
├── go.mod                 # Go 模块定义
├── main.go                # 主程序入口
└── greetservice.go        # 示例服务
```

---

## ⚡ 快捷命令别名（可选）

可以在 `~/.bashrc` 或 `~/.zshrc` 中添加：

```bash
# Wails 别名
alias wdev='wails3 dev'
alias wbuild='wails3 build'
alias wpack='wails3 package'
alias wdoc='wails3 doctor'

# 构建别名
alias tbuild='task build'
alias tpack='task package'

# Git 发布
alias grelease='git tag -a "$1" -m "Release $1" && git push --tags'
```

使用示例：
```bash
wdev           # 开发模式
wbuild         # 生产构建
wpack          # 打包
grelease v1.0.0  # 创建发布标签
```

---

**祝开发愉快！** 🎉
