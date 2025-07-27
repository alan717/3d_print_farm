<template>
  <div class="login-page">
    <div class="login-container">
      <!-- 大屏幕左侧展示区域 -->
      <div class="login-showcase">
        <div class="showcase-content">
          <h1 class="showcase-title">3D打印农场</h1>
          <p class="showcase-subtitle">专业的云端3D打印服务平台</p>
          <div class="showcase-features">
            <div class="feature-item">
              <div class="feature-icon">🎯</div>
              <span>精准建模预览</span>
            </div>
            <div class="feature-item">
              <div class="feature-icon">⚡</div>
              <span>快速打印服务</span>
            </div>
            <div class="feature-item">
              <div class="feature-icon">🔧</div>
              <span>多种材料选择</span>
            </div>
            <div class="feature-item">
              <div class="feature-icon">💎</div>
              <span>高品质成品</span>
            </div>
          </div>
          <div class="showcase-stats">
            <div class="stat-item">
              <div class="stat-number">10,000+</div>
              <div class="stat-label">完成订单</div>
            </div>
            <div class="stat-item">
              <div class="stat-number">5,000+</div>
              <div class="stat-label">活跃用户</div>
            </div>
            <div class="stat-item">
              <div class="stat-number">99.9%</div>
              <div class="stat-label">成功率</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧登录表单区域 -->
      <div class="login-form-section">
        <el-row justify="center">
          <el-col :span="24" :xs="22" :sm="18" :md="14" :lg="20" :xl="18">
            <el-card class="login-card">
              <template #header>
                <div class="card-header">
                  <h2>用户登录</h2>
                  <p class="header-subtitle">登录您的账户，开始3D打印之旅</p>
                </div>
              </template>
              
              <el-form
                :model="loginForm"
                :rules="rules"
                ref="loginFormRef"
                :label-width="labelWidth"
                label-position="top"
                @submit.prevent="handleLogin"
                class="login-form"
              >
                <el-form-item label="邮箱" prop="email">
                  <el-input
                    v-model="loginForm.email"
                    type="email"
                    placeholder="请输入邮箱"
                    prefix-icon="Message"
                    size="large"
                    class="form-input"
                  />
                </el-form-item>
                
                <el-form-item label="密码" prop="password">
                  <el-input
                    v-model="loginForm.password"
                    type="password"
                    placeholder="请输入密码"
                    prefix-icon="Lock"
                    show-password
                    size="large"
                    class="form-input"
                  />
                </el-form-item>
                
                <el-form-item>
                  <el-button
                    type="primary"
                    size="large"
                    :loading="loading"
                    @click="handleLogin"
                    class="login-button"
                  >
                    登录
                  </el-button>
                </el-form-item>
                
                <el-form-item>
                  <div class="login-footer">
                    <span>还没有账号？</span>
                    <router-link to="/register" class="register-link">立即注册</router-link>
                  </div>
                </el-form-item>
              </el-form>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

export default {
  name: 'Login',
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    const loginFormRef = ref()
    const loading = ref(false)
    
    // 响应式标签宽度
    const labelWidth = computed(() => {
      // 在移动端使用较小的标签宽度
      return window.innerWidth < 768 ? '60px' : '80px'
    })
    
    const loginForm = reactive({
      email: '',
      password: ''
    })
    
    const rules = {
      email: [
        { required: true, message: '请输入邮箱', trigger: 'blur' },
        { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, message: '密码长度至少6位', trigger: 'blur' }
      ]
    }
    
    const handleLogin = async () => {
      if (!loginFormRef.value) return
      
      try {
        await loginFormRef.value.validate()
        loading.value = true
        
        const result = await userStore.login(loginForm)
        
        if (result.success) {
          ElMessage.success('登录成功')
          router.push('/')
        } else {
          ElMessage.error(result.message)
        }
      } catch (error) {
        console.error('登录失败:', error)
      } finally {
        loading.value = false
      }
    }
    
    return {
      loginForm,
      rules,
      loginFormRef,
      loading,
      labelWidth,
      handleLogin
    }
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px 16px;
}

.login-container {
  width: 100%;
  max-width: 1400px;
  display: flex;
  align-items: center;
  gap: 60px;
}

/* 左侧展示区域 - 仅在大屏显示 */
.login-showcase {
  flex: 1;
  display: none;
  color: white;
  max-width: 600px;
}

.showcase-content {
  padding: 60px 40px;
}

