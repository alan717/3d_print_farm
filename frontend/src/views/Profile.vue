<template>
  <div class="profile-page">
    <el-row :gutter="20">
      <el-col :span="24">
        <h1>个人中心</h1>
      </el-col>
    </el-row>
    
    <el-row :gutter="20">
      <el-col :span="12">
        <el-card class="profile-card">
          <template #header>
            <span>基本信息</span>
          </template>
          
          <el-form
            :model="profileForm"
            :rules="rules"
            ref="profileFormRef"
            label-width="80px"
          >
            <el-form-item label="用户名" prop="username">
              <el-input v-model="profileForm.username" />
            </el-form-item>
            
            <el-form-item label="邮箱">
              <el-input v-model="profileForm.email" disabled />
            </el-form-item>
            
            <el-form-item label="手机号" prop="phone">
              <el-input v-model="profileForm.phone" />
            </el-form-item>
            
            <el-form-item label="地址" prop="address">
              <el-input
                v-model="profileForm.address"
                type="textarea"
                :rows="3"
                placeholder="请输入收货地址"
              />
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="updateProfile" :loading="loading">
                保存修改
              </el-button>
              <el-button @click="resetForm">重置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card class="stats-card">
          <template #header>
            <span>账户统计</span>
          </template>
          
          <div class="stats-grid">
            <div class="stat-item">
              <div class="stat-number">{{ userStats.totalOrders }}</div>
              <div class="stat-label">总订单数</div>
            </div>
            <div class="stat-item">
              <div class="stat-number">{{ userStats.completedOrders }}</div>
              <div class="stat-label">已完成订单</div>
            </div>
            <div class="stat-item">
              <div class="stat-number">¥{{ userStats.totalSpent }}</div>
              <div class="stat-label">累计消费</div>
            </div>
            <div class="stat-item">
              <div class="stat-number">{{ userStats.uploadedFiles }}</div>
              <div class="stat-label">上传文件数</div>
            </div>
          </div>
        </el-card>
        
        <el-card class="recent-orders-card" style="margin-top: 20px;">
          <template #header>
            <span>最近订单</span>
          </template>
          
          <el-table :data="recentOrders" style="width: 100%">
            <el-table-column prop="id" label="订单号" width="120">
              <template #default="scope">
                {{ scope.row.id.substring(0, 8) }}...
              </template>
            </el-table-column>
            <el-table-column prop="total_amount" label="金额" width="80">
              <template #default="scope">
                ¥{{ scope.row.total_amount }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态">
              <template #default="scope">
                <el-tag :type="getStatusType(scope.row.status)" size="small">
                  {{ getStatusText(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
          
          <div style="text-align: center; margin-top: 15px;">
            <router-link to="/orders">
              <el-button type="text">查看全部订单</el-button>
            </router-link>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import axios from 'axios'

export default {
  name: 'Profile',
  setup() {
    const userStore = useUserStore()
    const profileFormRef = ref()
    const loading = ref(false)
    
    const profileForm = reactive({
      username: '',
      email: '',
      phone: '',
      address: ''
    })
    
    const userStats = ref({
      totalOrders: 0,
      completedOrders: 0,
      totalSpent: 0,
      uploadedFiles: 0
    })
    
    const recentOrders = ref([])
    
    const rules = {
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
        { min: 3, max: 50, message: '用户名长度在3到50个字符', trigger: 'blur' }
      ],
      phone: [
        { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号格式', trigger: 'blur' }
      ]
    }
    
    const getStatusType = (status) => {
      const types = {
        pending: 'warning',
        printing: 'primary',
        completed: 'success',
        shipped: 'info',
        cancelled: 'danger'
      }
      return types[status] || 'info'
    }
    
    const getStatusText = (status) => {
      const texts = {
        pending: '待处理',
        printing: '打印中',
        completed: '已完成',
        shipped: '已发货',
        cancelled: '已取消'
      }
      return texts[status] || status
    }
    
    const loadUserProfile = async () => {
      try {
        // 从store获取用户信息
        const user = userStore.userInfo
        Object.assign(profileForm, {
          username: user.username,
          email: user.email,
          phone: user.phone || '',
          address: user.address || ''
        })
      } catch (error) {
        ElMessage.error('加载用户信息失败')
      }
    }
    
    const loadUserStats = async () => {
      try {
        // TODO: 实际API调用
        // const response = await axios.get('/api/user/stats')
        // userStats.value = response.data
        
        // 模拟数据
        userStats.value = {
          totalOrders: 5,
          completedOrders: 3,
          totalSpent: 125.50,
          uploadedFiles: 8
        }
      } catch (error) {
        console.error('加载用户统计失败:', error)
      }
    }
    
    const loadRecentOrders = async () => {
      try {
        // TODO: 实际API调用
        // const response = await axios.get('/api/orders?limit=5')
        // recentOrders.value = response.data
        
        // 模拟数据
        recentOrders.value = [
          {
            id: '550e8400-e29b-41d4-a716-446655440000',
            total_amount: 25.50,
            status: 'printing'
          },
          {
            id: '550e8400-e29b-41d4-a716-446655440001',
            total_amount: 18.00,
            status: 'completed'
          }
        ]
      } catch (error) {
        console.error('加载最近订单失败:', error)
      }
    }
    
    const updateProfile = async () => {
      if (!profileFormRef.value) return
      
      try {
        await profileFormRef.value.validate()
        loading.value = true
        
        // TODO: 实际API调用
        // await axios.put('/api/user/profile', profileForm)
        
        // 更新store中的用户信息
        userStore.userInfo.username = profileForm.username
        userStore.userInfo.phone = profileForm.phone
        userStore.userInfo.address = profileForm.address
        
        ElMessage.success('个人信息更新成功')
      } catch (error) {
        ElMessage.error('更新失败')
      } finally {
        loading.value = false
      }
    }
    
    const resetForm = () => {
      loadUserProfile()
    }
    
    onMounted(() => {
      loadUserProfile()
      loadUserStats()
      loadRecentOrders()
    })
    
    return {
      profileForm,
      rules,
      profileFormRef,
      loading,
      userStats,
      recentOrders,
      getStatusType,
      getStatusText,
      updateProfile,
      resetForm
    }
  }
}
</script>

<style scoped>
.profile-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.profile-card, .stats-card, .recent-orders-card {
  height: auto;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.stat-item {
  text-align: center;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
}

.stat-number {
  font-size: 24px;
  font-weight: bold;
  color: #409eff;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  color: #666;
}
</style>