# 📦 IntellijApp 打包指南

本文档详细说明如何为不同平台构建和打包 IntellijApp。

## 🎯 快速开始

### 基础构建
```bash
# 开发模式（支持热重载）
wails3 dev

# 生产构建
wails3 build
```

### 快速打包（推荐）
```bash
# 当前平台默认格式
wails3 package
```

---

## 🪟 Windows 打包

### NSIS 安装程序（推荐）
```bash
# 使用默认架构（当前系统架构）
wails3 task windows:create:nsis:installer

# 指定架构
wails3 task windows:create:nsis:installer ARCH=amd64
wails3 task windows:create:nsis:installer ARCH=arm64
```

**输出文件：**
- `bin/intellijapp-{arch}-installer.exe`

**特点：**
- 标准 Windows 安装程序
- 自动处理 WebView2 运行时
- 支持卸载程序
- 适合大多数 Windows 用户

### MSIX 包（Microsoft Store）
```bash
# 使用默认架构
wails3 task windows:create:msix:package

# 指定架构和证书
wails3 task windows:create:msix:package ARCH=amd64 CERT_PATH=/path/to/cert.pfx
```

**输出文件：**
- `bin/intellijapp-{arch}.msix`

**特点：**
- Microsoft Store 兼容
- 需要代码签名证书
- 沙箱化安全环境

### 安装 MSIX 工具（首次使用）
```bash
wails3 task windows:install:msix:tools
```

---

## 🍎 macOS 打包

### .app 应用包
```bash
# 默认架构
wails3 package

# 生产构建
wails3 task darwin:build PRODUCTION=true
```

**输出文件：**
- `bin/intellijapp.app`

**支持的架构：**
- `amd64` - Intel 芯片
- `arm64` - Apple Silicon (M1/M2/M3)

**特点：**
- 原生 macOS 应用格式
- 拖拽安装到 Applications
- 支持通用二进制（Universal Binary）

### DMG 磁盘镜像（如果配置）
如果项目配置了 DMG 打包，可以创建：
```bash
# 需要额外的 DMG 构建工具
wails3 task darwin:create:dmg
```

---

## 🐧 Linux 打包

### 全部格式（一键打包）
```bash
# 创建所有 Linux 格式
wails3 package
```

这会生成：
- AppImage
- .deb（Debian/Ubuntu）
- .rpm（Red Hat/Fedora/CentOS）
- Arch Linux 包

### AppImage（推荐 - 通用格式）
```bash
wails3 task linux:create:appimage
```

**输出文件：**
- `bin/intellijapp-{version}-{arch}.AppImage`

**特点：**
- 单文件，开箱即用
- 兼容大多数 Linux 发行版
- 无需安装，直接运行
- 支持自动更新

### Debian/Ubuntu 包（.deb）
```bash
wails3 task linux:create:deb
```

**输出文件：**
- `bin/intellijapp_{version}_{arch}.deb`

**安装方式：**
```bash
sudo dpkg -i intellijapp_*.deb
sudo apt-get install -f  # 解决依赖
```

**支持的发行版：**
- Debian 10+
- Ubuntu 20.04+
- Linux Mint 20+

### Red Hat 包（.rpm）
```bash
wails3 task linux:create:rpm
```

**输出文件：**
- `bin/intellijapp-{version}-{arch}.rpm`

**安装方式：**
```bash
sudo rpm -ivh intellijapp-*.rpm
# 或使用 dnf/yum
sudo dnf install intellijapp-*.rpm
```

**支持的发行版：**
- Red Hat Enterprise Linux 8+
- Fedora 33+
- CentOS 8+

### Arch Linux 包
```bash
wails3 task linux:create:aur
```

**输出文件：**
- `bin/intellijapp-{version}-{arch}.pkg.tar.zst`

**安装方式：**
```bash
sudo pacman -U intellijapp-*.pkg.tar.zst
```

---

## 🚀 GoReleaser 集成

项目已配置 GoReleaser 用于自动化多平台发布。

