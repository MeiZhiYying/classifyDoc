<<<<<<< HEAD
# classifyDoc
智能文件分类
=======
# 智能文件分类系统

一个基于 Node.js 和前端技术的智能文件分类系统，能够根据文件名和内容自动分类文件。

## 功能特性

- 📁 **批量文件上传**: 支持一次性上传100-200个文件
- 🔍 **智能分类**: 两步分类策略
  - 第一步：基于文件名关键词匹配
  - 第二步：AI内容分析（占位符实现）
- 📊 **可视化界面**: 美观的卡片式分类展示
- 📱 **响应式设计**: 支持桌面和移动设备

## 分类类型

- 📄 **合同**: 合同、协议等商务文件
- 👤 **简历**: 个人简历、求职相关文件
- 🧾 **发票**: 发票、收据等财务文件
- 🎓 **论文**: 学术论文、研究报告
- ❓ **未分类**: 未能识别的文件
- ➕ **新增分类**: 可扩展的自定义分类

## 快速开始

### 1. 安装依赖

```bash
npm install
```

### 2. 启动服务

```bash
npm start
```

或者使用开发模式（自动重启）：

```bash
npm run dev
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
- **Express.js**: Web 服务器框架
- **Multer**: 文件上传处理
- **fs-extra**: 文件系统操作

### 前端
- **原生 JavaScript**: 无框架依赖
- **CSS3**: 现代样式和动画
- **Font Awesome**: 图标库

### 分类算法
- **关键词匹配**: 基于预定义关键词库
- **AI 分析**: 占位符接口（可接入真实AI服务）

## 文件结构

```
classifyDocuments/
├── package.json          # 项目配置
├── server.js            # 后端服务器
├── public/              # 前端文件
│   ├── index.html       # 主页面
│   ├── style.css        # 样式文件
│   └── script.js        # 前端逻辑
├── uploads/             # 上传文件存储（自动创建）
└── README.md           # 项目说明
```

## 自定义配置

### 修改分类关键词

编辑 `server.js` 中的 `CLASSIFICATION_KEYWORDS` 对象：

```javascript
const CLASSIFICATION_KEYWORDS = {
  '合同': ['合同', '协议', 'contract', 'agreement'],
  '简历': ['简历', 'resume', 'cv'],
  // 添加更多关键词...
};
```

### 接入真实AI服务

替换 `server.js` 中的 `classifyByAI` 函数：

```javascript
async function classifyByAI(filePath, filename) {
  // 在这里调用真实的AI接口
  // 例如：OpenAI、百度AI、阿里云AI等
  const response = await fetch('YOUR_AI_API_ENDPOINT', {
    method: 'POST',
    headers: { 'Authorization': 'Bearer YOUR_API_KEY' },
    body: JSON.stringify({ file: filePath, name: filename })
  });
  
  const result = await response.json();
  return result.category;
}
```

## 注意事项

- 支持的最大文件数量：200个
- 单个文件大小限制：10MB
- 支持的文件类型：所有类型
- AI分析功能当前为占位符实现

## 后续扩展

- [ ] 接入真实AI大模型服务
- [ ] 添加更多文件分类
- [ ] 支持文件内容预览
- [ ] 添加用户自定义分类
- [ ] 实现文件搜索功能
- [ ] 支持批量导出分类结果
>>>>>>> 31eaaab (feat: Go智能文件分类系统首版)
