<template>
  <div class="admin-dashboard">
    <div class="dashboard-header">
      <h1>管理后台</h1>
      <p>欢迎回来，{{ userStore.userInfo.username }}</p>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon users">
              <el-icon size="30"><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-number">{{ stats.total_users || 0 }}</div>
              <div class="stat-label">总用户数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon orders">
              <el-icon size="30"><Document /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-number">{{ stats.total_orders || 0 }}</div>
              <div class="stat-label">总订单数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon revenue">
              <el-icon size="30"><Money /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-number">¥{{ (stats.total_revenue || 0).toFixed(2) }}</div>
              <div class="stat-label">总收入</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon files">
              <el-icon size="30"><Folder /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-number">{{ stats.total_files || 0 }}</div>
              <div class="stat-label">总文件数</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 今日数据 -->
    <el-row :gutter="20" class="today-stats">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>今日数据</span>
              <el-button @click="refreshStats" type="primary" size="small">
                <el-icon><Refresh /></el-icon>
                刷新
              </el-button>
            </div>
          </template>
          
          <el-row :gutter="20">
            <el-col :xs="12" :sm="6">
              <div class="today-stat">
                <div class="today-number">{{ stats.today_orders || 0 }}</div>
                <div class="today-label">今日订单</div>
              </div>
            </el-col>
            <el-col :xs="12" :sm="6">
              <div class="today-stat">
                <div class="today-number">¥{{ (stats.today_revenue || 0).toFixed(2) }}</div>
                <div class="today-label">今日收入</div>
              </div>
            </el-col>
            <el-col :xs="12" :sm="6">
              <div class="today-stat">
                <div class="today-number">{{ stats.active_orders || 0 }}</div>
                <div class="today-label">进行中订单</div>
              </div>
            </el-col>
            <el-col :xs="12" :sm="6">
              <div class="today-stat">
                <div class="today-number">{{ stats.completed_orders || 0 }}</div>
                <div class="today-label">已完成订单</div>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>

    <!-- 最近订单和热门材料 -->
    <el-row :gutter="20" class="data-tables">
      <el-col :xs="24" :lg="14">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>最近订单</span>
              <el-button @click="$router.push('/admin/orders')" text type="primary">
                查看全部
              </el-button>
            </div>
          </template>
          
          <el-table :data="stats.recent_orders || []" style="width: 100%" size="small">
            <el-table-column prop="id" label="订单ID" width="120">
              <template #default="{ row }">
                <span class="order-id">{{ row.id.substring(0, 8) }}...</span>
              </template>
            </el-table-column>
            <el-table-column prop="user.username" label="用户" />
            <el-table-column prop="status" label="状态">
              <template #default="{ row }">
                <el-tag :type="getOrderStatusType(row.status)" size="small">
                  {{ getOrderStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="total_amount" label="金额">
              <template #default="{ row }">
                ¥{{ row.total_amount.toFixed(2) }}
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="100">
              <template #default="{ row }">
                {{ formatDate(row.created_at) }}
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :lg="10">
        <el-card>
          <template #header>
            <span>热门材料</span>
          </template>
          
          <div class="material-stats">
            <div 
              v-for="material in stats.top_materials || []" 
              :key="material.name"
              class="material-item"
            >
              <div class="material-name">{{ material.name }}</div>
              <div class="material-count">{{ material.count }} 次使用</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { User, Document, Money, Folder, Refresh } from '@element-plus/icons-vue'
import axios from 'axios'

export default {
  name: 'AdminDashboard',
  components: {
    User,
    Document,
    Money,
    Folder,
    Refresh
  },
  setup() {
    const userStore = useUserStore()
    const stats = ref({})
    const loading = ref(false)

    const loadStats = async () => {
      try {
        loading.value = true
        const response = await axios.get('/api/admin/dashboard')
        stats.value = response.data
      } catch (error) {
        console.error('加载统计数据失败:', error)
        ElMessage.error('加载统计数据失败')
      } finally {
        loading.value = false
      }
    }

    const refreshStats = () => {
      loadStats()
      ElMessage.success('数据已刷新')
    }

    const getOrderStatusType = (status) => {
      const statusMap = {
        'pending': 'warning',
        'printing': 'primary',
        'completed': 'success',
        'shipped': 'info',
        'cancelled': 'danger'
      }
      return statusMap[status] || 'info'
    }

    const getOrderStatusText = (status) => {
      const statusMap = {
        'pending': '待处理',
        'printing': '打印中',
        'completed': '已完成',
        'shipped': '已发货',
        'cancelled': '已取消'
      }
      return statusMap[status] || status
    }

    const formatDate = (dateString) => {
      const date = new Date(dateString)
      return date.toLocaleDateString('zh-CN')
    }

    onMounted(() => {
      loadStats()
    })

    return {
      userStore,
      stats,
      loading,
      refreshStats,
      getOrderStatusType,
      getOrderStatusText,
      formatDate
    }
  }
}
</script>

<style scoped>
.admin-dashboard {
  padding: 20px;
  background: #f5f7fa;
  min-height: 100vh;
}

.dashboard-header {
  margin-bottom: 24px;
}

.dashboard-header h1 {
  margin: 0 0 8px 0;
  color: #333;
  font-size: 2rem;
}

.dashboard-header p {
  margin: 0;
  color: #666;
  font-size: 1rem;
}

.stats-row {
  margin-bottom: 24px;
}

.stat-card {
  height: 120px;
  cursor: pointer;
  transition: transform 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
}

.stat-content {
  display: flex;
  align-items: center;
  height: 100%;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  color: white;
}

.stat-icon.users {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.orders {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.revenue {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.files {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-info {
  flex: 1;
}

.stat-number {
  font-size: 2rem;
  font-weight: bold;
  color: #333;
  line-height: 1;
}

.stat-label {
  font-size: 0.9rem;
  color: #666;
  margin-top: 4px;
}

.today-stats {
  margin-bottom: 24px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.today-stat {
  text-align: center;
  padding: 16px 0;
}

.today-number {
  font-size: 1.5rem;
  font-weight: bold;
  color: #409eff;
}

.today-label {
  font-size: 0.85rem;
  color: #666;
  margin-top: 4px;
}

.data-tables {
  margin-bottom: 24px;
}

.order-id {
  font-family: monospace;
  font-size: 0.85rem;
  color: #666;
}

.material-stats {
  padding: 16px 0;
}

.material-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.material-item:last-child {
  border-bottom: none;
}

.material-name {
  font-weight: 500;
  color: #333;
}

.material-count {
  color: #666;
  font-size: 0.9rem;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .admin-dashboard {
    padding: 15px;
  }

  .dashboard-header h1 {
    font-size: 1.5rem;
  }

  .stat-card {
    height: 100px;
    margin-bottom: 16px;
  }

  .stat-icon {
    width: 50px;
    height: 50px;
    margin-right: 12px;
  }

  .stat-number {
    font-size: 1.5rem;
  }

  .today-number {
    font-size: 1.2rem;
  }
}

@media (max-width: 480px) {
  .admin-dashboard {
    padding: 10px;
  }

  .stat-content {
    flex-direction: column;
    text-align: center;
  }

  .stat-icon {
    margin-right: 0;
    margin-bottom: 8px;
  }

  .stat-number {
    font-size: 1.3rem;
  }
}
</style>