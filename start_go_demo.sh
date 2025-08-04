#!/bin/bash

echo "======================================"
echo "🚀 Go智能文件分类系统启动中..."
echo "======================================"

# 检查Go是否已安装
if ! command -v go &> /dev/null; then
    echo "❌ 错误：Go 未安装"
    echo "请先安装 Go: https://golang.org/dl/"
    exit 1
fi

# 检查依赖是否已安装
if [ ! -f "go.sum" ]; then
    echo "📦 正在安装Go依赖..."
    go mod tidy
fi

# 停止可能运行的Node.js服务
pkill -f "node server.js" 2>/dev/null || true

# 启动Go服务器
echo "🌐 启动Go服务器..."
echo "访问地址: http://localhost:3000"
echo "按 Ctrl+C 退出服务器"
echo "======================================"

# 等待用户确认
echo "按任意键继续..."
read -n 1 -s

# 启动服务
go run main.go