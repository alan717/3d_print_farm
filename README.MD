# 3D打印农场系统

一个基于Vue.js + Go + PostgreSQL的云端3D打印服务平台。

## 功能特性

- 用户注册/登录系统
- 3D模型文件上传（STL/OBJ格式）
- 材料选择和打印质量配置
- 自动价格计算
- 订单管理和状态跟踪
- 响应式设计支持移动端

## 技术栈

### 前端
- Vue 3 + Vite
- Element Plus UI组件库
- Pinia状态管理
- Vue Router路由管理
- Three.js 3D渲染（计划中）

### 后端
- Go + Gin框架
- GORM数据库ORM
- JWT身份认证
- PostgreSQL数据库

## 项目结构

```
/
├── frontend/              # Vue.js前端项目
│   ├── src/
│   │   ├── components/    # 公共组件
│   │   ├── views/         # 页面组件
│   │   ├── stores/        # Pinia状态管理
│   │   ├── router/        # 路由配置
│   │   └── utils/         # 工具函数
│   ├── package.json
│   └── vite.config.js
├── backend/               # Go后端项目
│   ├── cmd/               # 应用入口
│   ├── internal/          # 内部包
│   │   ├── handlers/      # HTTP处理器
│   │   ├── models/        # 数据模型
│   │   ├── services/      # 业务逻辑
│   │   ├── middleware/    # 中间件
│   │   └── config/        # 配置
│   ├── pkg/               # 公共包
│   └── go.mod
├── database/              # 数据库相关
│   └── schema.sql         # 数据库架构
├── docker-compose.yml     # Docker编排配置
└── README.md
```

## 快速开始

### 使用Docker Compose（推荐）

1. 克隆项目
```bash
git clone <repository>
cd sys_3d
```

2. 启动服务
```bash
docker-compose up -d
```

3. 访问应用
- 前端：http://localhost:3000
- 后端API：http://localhost:8080

### 手动部署

#### 数据库设置

1. 安装PostgreSQL
2. 创建数据库：
```sql
CREATE DATABASE printing_farm;
```

3. 执行数据库架构：
```bash
psql -d printing_farm -f database/schema.sql
```

#### 后端启动

1. 进入后端目录：
```bash
cd backend
```

2. 复制环境配置：
```bash
cp .env.example .env
```

3. 修改.env文件中的数据库连接信息

4. 安装依赖并启动：
```bash
go mod tidy
go run cmd/main.go
```

#### 前端启动

1. 进入前端目录：
```bash
cd frontend
```

2. 安装依赖：
```bash
npm install
```

3. 启动开发服务器：
```bash
npm run dev
```

## API文档

### 认证接口
- POST `/api/auth/register` - 用户注册
- POST `/api/auth/login` - 用户登录

### 用户接口
- GET `/api/user/profile` - 获取用户信息
- PUT `/api/user/profile` - 更新用户信息

### 材料接口
- GET `/api/materials` - 获取材料列表

### 文件接口
- POST `/api/files/upload` - 上传文件
- GET `/api/files/:id` - 获取文件
- DELETE `/api/files/:id` - 删除文件

### 订单接口
- POST `/api/orders` - 创建订单
- GET `/api/orders` - 获取用户订单列表
- GET `/api/orders/:id` - 获取订单详情
- PUT `/api/orders/:id/status` - 更新订单状态

## 开发状态

✅ 已完成：
- 项目架构设计
- 用户注册/登录系统
- 基础页面和组件
- 数据库设计
- API接口框架
- 材料管理系统
- 订单状态跟踪

🚧 开发中：
- 文件上传功能
- 3D模型预览
- 价格计算系统

📋 计划中：
- 支付系统集成
- 3D模型修复提示
- 运费计算
- 性能优化

## 贡献

欢迎提交Issue和Pull Request来改进这个项目。

## 许可证

MIT License