.showcase-title {
  font-size: 3.5rem;
  font-weight: 700;
  margin: 0 0 20px 0;
  line-height: 1.2;
  background: linear-gradient(45deg, #ffffff, #e8f4fd);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.showcase-subtitle {
  font-size: 1.25rem;
  margin: 0 0 40px 0;
  opacity: 0.9;
  line-height: 1.6;
}

.showcase-features {
  margin: 50px 0;
}

.feature-item {
  display: flex;
  align-items: center;
  margin: 20px 0;
  font-size: 1.1rem;
  opacity: 0.95;
}

.feature-icon {
  font-size: 1.5rem;
  margin-right: 15px;
  width: 40px;
  text-align: center;
}

.showcase-stats {
  display: flex;
  gap: 40px;
  margin-top: 60px;
}

.stat-item {
  text-align: center;
}

.stat-number {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 8px;
  background: linear-gradient(45deg, #ffffff, #e8f4fd);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.stat-label {
  font-size: 0.9rem;
  opacity: 0.8;
}

/* 右侧表单区域 */
.login-form-section {
  flex: 0 0 auto;
  width: 100%;
  max-width: 500px;
}

.login-card {
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
  border-radius: 16px;
  overflow: hidden;
  backdrop-filter: blur(20px);
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.card-header {
  text-align: center;
  padding: 20px 0 10px 0;
}

.card-header h2 {
  margin: 0 0 8px 0;
  color: #333;
  font-weight: 600;
  font-size: 1.75rem;
}

.header-subtitle {
  margin: 0;
  color: #666;
  font-size: 0.95rem;
  opacity: 0.8;
}

.login-form {
  padding: 10px 0;
}

.form-input {
  border-radius: 8px;
}

.login-button {
  width: 100%;
  height: 52px;
  font-size: 16px;
  font-weight: 500;
  border-radius: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
}

.login-button:hover {
  background: linear-gradient(135deg, #5a6fd8 0%, #6a4190 100%);
  transform: translateY(-1px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.3);
}

.login-footer {
  text-align: center;
  color: #666;
  font-size: 14px;
}

.register-link {
  color: #667eea;
  text-decoration: none;
  margin-left: 5px;
  font-weight: 500;
}

.register-link:hover {
  text-decoration: underline;
  color: #5a6fd8;
}

/* 大屏幕优化 (1200px+) */
@media (min-width: 1200px) {
  .login-showcase {
    display: block;
  }
  
  .login-container {
    max-width: 1400px;
  }
  
  .showcase-title {
    font-size: 4rem;
  }
  
  .showcase-subtitle {
    font-size: 1.4rem;
  }
}

/* 超大屏幕优化 (1600px+) */
@media (min-width: 1600px) {
  .login-container {
    max-width: 1600px;
    gap: 80px;
  }
  
  .showcase-content {
    padding: 80px 60px;
  }
  
  .showcase-title {
    font-size: 4.5rem;
  }
  
  .showcase-features {
    margin: 60px 0;
  }
  
  .feature-item {
    margin: 25px 0;
    font-size: 1.2rem;
  }
  
  .showcase-stats {
    gap: 60px;
    margin-top: 80px;
  }
  
  .stat-number {
    font-size: 3rem;
  }
}

/* 4K屏幕优化 (2000px+) */
@media (min-width: 2000px) {
  .login-container {
    max-width: 1800px;
    gap: 100px;
  }
  
  .showcase-title {
    font-size: 5rem;
  }
  
  .showcase-subtitle {
    font-size: 1.5rem;
  }
  
  .login-card {
    transform: scale(1.1);
  }
}

/* 中等屏幕 (768px - 1199px) */
@media (max-width: 1199px) {
  .login-container {
    justify-content: center;
  }
  
  .login-form-section {
    max-width: 500px;
  }
}

/* 平板和移动端响应式 */
@media (max-width: 768px) {
  .login-page {
    padding: 16px 12px;
  }
  
  .login-container {
    gap: 0;
  }
  
  .card-header h2 {
    font-size: 1.5rem;
  }
  
  .header-subtitle {
    font-size: 0.9rem;
  }
  
  .login-button {
    height: 48px;
    font-size: 15px;
  }
}

@media (max-width: 480px) {
  .login-page {
    padding: 12px 8px;
  }
  
  .card-header {
    padding: 15px 0 8px 0;
  }
  
  .card-header h2 {
    font-size: 1.25rem;
  }
  
  .header-subtitle {
    font-size: 0.85rem;
  }
  
  .login-footer {
    font-size: 13px;
  }
}

/* 超小屏幕优化 */
@media (max-width: 360px) {
  .login-page {
    padding: 8px 4px;
  }
  
  :deep(.el-card__body) {
    padding: 16px;
  }
  
  :deep(.el-form-item) {
    margin-bottom: 18px;
  }
  
  :deep(.el-input__wrapper) {
    padding: 8px 12px;
  }
}

/* 横屏优化 */
@media (max-height: 700px) and (orientation: landscape) and (min-width: 1200px) {
  .showcase-content {
    padding: 40px;
  }
  
  .showcase-title {
    font-size: 3rem;
  }
  
  .showcase-subtitle {
    font-size: 1.1rem;
  }
  
  .showcase-features {
    margin: 30px 0;
  }
  
  .showcase-stats {
    margin-top: 40px;
    gap: 30px;
  }
  
  .stat-number {
    font-size: 2rem;
  }
}

/* 深色模式适配 */
@media (prefers-color-scheme: dark) {
  .login-card {
    background: rgba(28, 28, 30, 0.95);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }
  
  .card-header h2 {
    color: #ffffff;
  }
  
  .header-subtitle {
    color: #8e8e93;
  }
  
  .login-footer {
    color: #8e8e93;
  }
}

/* 高分辨率屏幕优化 */
@media (-webkit-min-device-pixel-ratio: 2), (min-resolution: 192dpi) {
  .login-card {
    border: 0.5px solid rgba(255, 255, 255, 0.3);
  }
  
  .showcase-title,
  .stat-number {
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }
}
</style>