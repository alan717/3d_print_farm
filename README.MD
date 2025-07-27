# 3D打印农场系统

一个基于Vue.js + Go + PostgreSQL的云端3D打印服务平台。

## 功能特性

- ✅ 用户注册/登录系统（JWT认证）
- ✅ 3D模型文件上传（STL/OBJ格式，<50MB）
- ✅ Three.js实时3D模型预览
- ✅ 材料选择系统（PLA/ABS/PETG/TPU/WOOD）
- ✅ 打印质量配置（快速/标准/精细）
- ✅ 自动支撑结构生成选项
- ✅ 基于体积+材料的自动报价
- ✅ 订单管理和状态跟踪
- ✅ 响应式设计支持移动端和大屏
- ✅ 完整购物车功能（本地存储、智能合并、批量操作）
- ✅ 演示模型生成器

## 技术栈

### 前端
- Vue 3 + Composition API + Vite
- Element Plus UI组件库 + Icons
- Pinia状态管理
- Vue Router路由管理
- Three.js 3D渲染和模型加载
- Axios HTTP客户端
- 响应式设计（支持移动端到4K屏幕）

### 后端
- Go + Gin框架
- GORM数据库ORM
- JWT身份认证
- PostgreSQL数据库
- 文件上传处理
- CORS中间件

### 数据库
- PostgreSQL 14+
- UUID主键
- JSONB字段类型
- 触发器自动更新时间戳
- 完整的关系约束

## 项目结构

```
/
├── frontend/              # Vue.js前端项目
│   ├── src/
│   │   ├── components/    # 公共组件
│   │   │   ├── AppFooter.vue
│   │   │   ├── NavHeader.vue
│   │   │   └── ThreeViewer.vue  # 3D模型查看器
│   │   ├── views/         # 页面组件
│   │   │   ├── Home.vue
│   │   │   ├── Login.vue         # 登录页（支持大屏展示）
│   │   │   ├── Register.vue
│   │   │   ├── Upload.vue        # 文件上传和配置
│   │   │   ├── Cart.vue          # 购物车页面
│   │   │   ├── Orders.vue
│   │   │   └── Profile.vue
│   │   ├── stores/        # Pinia状态管理
│   │   │   ├── user.js           # 用户状态
│   │   │   ├── model.js          # 模型状态
│   │   │   └── cart.js           # 购物车状态
│   │   ├── router/        # 路由配置
│   │   │   └── index.js
│   │   └── utils/         # 工具函数
│   │       ├── fileUtils.js
│   │       └── modelGenerator.js # 演示模型生成器
│   ├── package.json
│   └── vite.config.js
├── backend/               # Go后端项目
│   ├── cmd/               # 应用入口
│   │   └── main.go
│   ├── internal/          # 内部包
│   │   ├── handlers/      # HTTP处理器
│   │   │   ├── auth.go         # 认证处理
│   │   │   ├── user.go         # 用户管理
│   │   │   ├── file.go         # 文件上传
│   │   │   ├── material.go     # 材料管理
│   │   │   └── order.go        # 订单管理
│   │   ├── models/        # 数据模型
│   │   │   └── models.go       # 所有数据模型
│   │   ├── middleware/    # 中间件
│   │   │   └── middleware.go   # CORS和认证中间件
│   │   └── config/        # 配置
│   │       └── config.go
│   ├── pkg/               # 公共包
│   │   └── database/      # 数据库连接
│   │       └── database.go
│   ├── uploads/           # 文件上传目录
│   ├── .env              # 环境配置
│   └── go.mod
├── database/              # 数据库相关
│   └── schema.sql         # 数据库架构（包含初始数据）
├── docker-compose.yml     # Docker编排配置
├── CLAUDE.md             # Claude Code 项目文档
└── README.md
```

## 系统状态

### ✅ 已部署运行的服务

1. **PostgreSQL数据库** (端口5432)
   - 包含完整表结构：users, materials, files, orders, order_items
   - 预装5种打印材料数据
   - 支持UUID主键和JSONB字段

2. **Go后端服务** (端口8080)
   - Release模式运行，已配置安全代理
   - 完整的REST API接口
   - JWT身份认证
   - 文件上传处理

3. **Vue前端服务** (端口3000)
   - 开发服务器运行
   - 代理配置正确，API调用正常
   - 响应式设计支持各种屏幕尺寸

### 🔧 已解决的技术问题

1. **数据库类型问题** - 修复了JSONB字段类型错误
2. **代理安全警告** - 配置了信任的代理列表
3. **路由重定向** - 修复了材料接口的301重定向
4. **认证错误** - 优化了文件访问的认证策略
5. **JSON格式问题** - 修复了材料properties字段的序列化

### 📋 当前功能状态

