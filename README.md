# ApplyPages - 高校招生申请系统

ApplyPages 是一个现代化的高校招生申请管理系统，为高校和申请者提供一个便捷的在线申请平台。

## 功能特性

- 🔐 支持多种登录方式（OAuth认证）
- 📝 在线申请表单提交
- 👤 用户信息管理
- 📊 申请进度追踪
- 🔄 实时状态更新

## 技术栈

### 后端
- Go
- Gin Web Framework
- GORM
- MySQL

### 前端
- Vue.js 3
- Vue Router
- Modern CSS

## 系统要求

- Go 1.16+
- MySQL 5.7+
- Node.js 16+
- npm 或 yarn

## 快速开始

### 1. 克隆项目


### 2. 配置环境

1. 复制配置文件模板：
```bash
cp config/example.config.yaml config/config.yaml
```

2. 修改 `config/config.yaml` 中的配置：
   - 数据库连接信息
   - 服务器端口
   - OAuth 配置

3. 修改 `web/src/assets/BigLOGO.png` 作为学校logo
4. 修改 `web/index.html` 中的校名
5. 修改 `web/src/views/RegisterView.vue` 中的邮箱后缀
6. 修改 `web/src/views/LoginView.vue` 中的subtitle

### 3. 安装依赖

后端依赖：
```bash
go mod download
```

前端依赖：
```bash
cd web
npm install
```

### 4. 构建前端

```bash
cd web
npm run build
```

### 5. 启动服务

```bash
go run main.go
```

服务将在配置的端口上启动（默认为 :8080）

## 项目结构

```
.
├── config/          # 配置文件
├── internal/        # 内部包
│   ├── database/   # 数据库相关
│   ├── handler/    # 请求处理器
│   └── model/      # 数据模型
├── web/            # 前端代码
│   ├── src/        # Vue.js 源代码
│   └── public/     # 静态资源
└── main.go         # 程序入口
```

## 开发

### 后端开发
```bash
go run main.go
```

### 前端开发
```bash
cd web
npm run dev
```
