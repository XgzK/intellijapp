#!/bin/bash

# IntelliJ Config Helper - macOS 安装脚本
# 自动移除隔离属性，允许应用运行

set -e

echo "=========================================="
echo "  IntelliJ Config Helper - macOS 安装"
echo "=========================================="
echo ""

# 检测是否在 macOS 上运行
if [[ "$OSTYPE" != "darwin"* ]]; then
    echo "❌ 错误：此脚本仅适用于 macOS"
    exit 1
fi

# 查找 .app 文件
APP_NAME="intellijapp.app"
CURRENT_DIR="$(cd "$(dirname "$0")" && pwd)"
APP_PATH=""

# 1. 检查当前目录
if [ -d "$CURRENT_DIR/$APP_NAME" ]; then
    APP_PATH="$CURRENT_DIR/$APP_NAME"
# 2. 检查上层目录
elif [ -d "$CURRENT_DIR/../$APP_NAME" ]; then
    APP_PATH="$CURRENT_DIR/../$APP_NAME"
# 3. 让用户手动指定
else
    echo "⚠️  未找到 $APP_NAME"
    echo ""
    echo "请将此脚本放在与 $APP_NAME 相同的目录中"
    echo "或者手动运行以下命令："
    echo ""
    echo "    xattr -cr /path/to/$APP_NAME"
    echo ""
    exit 1
fi

echo "✓ 找到应用: $APP_PATH"
echo ""

# 检查是否有隔离属性
if xattr "$APP_PATH" | grep -q "com.apple.quarantine"; then
    echo "正在移除隔离属性..."
    xattr -cr "$APP_PATH"
    echo "✓ 隔离属性已移除"
else
    echo "✓ 应用未被隔离，无需处理"
fi

echo ""
echo "=========================================="
echo "  安装完成！"
echo "=========================================="
echo ""
echo "现在可以双击打开 $APP_NAME 了"
echo ""

# 询问是否立即打开应用
read -p "是否立即打开应用？(y/n) " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "正在启动应用..."
    open "$APP_PATH"
fi