| 功能模块 | 状态 | 说明 |
|---------|------|------|
| 用户注册/登录 | ✅ 完成 | JWT认证，密码加密 |
| 3D模型上传 | ✅ 完成 | 支持STL/OBJ，文件大小验证 |
| 3D模型预览 | ✅ 完成 | Three.js渲染，交互控制 |
| 材料选择系统 | ✅ 完成 | 5种材料，价格计算 |
| 打印质量配置 | ✅ 完成 | 三档质量选择 |
| 支撑结构选项 | ✅ 完成 | 可选启用/禁用 |
| 自动报价系统 | ✅ 完成 | 基于体积和材料计算 |
| 购物车功能 | ✅ 完成 | 本地存储，智能合并，批量操作 |
| 订单管理 | ✅ 完成 | 创建、查询、状态更新 |
| 响应式设计 | ✅ 完成 | 支持移动端到4K屏幕 |
| 模型修复提示 | 🚧 待实现 | 计划中 |
| 运费估算 | 🚧 待实现 | 计划中 |
| 支付系统 | 🚧 待实现 | 计划中 |

## 快速开始

### 前置要求
- Node.js 16+
- Go 1.19+
- PostgreSQL 14+

### 本地开发部署

#### 1. 数据库设置

```bash
# 安装并启动PostgreSQL
brew install postgresql@14
brew services start postgresql@14

# 创建数据库
createdb printing_farm

# 执行数据库架构
psql -d printing_farm -f database/schema.sql
```

#### 2. 后端启动

```bash
cd backend

# 创建环境配置
cat > .env << EOF
DATABASE_URL=postgres://$(whoami)@localhost:5432/printing_farm?sslmode=disable
JWT_SECRET=your-secret-key-change-in-production
PORT=8080
UPLOAD_DIR=./uploads
EOF

# 安装依赖
go mod tidy

# 启动服务
go run cmd/main.go
```

#### 3. 前端启动

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

#### 4. 访问应用

- 前端应用：http://localhost:3000
- 后端API：http://localhost:8080
- 数据库：localhost:5432

### 使用Docker Compose

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

## API文档

### 认证接口
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录

### 用户接口
- `GET /api/user/profile` - 获取用户信息（需认证）
- `PUT /api/user/profile` - 更新用户信息（需认证）

### 材料接口
- `GET /api/materials` - 获取材料列表

### 文件接口
- `POST /api/files/upload` - 上传文件（需认证）
- `GET /api/files/:id` - 获取文件
- `DELETE /api/files/:id` - 删除文件（需认证）

### 购物车功能（前端）
- 添加商品到购物车 - 在上传页面配置完成后添加
- 查看购物车 - `/cart` 路由访问
- 管理购物车 - 调整数量、删除商品、批量选择
- 本地存储 - 购物车数据持久化保存

### 订单接口
- `POST /api/orders` - 创建订单（需认证）
- `GET /api/orders` - 获取用户订单列表（需认证）
- `GET /api/orders/:id` - 获取订单详情（需认证）
- `PUT /api/orders/:id/status` - 更新订单状态（需认证）

## 材料信息

| 材料 | 价格(¥/克) | 特性 | 适用场景 |
|------|-----------|------|----------|
| PLA | 0.80 | 环保易用，适合初学者 | 原型制作、装饰品 |
| ABS | 1.00 | 耐用性强，适合功能件 | 机械零件、工具 |
| PETG | 1.20 | 透明度好，化学稳定 | 透明外壳、容器 |
| TPU | 2.00 | 柔性材料，适合弹性件 | 手机壳、垫片 |
| WOOD | 1.50 | 木质纤维复合材料 | 装饰品、艺术品 |

## 开发说明

### 项目特点
- 完全响应式设计，支持移动端到4K屏幕
- 实时3D模型预览和交互
- 智能购物车系统（本地存储、自动合并、批量操作）
- 自动价格计算系统
- 完整的用户认证和权限管理
- 现代化的前后端分离架构

### 技术亮点
- Vue 3 Composition API + Pinia状态管理
- Three.js 3D渲染引擎集成
- 智能购物车本地存储和状态同步
- Go Gin高性能后端框架
- PostgreSQL关系型数据库
- JWT无状态认证
- 文件上传和处理系统

### 购物车系统特性
- **智能合并**: 自动检测相同配置商品并合并数量
- **本地存储**: 购物车数据持久化，刷新页面不丢失
- **实时更新**: 导航栏购物车图标实时显示商品数量
- **批量操作**: 支持选择多个商品批量结算
- **价格计算**: 智能计算（材料价格 + 质量加成 + 支撑费用）
- **响应式界面**: 完美适配各种屏幕尺寸

## 贡献

欢迎提交Issue和Pull Request来改进这个项目。

## 许可证

MIT License

---

*最后更新: 2025-07-27*  
*系统状态: 运行正常*  
*版本: v1.1.0-beta*  
*新增功能: 完整购物车系统*