### 手动发布
```bash
# 创建标签
git tag v1.0.0
git push --tags

# 本地测试发布（不上传）
goreleaser release --snapshot --clean

# 正式发布到 GitHub
goreleaser release --clean
```

### GitHub Actions 自动发布
推送标签会自动触发 CI/CD 流程：
```bash
git tag v1.0.0
git push --tags
```

CI 会自动：
1. ✅ 运行测试
2. ✅ 多平台构建（Linux, macOS, Windows）
3. ✅ 收集所有二进制文件和打包格式
4. ✅ 创建 GitHub Release（草稿）
5. ✅ 上传所有资产文件

**发布资产示例：**
```
├── intellijapp_v1.0.0_windows_amd64.exe
├── intellijapp_v1.0.0_windows_arm64.exe
├── intellijapp_v1.0.0_linux_amd64
├── intellijapp_v1.0.0_linux_arm64
├── intellijapp_v1.0.0_darwin_amd64
├── intellijapp_v1.0.0_darwin_arm64
├── intellijapp-amd64-installer.exe      # NSIS 安装程序
├── intellijapp-amd64.msix               # MSIX 包
├── intellijapp-amd64.AppImage           # Linux AppImage
├── intellijapp_1.0.0_amd64.deb         # Debian 包
├── intellijapp-1.0.0-1.x86_64.rpm      # Red Hat 包
├── intellijapp-1.0.0.dmg                # macOS 磁盘镜像
└── intellijapp_v1.0.0_checksums.txt    # 校验和
```

---

## 🛠️ 构建选项

### 环境变量
```bash
# 指定架构
ARCH=amd64 wails3 task build
ARCH=arm64 wails3 task build

# 生产构建
PRODUCTION=true wails3 task build

# 禁用 CGO（静态编译）
CGO_ENABLED=0 go build
```

### 调试与优化
```bash
# 调试构建（保留符号信息）
wails3 build

# 生产构建（优化体积）
wails3 task build PRODUCTION=true

# 查看构建输出
ls -lh bin/
```

---

## 📋 平台要求与依赖

### 🪟 Windows

**最低系统要求：**
- Windows 10 1809+ 或 Windows 11
- 支持架构：AMD64 / ARM64

