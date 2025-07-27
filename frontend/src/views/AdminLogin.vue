<template>
  <div class="admin-login">
    <div class="login-container">
      <div class="login-header">
        <div class="logo">
          <el-icon size="40"><Lock /></el-icon>
        </div>
        <h1>管理后台</h1>
        <p>3D打印农场管理系统</p>
      </div>

      <el-card class="login-card" shadow="always">
        <el-form
          ref="loginFormRef"
          :model="loginForm"
          :rules="loginRules"
          label-position="top"
          size="large"
          @keyup.enter="handleLogin"
        >
          <el-form-item label="管理员账号" prop="username">
            <el-input
              v-model="loginForm.username"
              placeholder="请输入管理员用户名"
              prefix-icon="User"
              clearable
            />
          </el-form-item>

          <el-form-item label="密码" prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="请输入密码"
              prefix-icon="Lock"
              show-password
              clearable
            />
          </el-form-item>

          <el-form-item>
            <el-checkbox v-model="loginForm.remember">记住登录状态</el-checkbox>
          </el-form-item>

          <el-form-item>
            <el-button
              type="primary"
              style="width: 100%;"
              size="large"
              :loading="loading"
              @click="handleLogin"
            >
              <span v-if="!loading">登录管理后台</span>
              <span v-else>登录中...</span>
            </el-button>
          </el-form-item>
        </el-form>

        <div class="login-footer">
          <el-divider>
            <span class="divider-text">其他操作</span>
          </el-divider>
          
          <div class="footer-links">
            <el-button text @click="$router.push('/')">
              <el-icon><House /></el-icon>
              返回网站首页
            </el-button>
            <el-button text @click="$router.push('/login')">
              <el-icon><User /></el-icon>
              用户登录
            </el-button>
          </div>
        </div>
      </el-card>

      <!-- 系统信息 -->
      <div class="system-info">
        <p>© 2025 3D打印农场管理系统</p>
        <p>请使用管理员账户登录</p>
      </div>
    </div>

    <!-- 背景装饰 -->
    <div class="background-decoration">
      <div class="decoration-item decoration-1"></div>
      <div class="decoration-item decoration-2"></div>
      <div class="decoration-item decoration-3"></div>
      <div class="decoration-item decoration-4"></div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { Lock, User, House } from '@element-plus/icons-vue'
import axios from 'axios'

export default {
  name: 'AdminLogin',
  components: {
    Lock, User, House
  },
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    const loginFormRef = ref()
    const loading = ref(false)

    const loginForm = reactive({
      username: 'admin',
      password: '',
      remember: true
    })

    const loginRules = {
      username: [
        { required: true, message: '请输入管理员用户名', trigger: 'blur' },
        { min: 2, max: 50, message: '用户名长度在 2 到 50 个字符', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, message: '密码长度不能少于 6 个字符', trigger: 'blur' }
      ]
    }

    const handleLogin = async () => {
      try {
        const valid = await loginFormRef.value.validate()
        if (!valid) return

        loading.value = true

        const response = await axios.post('/api/auth/login', {
          username: loginForm.username,
          password: loginForm.password
        })

        const { token, user } = response.data

        // 检查用户是否有管理员权限
        if (!user.role || (user.role !== 'admin' && user.role !== 'super_admin')) {
          throw new Error('当前账户没有管理员权限')
        }

        // 保存登录信息
        userStore.setToken(token)
        userStore.setUserInfo(user)

        if (loginForm.remember) {
          localStorage.setItem('admin_remember', 'true')
        }

        ElMessage.success(`欢迎回来，${user.username}`)
        
        // 跳转到管理后台
        router.push('/admin/dashboard')

      } catch (error) {
        console.error('管理员登录失败:', error)
        const message = error.response?.data?.error || error.message || '登录失败'
        ElMessage.error(message)
      } finally {
        loading.value = false
      }
    }

    // 检查是否已经登录且有管理员权限
    onMounted(() => {
      const userInfo = userStore.userInfo
      if (userStore.isLoggedIn && userInfo?.role && 
          (userInfo.role === 'admin' || userInfo.role === 'super_admin')) {
        router.push('/admin/dashboard')
      }

      // 恢复记住的登录状态
      const remembered = localStorage.getItem('admin_remember')
      if (remembered === 'true') {
        loginForm.remember = true
      }
    })

    return {
      loginFormRef,
      loginForm,
      loginRules,
      loading,
      handleLogin
    }
  }
}
</script>

