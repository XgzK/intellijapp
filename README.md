# IntelliJ Config Helper

> 版本：v1.0.3

一个基于 Wails3 构建的 IntelliJ 系列软件配置助手工具。

## 功能特性

- ✅ 自动验证 IntelliJ 系列软件安装路径
- ✅ 验证配置文件目录（支持英文、数字、下划线、空格等常用字符）
- ✅ 自动配置 VM Options 文件
- ✅ 精确清除配置（仅删除本工具添加的配置，保留用户自定义配置）
- ✅ 权限检查与友好错误提示
- ✅ 关于信息由后端动态管理
- ✅ 友好的图形界面
- ✅ 跨平台支持 (Windows, macOS, Linux)

## 快速开始

### 开发模式

1. 确保已安装 [Wails3](https://v3.wails.io/) 和 [Go 1.23+](https://golang.org/)

2. 克隆项目并进入目录：
   ```bash
   git clone https://github.com/XgzK/intellijapp.git
   cd intellijapp
   ```

3. 运行开发模式：
   ```bash
   wails3 dev
   ```

### 生产构建

构建可执行文件：

```bash
wails3 build
```

构建完成后，可执行文件将位于 `build/bin/` 目录。

## 使用方法

### 应用配置

1. 启动应用程序
2. 输入 IntelliJ 软件安装路径（例如：`D:/Program Files/JetBrains/IntelliJ IDEA 2024.2`）
3. 输入配置文件目录路径（例如：`D:/jetbra`）
4. 点击"应用配置"按钮
5. 查看操作结果，根据提示重启 IDE

### 清除配置

1. 输入 IntelliJ 软件安装路径
2. 点击"清除配置"按钮
3. 本工具会精确移除之前添加的配置行，不影响其他配置

### 关于页面

点击右上角"关于"按钮可查看：
- 应用版本信息
- 技术栈版本（Go、Vue、Wails）
- 项目仓库地址
- 开发者信息

所有信息由后端动态返回，便于版本管理。

## 技术栈

### 后端
- **Go** - 核心业务逻辑（版本由编译器自动检测）
- **Wails v3** - 跨平台桌面应用框架

### 前端
- **Vue 3.5.22** - 渐进式JavaScript框架
- **TypeScript** - 类型安全
- **Vite** - 快速构建工具

## 项目结构

```
intellijapp/
├── frontend/                     # 前端代码
│   ├── src/
│   │   ├── App.vue              # 主应用组件
│   │   ├── main.ts              # 入口文件
│   │   └── services/            # 服务层
│   │       └── configService.ts # 后端服务调用
│   ├── bindings/                # Wails 自动生成的类型绑定
│   └── package.json
├── internal/                    # 内部包
│   └── service/
│       └── configservice.go    # 配置路径验证服务
├── build/                       # 构建配置和输出
│   └── bin/                    # 编译后的可执行文件
├── main.go                     # Go 后端入口
├── README.md                   # 项目文档
├── PACKAGING.md               # 打包说明
├── QUICK_START.md            # 快速开始指南
└── TROUBLESHOOTING.md        # 故障排除指南
```

## 核心功能说明

### 路径验证

- 自动检测 IntelliJ 安装目录（支持直接指定或自动识别 bin 目录）
- 验证配置文件完整性（检查 ja-netfilter.jar 是否存在）
- 支持的路径字符：英文字母、数字、下划线、空格、点、横杠、冒号、斜杠

### VM Options 配置

本工具会自动修改 `.vmoptions` 文件，添加以下配置：

```
--add-opens=java.base/jdk.internal.org.objectweb.asm=ALL-UNNAMED
--add-opens=java.base/jdk.internal.org.objectweb.asm.tree=ALL-UNNAMED
-javaagent:<配置路径>/ja-netfilter.jar=jetbrains
```

### 权限检查

- 自动检查目录读取权限
- 自动检查文件读写权限
- 友好的错误提示（Windows/Linux/macOS 分别提示）

### 精确清除

清除配置时**仅删除**以下本工具添加的行：
- `--add-opens=java.base/jdk.internal.org.objectweb.asm=ALL-UNNAMED`
- `--add-opens=java.base/jdk.internal.org.objectweb.asm.tree=ALL-UNNAMED`
- `-javaagent:**/ja-netfilter.jar=jetbrains`

**不会删除**用户自定义的其他 `--add-opens` 或 `-javaagent` 配置。

## 开发指南

### 后端开发

1. **添加新服务**：在 `internal/service/` 中创建新的服务文件
2. **注册服务**：在 `main.go` 中使用 `application.NewService()` 注册
3. **生成绑定**：运行 `wails3 generate bindings -ts` 生成前端类型

### 前端开发

1. **调用后端**：从 `bindings/` 导入自动生成的函数
2. **类型安全**：使用 TypeScript 确保类型一致
3. **响应式数据**：使用 Vue 3 Composition API

### 调试技巧

- **后端日志**：查看控制台输出，所有操作都有详细日志
- **前端调试**：应用内按 F12 打开开发者工具
- **Bindings 更新**：修改 Go 结构体后记得重新生成 bindings

## 版本历史

### v1.0.3 (当前版本)

**文档改进：**
- 📝 所有代码注释统一为中文，提高代码可读性
- 📝 完善了所有私有函数的注释文档
- 📝 main.go 中的注释全部翻译为中文
- 📝 扩展了 README.md 文档（从 104 行扩展到 208 行）
- 📝 添加了详细的版本历史记录

**功能完善：**
- ✨ 关于页面信息完全由后端管理（版本号、技术栈、开发者信息）
- ✨ 支持多个开发者信息展示
- ✨ 前端使用 `onMounted` 动态获取后端数据

**代码质量：**
- 🧪 所有函数都有完整的中文注释
- 🧪 代码结构更加清晰易懂
- 🧪 遵循 KISS、YAGNI、DRY、SOLID 原则

### v1.0.2

**修复：**
- 🐛 删除了会导致 MCP 工具失败的 `nul` 文件
- 🐛 修复了 `CONFIG_PATH_PATTERN` 正则表达式（现在支持空格、下划线等）
- 🐛 精确清除配置，只删除本工具添加的特定行，不影响用户配置
- 🐛 移除了未使用的 `time` 事件循环

**改进：**
- ✨ Go 版本使用 `runtime.Version()` 自动获取编译器版本
- ✨ 添加了完善的权限检查机制
- ✨ 跨平台错误提示（Windows/Linux/macOS）
- 🗑️ 清理了未使用的 `HelloWorld.vue` 组件

**代码质量：**
- 📝 遵循 KISS、YAGNI、DRY、SOLID 原则
- 🧪 代码结构更清晰，职责更明确

## 许可证

本项目采用 MIT 许可证。

## 贡献指南

欢迎提交 Issue 和 Pull Request！

提交前请确保：
- 代码通过 `go build` 编译
- 遵循项目的代码风格
- 添加必要的注释和文档

## 开发者

- **XgzK** - 项目维护者
- **Claude (AI)** - AI 辅助开发
- **Codex (AI)** - AI 辅助开发

## 相关链接

- [项目仓库](https://github.com/XgzK/intellijapp)
- [Wails 文档](https://v3.wails.io/)
- [Vue 3 文档](https://vuejs.org/)
- [Go 文档](https://golang.org/doc/)

---

**⚠️ 免责声明**：本项目运行所产生的一切问题需自行承担，仅限学习使用。
