# IntelliJ Config Helper

一个基于 Wails3 构建的 IntelliJ 系列软件配置路径验证工具。

## 功能特性

- ✅ 验证 IntelliJ 系列软件安装路径
- ✅ 验证配置文件目录
- ✅ 友好的图形界面
- ✅ 详细的错误提示
- ✅ 跨平台支持 (Windows, macOS, Linux)

## 快速开始

### 开发模式

1. 确保已安装 [Wails3](https://v3.wails.io/) 和 [Go 1.24+](https://golang.org/)

2. 克隆项目并进入目录：
   ```bash
   git clone <repository-url>
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

1. 启动应用程序
2. 输入 IntelliJ 软件安装路径（例如：`D:/Program Files/JetBrains/IntelliJ IDEA 2024.2`）
3. 输入配置文件目录路径（例如：`D:/jetbra`）
4. 点击"提交到 Go 后端"按钮进行验证
5. 查看验证结果

## 技术栈

### 后端
- **Go 1.24** - 核心业务逻辑
- **Wails v3** - 桌面应用框架

### 前端
- **Vue 3** - 渐进式JavaScript框架
- **TypeScript** - 类型安全
- **Vite (Rolldown)** - 快速构建工具

## 项目结构

```
intellijapp/
├── frontend/              # 前端代码
│   ├── src/
│   │   ├── App.vue       # 主应用组件
│   │   ├── main.ts       # 入口文件
│   │   └── services/     # 服务层
│   └── package.json
├── internal/              # 内部包
│   └── service/
│       └── greetservice.go  # 路径验证服务
├── build/                 # 构建配置
│   └── config.yml        # 应用配置
├── main.go               # Go 后端入口
└── README.md             # 项目文档
```

## 开发指南

### 添加新功能

1. **后端服务**：在 `internal/service/` 中添加新的服务
2. **前端组件**：在 `frontend/src/` 中创建 Vue 组件
3. **服务绑定**：在 `main.go` 中注册新服务

### 调试

- 后端日志：查看控制台输出
- 前端调试：使用浏览器开发者工具（应用内按 F12）

## 许可证

本项目采用 MIT 许可证。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 相关链接

- [Wails 文档](https://v3.wails.io/)
- [Vue 3 文档](https://vuejs.org/)
- [Go 文档](https://golang.org/doc/)
