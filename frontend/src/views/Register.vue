<template>
  <div class="register-page">
    <div class="register-container">
      <el-row justify="center">
        <el-col :span="24" :xs="22" :sm="20" :md="16" :lg="12" :xl="10">
          <el-card class="register-card">
          <template #header>
            <div class="card-header">
              <h2>用户注册</h2>
            </div>
          </template>
          
          <el-form
            :model="registerForm"
            :rules="rules"
            ref="registerFormRef"
            label-width="80px"
            @submit.prevent="handleRegister"
          >
            <el-form-item label="用户名" prop="username">
              <el-input
                v-model="registerForm.username"
                placeholder="请输入用户名"
                prefix-icon="User"
              />
            </el-form-item>
            
            <el-form-item label="邮箱" prop="email">
              <el-input
                v-model="registerForm.email"
                type="email"
                placeholder="请输入邮箱"
                prefix-icon="Message"
              />
            </el-form-item>
            
            <el-form-item label="密码" prop="password">
              <el-input
                v-model="registerForm.password"
                type="password"
                placeholder="请输入密码"
                prefix-icon="Lock"
                show-password
              />
            </el-form-item>
            
            <el-form-item label="确认密码" prop="confirmPassword">
              <el-input
                v-model="registerForm.confirmPassword"
                type="password"
                placeholder="请再次输入密码"
                prefix-icon="Lock"
                show-password
              />
            </el-form-item>
            
            <el-form-item label="手机号" prop="phone">
              <el-input
                v-model="registerForm.phone"
                placeholder="请输入手机号（可选）"
                prefix-icon="Phone"
              />
            </el-form-item>
            
            <el-form-item label="地址" prop="address">
              <el-input
                v-model="registerForm.address"
                type="textarea"
                placeholder="请输入收货地址（可选）"
                :rows="3"
              />
            </el-form-item>
            
            <el-form-item>
              <el-button
                type="primary"
                :loading="loading"
                @click="handleRegister"
                style="width: 100%;"
              >
                注册
              </el-button>
            </el-form-item>
            
            <el-form-item>
              <div class="register-footer">
                <span>已有账号？</span>
                <router-link to="/login" class="login-link">立即登录</router-link>
              </div>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

export default {
  name: 'Register',
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    const registerFormRef = ref()
    const loading = ref(false)
    
    const registerForm = reactive({
      username: '',
      email: '',
      password: '',
      confirmPassword: '',
      phone: '',
      address: ''
    })
    
    const validateConfirmPassword = (rule, value, callback) => {
      if (value !== registerForm.password) {
        callback(new Error('两次输入的密码不一致'))
      } else {
        callback()
      }
    }
    
    const rules = {
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
        { min: 3, max: 50, message: '用户名长度在3到50个字符', trigger: 'blur' }
      ],
      email: [
        { required: true, message: '请输入邮箱', trigger: 'blur' },
        { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, message: '密码长度至少6位', trigger: 'blur' }
      ],
      confirmPassword: [
        { required: true, message: '请确认密码', trigger: 'blur' },
        { validator: validateConfirmPassword, trigger: 'blur' }
      ],
      phone: [
        { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号格式', trigger: 'blur' }
      ]
    }
    
    const handleRegister = async () => {
      if (!registerFormRef.value) return
      
      try {
        await registerFormRef.value.validate()
        loading.value = true
        
        const { confirmPassword, ...registerData } = registerForm
        const result = await userStore.register(registerData)
        
        if (result.success) {
          ElMessage.success('注册成功，请登录')
          router.push('/login')
        } else {
          ElMessage.error(result.message)
        }
      } catch (error) {
        console.error('注册失败:', error)
      } finally {
        loading.value = false
      }
    }
    
    return {
      registerForm,
      rules,
      registerFormRef,
      loading,
      handleRegister
    }
  }
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px 16px;
}

.register-container {
  width: 100%;
  max-width: 600px;
}

.register-card {
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  overflow: hidden;
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.95);
}

.card-header {
  text-align: center;
  padding: 10px 0;
}

.card-header h2 {
  margin: 0;
  color: #333;
  font-weight: 600;
  font-size: 1.5rem;
}

.register-footer {
  text-align: center;
  color: #666;
  font-size: 14px;
}

.login-link {
  color: #409eff;
  text-decoration: none;
  margin-left: 5px;
  font-weight: 500;
}

.login-link:hover {
  text-decoration: underline;
  color: #337ecc;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .register-page {
    padding: 16px 12px;
  }
  
  .card-header h2 {
    font-size: 1.25rem;
  }
}

@media (max-width: 480px) {
  .register-page {
    padding: 12px 8px;
  }
  
  .card-header {
    padding: 8px 0;
  }
  
  .card-header h2 {
    font-size: 1.125rem;
  }
  
  .register-footer {
    font-size: 13px;
  }
}

/* 超小屏幕优化 */
@media (max-width: 360px) {
  .register-page {
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
@media (max-height: 600px) and (orientation: landscape) {
  .register-page {
    padding: 12px 16px;
  }
  
  .card-header {
    padding: 6px 0;
  }
  
  .card-header h2 {
    font-size: 1.125rem;
  }
  
  :deep(.el-form-item) {
    margin-bottom: 16px;
  }
}

/* 深色模式适配 */
@media (prefers-color-scheme: dark) {
  .register-card {
    background: rgba(28, 28, 30, 0.95);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }
  
  .card-header h2 {
    color: #ffffff;
  }
  
  .register-footer {
    color: #8e8e93;
  }
}
</style>