<style scoped>
.admin-login {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.login-container {
  width: 100%;
  max-width: 420px;
  padding: 20px;
  position: relative;
  z-index: 10;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
  color: white;
}

.logo {
  width: 80px;
  height: 80px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 20px;
  backdrop-filter: blur(10px);
  border: 2px solid rgba(255, 255, 255, 0.2);
}

.login-header h1 {
  margin: 0 0 8px 0;
  font-size: 2rem;
  font-weight: 300;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.login-header p {
  margin: 0;
  opacity: 0.9;
  font-size: 1rem;
}

.login-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  overflow: hidden;
}

.login-card :deep(.el-card__body) {
  padding: 32px;
}

.login-footer {
  margin-top: 24px;
}

.divider-text {
  color: #999;
  font-size: 0.9rem;
}

.footer-links {
  display: flex;
  justify-content: space-between;
  margin-top: 16px;
}

.footer-links .el-button {
  color: #666;
  font-size: 0.9rem;
}

.footer-links .el-button:hover {
  color: #409eff;
}

.system-info {
  text-align: center;
  margin-top: 32px;
  color: rgba(255, 255, 255, 0.8);
}

.system-info p {
  margin: 4px 0;
  font-size: 0.9rem;
}

/* 背景装饰 */
.background-decoration {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  overflow: hidden;
}

.decoration-item {
  position: absolute;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 50%;
  animation: float 6s ease-in-out infinite;
}

.decoration-1 {
  width: 100px;
  height: 100px;
  top: 10%;
  left: 10%;
  animation-delay: 0s;
}

.decoration-2 {
  width: 150px;
  height: 150px;
  top: 60%;
  right: 10%;
  animation-delay: 2s;
}

.decoration-3 {
  width: 80px;
  height: 80px;
  bottom: 20%;
  left: 15%;
  animation-delay: 4s;
}

.decoration-4 {
  width: 120px;
  height: 120px;
  top: 20%;
  right: 20%;
  animation-delay: 1s;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
    opacity: 0.3;
  }
  50% {
    transform: translateY(-20px) rotate(180deg);
    opacity: 0.6;
  }
}

/* 表单样式优化 */
.el-form-item {
  margin-bottom: 24px;
}

.el-form-item:last-child {
  margin-bottom: 0;
}

:deep(.el-input__inner) {
  border-radius: 8px;
  border: 1px solid #e0e0e0;
  transition: all 0.3s;
}

:deep(.el-input__inner:focus) {
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.1);
}

:deep(.el-button--primary) {
  border-radius: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  font-weight: 500;
  transition: all 0.3s;
}

:deep(.el-button--primary:hover) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(64, 158, 255, 0.3);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .login-container {
    max-width: 90%;
    padding: 15px;
  }

  .login-header h1 {
    font-size: 1.6rem;
  }

  .login-card :deep(.el-card__body) {
    padding: 24px;
  }

  .footer-links {
    flex-direction: column;
    gap: 8px;
  }

  .footer-links .el-button {
    width: 100%;
    justify-content: center;
  }

  .decoration-item {
    display: none;
  }
}

@media (max-width: 480px) {
  .login-container {
    max-width: 95%;
    padding: 10px;
  }

  .logo {
    width: 60px;
    height: 60px;
  }

  .login-header h1 {
    font-size: 1.4rem;
  }

  .login-card :deep(.el-card__body) {
    padding: 20px;
  }
}
</style>