**开发环境依赖：**
- [WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)（几乎所有 Windows 系统自带）
- [NSIS](https://nsis.sourceforge.io/)（用于创建安装程序）
- [Chocolatey](https://chocolatey.org/)（可选，CI 环境推荐）

**安装 NSIS：**
```bash
# 使用 Chocolatey
choco install nsis -y

# 或手动下载安装
# https://nsis.sourceforge.io/Download
```

**检查依赖：**
```bash
wails3 doctor
```

**运行时依赖（最终用户）：**
- WebView2 Runtime（NSIS 安装程序会自动安装）

---

### 🍎 macOS

**最低系统要求：**
- macOS 10.15+ (Catalina) 用于开发
- macOS 10.13+ (High Sierra) 用于部署
- macOS 11.0+ (Big Sur) 用于 Apple Silicon (ARM64)

**开发环境依赖：**
- Xcode Command Line Tools

**安装命令：**
```bash
xcode-select --install
```

**检查依赖：**
```bash
wails3 doctor
```

**运行时依赖（最终用户）：**
- 无需额外依赖（WebKit 内置于系统）

---

### 🐧 Linux

**推荐发行版：**
- Ubuntu 24.04 AMD64/ARM64（官方支持）
- 其他发行版理论上也可工作

**开发环境依赖：**

#### Debian/Ubuntu
```bash
# 基础构建工具
sudo apt-get update
sudo apt-get install -y build-essential

# GTK3 + WebKit2GTK
sudo apt-get install -y libgtk-3-dev libwebkit2gtk-4.1-dev

# 打包工具（可选）
sudo apt-get install -y nsis rpm
```

#### Red Hat/Fedora/CentOS
```bash
# 基础构建工具
sudo dnf groupinstall "Development Tools"

# GTK3 + WebKit2GTK
sudo dnf install gtk3-devel webkit2gtk4.1-devel

# 打包工具（可选）
sudo dnf install mingw32-nsis rpm-build
```

#### Arch Linux
```bash
# 基础构建工具
sudo pacman -S base-devel

# GTK3 + WebKit2GTK
sudo pacman -S gtk3 webkit2gtk

# 打包工具（可选）
sudo pacman -S nsis rpm-tools
```

**检查依赖：**
```bash
wails3 doctor
```

**运行时依赖（最终用户）：**
- GTK 3.18+
- WebKit2GTK 2.24+

**注意：** AppImage 格式包含所有依赖，可在大多数 Linux 发行版直接运行。

---

### 🛠️ 通用开发依赖

#### Go（必需）
```bash
# 最低版本：Go 1.23+
# 推荐版本：最新稳定版

# 检查安装
go version

# 检查 GOPATH
echo $GOPATH
go env GOPATH
```

**安装：** [Go 官方下载](https://go.dev/dl/)

#### Node.js + npm（强烈推荐）
```bash
# 推荐使用 LTS 版本

# 检查安装
node --version
npm --version
```

**安装：** [Node.js 官方下载](https://nodejs.org/)

#### Wails CLI（必需）
```bash
# 安装最新稳定版
go install github.com/wailsapp/wails/v3/cmd/wails3@latest

# 或安装开发版
git clone https://github.com/wailsapp/wails.git
cd wails
git checkout v3-alpha
cd v3/cmd/wails3
go install

# 检查安装
wails3 version
```

#### Task CLI（推荐）
```bash
# 安装 Task 构建工具
go install github.com/go-task/task/v3/cmd/task@latest

# 检查安装
task --version
```

---

### ✅ 依赖检查清单

运行以下命令检查所有依赖是否正确安装：

```bash
# Wails 系统诊断（最重要）
wails3 doctor

# Go 环境
go version
go env

# Node 环境
node --version
npm --version

# 构建工具
task --version

# 特定平台
# Windows:
where makensis

# macOS:
xcode-select -p

# Linux:
pkg-config --modversion gtk+-3.0
pkg-config --modversion webkit2gtk-4.1
```

---

## 💡 最佳实践

### 1. 版本号管理
使用语义化版本号：
```bash
git tag v1.0.0   # 主要版本
git tag v1.1.0   # 次要版本
git tag v1.1.1   # 补丁版本
```

### 2. 代码签名（推荐）

#### Windows
```bash
# 使用证书签名 NSIS 安装程序
signtool sign /f cert.pfx /p password /tr http://timestamp.digicert.com installer.exe
```

#### macOS
```bash
# 签名应用
codesign --deep --force --verify --verbose --sign "Developer ID" intellijapp.app

# 公证
xcrun notarytool submit intellijapp.zip --apple-id <id> --password <password> --team-id <team>
```

### 3. 测试矩阵
在发布前测试所有格式：
- ✅ Windows: NSIS 安装程序 + MSIX
- ✅ macOS: .app 包（Intel + Apple Silicon）
- ✅ Linux: AppImage + deb + rpm

### 4. 发布检查清单
- [ ] 更新版本号（`wails.json`, `go.mod`）
- [ ] 更新 CHANGELOG
- [ ] 运行完整测试套件
- [ ] 本地构建所有格式
- [ ] 测试主要平台的安装程序
- [ ] 创建并推送标签
- [ ] 验证 CI 构建成功
- [ ] 审核 GitHub Release 草稿
- [ ] 发布 Release

---

## 🐛 故障排查

### Windows 构建失败
```bash
# 确保安装了 NSIS
where makensis

# 重新安装 Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

### macOS 签名问题
```bash
# 检查证书
security find-identity -v -p codesigning

# 清理构建缓存
wails3 clean
```

### Linux 依赖问题
```bash
# 检查依赖
ldd bin/intellijapp

# 安装缺失的库
sudo apt-get install -f
```

---

## 📚 相关资源

- [Wails v3 文档](https://v3alpha.wails.io)
- [GoReleaser 文档](https://goreleaser.com)
- [Task 构建工具](https://taskfile.dev)
- [代码签名指南](/guides/signing/)

---

**祝打包顺利！** 🎉
