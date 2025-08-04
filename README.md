# 智能文件分类系统

一个基于 Go 语言开发的智能文件分类系统，能够根据文件名和内容自动分类文件，提供美观的Web界面。

## 功能特性

- 📁 **批量文件上传**: 支持一次性上传100-200个文件
- 🔍 **智能分类**: 两步分类策略
  - 第一步：基于文件名关键词匹配
  - 第二步：AI内容分析（占位符实现）
- 📊 **可视化界面**: 美观的卡片式分类展示
- 📱 **响应式设计**: 支持桌面和移动设备
- 🏗️ **企业级架构**: 模块化设计，易于维护和扩展

## 分类类型

- 📄 **合同**: 合同、协议等商务文件
- 👤 **简历**: 个人简历、求职相关文件
- 🧾 **发票**: 发票、收据等财务文件
- 🎓 **论文**: 学术论文、研究报告
- ❓ **未分类**: 未能识别的文件
- ➕ **新增分类**: 可扩展的自定义分类

## 项目结构

```
file-classifier/
├── cmd/
│   └── server/          # 应用程序入口
│       └── main.go
├── internal/            # 内部包，不对外暴露
│   ├── config/          # 配置文件
│   ├── handlers/        # HTTP处理器
│   ├── models/          # 数据模型
│   ├── router/          # 路由配置
│   ├── service/         # 业务逻辑
│   └── utils/           # 工具函数
├── public/              # 前端静态文件
│   ├── index.html
│   ├── style.css
│   └── script.js
├── Makefile            # 构建脚本
├── go.mod              # Go模块文件
└── README.md           # 项目文档
```

## 快速开始

### 方法一：使用 Makefile (推荐)

```bash
# 查看所有可用命令
make help

# 运行应用程序
make run

# 构建应用程序
make build

# 开发模式（自动重启）
make dev
```

### 方法二：直接使用 Go 命令

```bash
# 安装依赖
go mod tidy

# 运行应用程序
go run cmd/server/main.go

# 构建应用程序
go build -o bin/file-classifier cmd/server/main.go
```

### 3. 访问应用

打开浏览器访问: `http://localhost:3000`

## 使用说明

1. **上传文件**: 点击上传区域或拖拽文件到页面
2. **选择文件夹**: 选择包含多个文件的文件夹进行批量上传
3. **等待分类**: 系统将自动进行两步分类处理
4. **查看结果**: 点击分类卡片查看具体的文件列表

## 技术架构

### 后端
- **Go 1.21+**: 主要编程语言
- **Gin**: Web框架
- **企业级架构**: 模块化设计，分离关注点

### 前端
- **原生 JavaScript**: 无框架依赖
- **CSS3**: 现代样式和动画
- **Font Awesome**: 图标库

### 分类算法
- **关键词匹配**: 基于预定义关键词库
- **AI 分析**: 占位符接口（可接入真实AI服务）

## 开发指南

### 环境要求
- Go 1.21 或更高版本
- 现代浏览器支持

### 本地开发

```bash
# 克隆仓库
git clone https://github.com/MeiZhiYing/classifyDoc.git
cd classifyDoc

# 安装依赖
go mod tidy

# 运行应用
make run
```

### 代码质量

```bash
# 格式化代码
make fmt

# 静态检查
make vet

# 运行测试
make test
```

## 自定义配置

### 修改分类关键词

编辑 `internal/config/config.go` 中的 `ClassificationKeywords`：

```go
var ClassificationKeywords = map[string][]string{
    "合同": {"合同", "协议", "contract", "agreement"},
    "简历": {"简历", "resume", "cv"},
    // 添加更多关键词...
}
```

### 接入真实AI服务

替换 `internal/service/classifier.go` 中的 `ClassifyByAI` 函数：

```go
func ClassifyByAI(filePath, filename string) string {
    // 在这里调用真实的AI接口
    // 例如：OpenAI、百度AI、阿里云AI等
    
    // 示例：调用外部API
    response, err := http.Post("YOUR_AI_API_ENDPOINT", "application/json", body)
    if err != nil {
        log.Printf("AI分析失败: %v", err)
        return "未分类"
    }
    
    // 解析响应并返回分类结果
    return parseAIResponse(response)
}
```

## 部署

### 构建生产版本

```bash
# 构建二进制文件
make build

# 运行
./bin/file-classifier
```

### Docker 部署

```bash
# 构建镜像
make docker-build

# 运行容器
make docker-run
```

## 注意事项

- 支持的最大文件数量：200个
- 单个文件大小限制：100MB
- 支持的文件类型：所有类型
- AI分析功能当前为占位符实现

## 后续扩展

- [ ] 接入真实AI大模型服务
- [ ] 添加更多文件分类
- [ ] 支持文件内容预览
- [ ] 添加用户自定义分类
- [ ] 实现文件搜索功能
- [ ] 支持批量导出分类结果
- [ ] 添加用户权限管理
- [ ] 实现分类规则学习

## 贡献

欢迎提交 Issue 和 Pull Request 来改进项目。

## 许可证

MIT License