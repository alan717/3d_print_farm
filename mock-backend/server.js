import express from 'express'
import cors from 'cors'
import multer from 'multer'
import jwt from 'jsonwebtoken'
import bcrypt from 'bcryptjs'
import { v4 as uuidv4 } from 'uuid'
import fs from 'fs'
import path from 'path'

const app = express()
const PORT = 8080

// 模拟数据库
const users = []
const files = []
const orders = []

// 材料数据
const materials = [
  {
    id: 1,
    name: 'PLA',
    description: '环保易用，适合初学者',
    price_per_gram: 0.8,
    properties: { temperature: '190-220°C', flexibility: '低', strength: '中等' },
    available: true
  },
  {
    id: 2,
    name: 'ABS',
    description: '耐用性强，适合功能件',
    price_per_gram: 1.0,
    properties: { temperature: '220-260°C', flexibility: '中', strength: '高' },
    available: true
  },
  {
    id: 3,
    name: 'PETG',
    description: '透明度好，化学稳定',
    price_per_gram: 1.2,
    properties: { temperature: '220-250°C', flexibility: '中', strength: '高' },
    available: true
  },
  {
    id: 4,
    name: 'TPU',
    description: '柔性材料，适合弹性件',
    price_per_gram: 2.0,
    properties: { temperature: '210-230°C', flexibility: '高', strength: '中等' },
    available: true
  },
  {
    id: 5,
    name: 'WOOD',
    description: '木质纤维复合材料',
    price_per_gram: 1.5,
    properties: { temperature: '190-220°C', flexibility: '低', strength: '中等', 特性: '木质纹理' },
    available: true
  }
]

// 中间件
app.use(cors())
app.use(express.json())
app.use('/uploads', express.static('uploads'))

// 确保uploads目录存在
if (!fs.existsSync('uploads')) {
  fs.mkdirSync('uploads')
}

// 文件上传配置
const storage = multer.diskStorage({
  destination: (req, file, cb) => {
    cb(null, 'uploads/')
  },
  filename: (req, file, cb) => {
    const ext = path.extname(file.originalname)
    const filename = uuidv4() + ext
    cb(null, filename)
  }
})

const upload = multer({ 
  storage,
  limits: {
    fileSize: 50 * 1024 * 1024 // 50MB
  },
  fileFilter: (req, file, cb) => {
    const allowedTypes = ['.stl', '.obj']
    const ext = path.extname(file.originalname).toLowerCase()
    if (allowedTypes.includes(ext)) {
      cb(null, true)
    } else {
      cb(new Error('只支持STL和OBJ格式'))
    }
  }
})

// JWT密钥
const JWT_SECRET = 'your-secret-key'

// 认证中间件
const authenticateToken = (req, res, next) => {
  const authHeader = req.headers['authorization']
  const token = authHeader && authHeader.split(' ')[1]

  if (!token) {
    return res.status(401).json({ error: 'Access token required' })
  }

  jwt.verify(token, JWT_SECRET, (err, user) => {
    if (err) {
      return res.status(403).json({ error: 'Invalid token' })
    }
    req.user = user
    next()
  })
}

// 路由

// 注册
app.post('/api/auth/register', async (req, res) => {
  try {
    const { username, email, password, phone, address } = req.body

    // 检查用户是否已存在
    const existingUser = users.find(u => u.username === username || u.email === email)
    if (existingUser) {
      return res.status(409).json({ error: '用户名或邮箱已存在' })
    }

    // 加密密码
    const hashedPassword = await bcrypt.hash(password, 10)

    // 创建用户
    const user = {
      id: uuidv4(),
      username,
      email,
      password: hashedPassword,
      phone: phone || '',
      address: address || '',
      created_at: new Date().toISOString()
    }

    users.push(user)

    // 生成token
    const token = jwt.sign({ user_id: user.id }, JWT_SECRET, { expiresIn: '24h' })

    res.status(201).json({
      token,
      user: {
        id: user.id,
        username: user.username,
        email: user.email,
        phone: user.phone,
        address: user.address
      }
    })
  } catch (error) {
    res.status(500).json({ error: '注册失败' })
  }
})

// 登录
app.post('/api/auth/login', async (req, res) => {
  try {
    const { email, password } = req.body

    // 查找用户
    const user = users.find(u => u.email === email)
    if (!user) {
      return res.status(401).json({ error: '邮箱或密码错误' })
    }

    // 验证密码
    const validPassword = await bcrypt.compare(password, user.password)
    if (!validPassword) {
      return res.status(401).json({ error: '邮箱或密码错误' })
    }

    // 生成token
    const token = jwt.sign({ user_id: user.id }, JWT_SECRET, { expiresIn: '24h' })

    res.json({
      token,
      user: {
        id: user.id,
        username: user.username,
        email: user.email,
        phone: user.phone,
        address: user.address
      }
    })
  } catch (error) {
    res.status(500).json({ error: '登录失败' })
  }
})

// 获取用户信息
app.get('/api/user/profile', authenticateToken, (req, res) => {
  const user = users.find(u => u.id === req.user.user_id)
  if (!user) {
    return res.status(404).json({ error: '用户不存在' })
  }

  res.json({
    id: user.id,
    username: user.username,
    email: user.email,
    phone: user.phone,
    address: user.address
  })
})

// 获取材料列表
app.get('/api/materials', (req, res) => {
  res.json(materials.filter(m => m.available))
})

// 文件上传
app.post('/api/files/upload', authenticateToken, upload.single('file'), (req, res) => {
  try {
    if (!req.file) {
      return res.status(400).json({ error: '没有文件上传' })
    }

    const fileRecord = {
      id: uuidv4(),
      user_id: req.user.user_id,
      filename: req.file.filename,
      original_name: req.file.originalname,
      size: req.file.size,
      mime_type: req.file.mimetype,
      path: req.file.path,
      url: `/uploads/${req.file.filename}`,
      created_at: new Date().toISOString()
    }

    files.push(fileRecord)

    res.json(fileRecord)
  } catch (error) {
    res.status(500).json({ error: '文件上传失败' })
  }
})

// 获取文件
app.get('/api/files/:id', (req, res) => {
  const file = files.find(f => f.id === req.params.id)
  if (!file) {
    return res.status(404).json({ error: '文件不存在' })
  }

  res.sendFile(path.resolve(file.path))
})

// 创建订单
app.post('/api/orders', authenticateToken, (req, res) => {
  const order = {
    id: uuidv4(),
    user_id: req.user.user_id,
    ...req.body,
    status: 'pending',
    created_at: new Date().toISOString()
  }

  orders.push(order)
  res.status(201).json(order)
})

// 获取用户订单
app.get('/api/orders', authenticateToken, (req, res) => {
  const userOrders = orders.filter(o => o.user_id === req.user.user_id)
  res.json(userOrders)
})

// 错误处理中间件
app.use((error, req, res, next) => {
  if (error instanceof multer.MulterError) {
    if (error.code === 'LIMIT_FILE_SIZE') {
      return res.status(400).json({ error: '文件大小超过限制' })
    }
  }
  res.status(500).json({ error: error.message })
})

app.listen(PORT, () => {
  console.log(`Mock backend server running on http://localhost:${PORT}`)
  console.log('Available endpoints:')
  console.log('- POST /api/auth/register')
  console.log('- POST /api/auth/login')
  console.log('- GET /api/user/profile')
  console.log('- GET /api/materials')
  console.log('- POST /api/files/upload')
  console.log('- GET /api/files/:id')
  console.log('- POST /api/orders')
  console.log('- GET /api/orders')
})