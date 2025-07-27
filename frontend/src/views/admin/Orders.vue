<template>
  <div class="admin-orders">
    <div class="page-header">
      <h2>订单管理</h2>
      <p>管理系统中的所有订单</p>
    </div>

    <!-- 搜索和筛选 -->
    <el-card class="filter-card" shadow="never">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="订单状态">
          <el-select v-model="searchForm.status" placeholder="全部状态" clearable @change="loadOrders">
            <el-option label="待处理" value="pending" />
            <el-option label="打印中" value="printing" />
            <el-option label="已完成" value="completed" />
            <el-option label="已发货" value="shipped" />
            <el-option label="已取消" value="cancelled" />
          </el-select>
        </el-form-item>
        <el-form-item label="用户ID">
          <el-input
            v-model="searchForm.user_id"
            placeholder="用户ID"
            style="width: 200px;"
            clearable
            @change="loadOrders"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadOrders">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 订单列表 -->
    <el-card shadow="never">
      <el-table
        :data="orders"
        v-loading="loading"
        style="width: 100%"
        row-key="id"
        :expand-row-keys="expandedRows"
        @expand-change="handleExpandChange"
      >
        <el-table-column type="expand">
          <template #default="{ row }">
            <div class="order-details">
              <h4>订单详情</h4>
              <el-descriptions :column="2" border>
                <el-descriptions-item label="订单ID">{{ row.id }}</el-descriptions-item>
                <el-descriptions-item label="用户">{{ row.user?.username }}</el-descriptions-item>
                <el-descriptions-item label="邮箱">{{ row.user?.email }}</el-descriptions-item>
                <el-descriptions-item label="手机">{{ row.user?.phone || '-' }}</el-descriptions-item>
                <el-descriptions-item label="总金额">¥{{ row.total_amount.toFixed(2) }}</el-descriptions-item>
                <el-descriptions-item label="运费">¥{{ row.shipping_cost.toFixed(2) }}</el-descriptions-item>
                <el-descriptions-item label="物流单号">{{ row.tracking_number || '-' }}</el-descriptions-item>
                <el-descriptions-item label="创建时间">{{ formatDateTime(row.created_at) }}</el-descriptions-item>
              </el-descriptions>
              
              <h4 style="margin-top: 20px;">订单项目</h4>
              <el-table :data="row.items" size="small">
                <el-table-column prop="file.filename" label="文件名" />
                <el-table-column prop="material.name" label="材料" />
                <el-table-column prop="quality" label="质量">
                  <template #default="{ row: item }">
                    {{ getQualityText(item.quality) }}
                  </template>
                </el-table-column>
                <el-table-column prop="quantity" label="数量" />
                <el-table-column prop="support_enabled" label="支撑">
                  <template #default="{ row: item }">
                    <el-tag :type="item.support_enabled ? 'success' : 'info'" size="small">
                      {{ item.support_enabled ? '启用' : '禁用' }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="price" label="价格">
                  <template #default="{ row: item }">
                    ¥{{ item.price.toFixed(2) }}
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="id" label="订单ID" width="120">
          <template #default="{ row }">
            <span class="order-id">{{ row.id.substring(0, 8) }}...</span>
          </template>
        </el-table-column>
        <el-table-column prop="user.username" label="用户" min-width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="total_amount" label="总金额" width="100">
          <template #default="{ row }">
            ¥{{ row.total_amount.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="items" label="项目数" width="80">
          <template #default="{ row }">
            {{ row.items?.length || 0 }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="120">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-dropdown @command="handleCommand" trigger="click">
              <el-button type="primary" size="small">
                操作 <el-icon><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item :command="{ action: 'printing', order: row }" v-if="row.status === 'pending'">
                    开始打印
                  </el-dropdown-item>
                  <el-dropdown-item :command="{ action: 'completed', order: row }" v-if="row.status === 'printing'">
                    标记完成
                  </el-dropdown-item>
                  <el-dropdown-item :command="{ action: 'shipped', order: row }" v-if="row.status === 'completed'">
                    标记发货
                  </el-dropdown-item>
                  <el-dropdown-item :command="{ action: 'cancelled', order: row }" v-if="['pending', 'printing'].includes(row.status)" divided>
                    取消订单
                  </el-dropdown-item>
                  <el-dropdown-item :command="{ action: 'tracking', order: row }" v-if="row.status === 'shipped'">
                    更新物流
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.limit"
          :page-sizes="[20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadOrders"
          @current-change="loadOrders"
        />
      </div>
    </el-card>

    <!-- 物流单号更新对话框 -->
    <el-dialog
      v-model="trackingDialog.visible"
      title="更新物流单号"
      width="400px"
    >
      <el-form :model="trackingDialog.form" label-width="80px">
        <el-form-item label="物流单号">
          <el-input v-model="trackingDialog.form.tracking_number" placeholder="请输入物流单号" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="trackingDialog.visible = false">取消</el-button>
          <el-button type="primary" @click="updateTracking">确定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'
import axios from 'axios'

export default {
  name: 'AdminOrders',
  components: {
    ArrowDown
  },
  setup() {
    const loading = ref(false)
    const orders = ref([])
    const expandedRows = ref([])
    
    const searchForm = reactive({
      status: '',
      user_id: ''
    })
    
    const pagination = reactive({
      page: 1,
      limit: 20,
      total: 0
    })

    const trackingDialog = reactive({
      visible: false,
      order: null,
      form: {
        tracking_number: ''
      }
    })

    const loadOrders = async () => {
      try {
        loading.value = true
        const params = {
          page: pagination.page,
          limit: pagination.limit,
          ...searchForm
        }
        
        const response = await axios.get('/api/admin/orders', { params })
        orders.value = response.data.orders
        pagination.total = response.data.total
      } catch (error) {
        console.error('加载订单列表失败:', error)
        ElMessage.error('加载订单列表失败')
      } finally {
        loading.value = false
      }
    }

    const resetSearch = () => {
      searchForm.status = ''
      searchForm.user_id = ''
      pagination.page = 1
      loadOrders()
    }

    const handleExpandChange = (row, expanded) => {
      if (expanded) {
        expandedRows.value.push(row.id)
      } else {
        const index = expandedRows.value.indexOf(row.id)
        if (index > -1) {
          expandedRows.value.splice(index, 1)
        }
      }
    }

    const handleCommand = async ({ action, order }) => {
      if (action === 'tracking') {
        trackingDialog.order = order
        trackingDialog.form.tracking_number = order.tracking_number || ''
        trackingDialog.visible = true
        return
      }

      const actionMap = {
        printing: { status: 'printing', text: '开始打印' },
        completed: { status: 'completed', text: '标记完成' },
        shipped: { status: 'shipped', text: '标记发货' },
        cancelled: { status: 'cancelled', text: '取消订单' }
      }

      const { status, text } = actionMap[action]
      
      try {
        await ElMessageBox.confirm(
          `确定要${text}订单 "${order.id.substring(0, 8)}..." 吗？`,
          `${text}`,
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )

        await axios.put(`/api/admin/orders/${order.id}/status`, { status })
        ElMessage.success(`订单${text}成功`)
        loadOrders()
      } catch (error) {
        if (error !== 'cancel') {
          console.error(`${text}失败:`, error)
          ElMessage.error(`${text}失败`)
        }
      }
    }

    const updateTracking = async () => {
      try {
        await axios.put(`/api/admin/orders/${trackingDialog.order.id}/status`, {
          status: 'shipped',
          tracking_number: trackingDialog.form.tracking_number
        })
        ElMessage.success('物流单号更新成功')
        trackingDialog.visible = false
        loadOrders()
      } catch (error) {
        console.error('更新物流单号失败:', error)
        ElMessage.error('更新物流单号失败')
      }
    }

    const getStatusType = (status) => {
      const statusMap = {
        'pending': 'warning',
        'printing': 'primary',
        'completed': 'success',
        'shipped': 'info',
        'cancelled': 'danger'
      }
      return statusMap[status] || ''
    }

    const getStatusText = (status) => {
      const statusMap = {
        'pending': '待处理',
        'printing': '打印中',
        'completed': '已完成',
        'shipped': '已发货',
        'cancelled': '已取消'
      }
      return statusMap[status] || status
    }

    const getQualityText = (quality) => {
      const qualityMap = {
        'fast': '快速',
        'standard': '标准',
        'high': '精细'
      }
      return qualityMap[quality] || quality
    }

    const formatDate = (dateString) => {
      const date = new Date(dateString)
      return date.toLocaleDateString('zh-CN')
    }

    const formatDateTime = (dateString) => {
      const date = new Date(dateString)
      return date.toLocaleString('zh-CN')
    }

    onMounted(() => {
      loadOrders()
    })

    return {
      loading,
      orders,
      expandedRows,
      searchForm,
      pagination,
      trackingDialog,
      loadOrders,
      resetSearch,
      handleExpandChange,
      handleCommand,
      updateTracking,
      getStatusType,
      getStatusText,
      getQualityText,
      formatDate,
      formatDateTime
    }
  }
}
</script>

<style scoped>
.admin-orders {
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
}

.page-header h2 {
  margin: 0 0 8px 0;
  color: #333;
  font-size: 1.5rem;
}

.page-header p {
  margin: 0;
  color: #666;
  font-size: 0.9rem;
}

.filter-card {
  margin-bottom: 16px;
}

.search-form {
  margin: 0;
}

.order-id {
  font-family: monospace;
  font-size: 0.85rem;
  color: #666;
}

.order-details {
  padding: 20px;
  background: #f8f9fa;
}

.order-details h4 {
  margin: 0 0 16px 0;
  color: #333;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

.dialog-footer {
  text-align: right;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .admin-orders {
    padding: 16px;
  }
  
  .search-form {
    display: block;
  }
  
  .search-form .el-form-item {
    display: block;
    margin-bottom: 16px;
  }
  
  .search-form .el-input,
  .search-form .el-select {
    width: 100%;
  }
  
  .order-details {
    padding: 15px;
  }
}
</style>