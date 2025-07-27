<template>
  <div class="orders-page">
    <el-row>
      <el-col :span="24">
        <h1>我的订单</h1>
      </el-col>
    </el-row>
    
    <el-row>
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="orders-header">
              <span>订单列表</span>
              <el-select v-model="statusFilter" placeholder="筛选状态" style="width: 120px;">
                <el-option label="全部" value="" />
                <el-option label="待处理" value="pending" />
                <el-option label="打印中" value="printing" />
                <el-option label="已完成" value="completed" />
                <el-option label="已发货" value="shipped" />
              </el-select>
            </div>
          </template>
          
          <el-table :data="filteredOrders" style="width: 100%">
            <el-table-column prop="id" label="订单号" width="180" />
            <el-table-column prop="created_at" label="下单时间" width="160">
              <template #default="scope">
                {{ formatDate(scope.row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column prop="total_amount" label="订单金额" width="100">
              <template #default="scope">
                ¥{{ scope.row.total_amount }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="订单状态" width="100">
              <template #default="scope">
                <el-tag :type="getStatusType(scope.row.status)">
                  {{ getStatusText(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="tracking_number" label="物流单号" width="140">
              <template #default="scope">
                <span v-if="scope.row.tracking_number">
                  {{ scope.row.tracking_number }}
                </span>
                <span v-else class="text-muted">-</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200">
              <template #default="scope">
                <el-button type="primary" size="small" @click="viewOrder(scope.row)">
                  查看详情
                </el-button>
                <el-button 
                  v-if="scope.row.status === 'pending'" 
                  type="danger" 
                  size="small" 
                  @click="cancelOrder(scope.row)"
                >
                  取消订单
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          
          <div class="pagination-container">
            <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :page-sizes="[10, 20, 50]"
              :total="totalOrders"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
            />
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 订单详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="订单详情" width="60%">
      <div v-if="selectedOrder">
        <h3>订单信息</h3>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="订单号">{{ selectedOrder.id }}</el-descriptions-item>
          <el-descriptions-item label="下单时间">{{ formatDate(selectedOrder.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="订单状态">
            <el-tag :type="getStatusType(selectedOrder.status)">
              {{ getStatusText(selectedOrder.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="物流单号">
            {{ selectedOrder.tracking_number || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="订单金额">¥{{ selectedOrder.total_amount }}</el-descriptions-item>
          <el-descriptions-item label="运费">¥{{ selectedOrder.shipping_cost }}</el-descriptions-item>
        </el-descriptions>
        
        <h3 style="margin-top: 20px;">订单项目</h3>
        <el-table :data="selectedOrder.items" style="width: 100%">
          <el-table-column prop="file.filename" label="文件名" />
          <el-table-column prop="material.name" label="材料" />
          <el-table-column prop="quality" label="质量">
            <template #default="scope">
              {{ getQualityText(scope.row.quality) }}
            </template>
          </el-table-column>
          <el-table-column prop="quantity" label="数量" />
          <el-table-column prop="support_enabled" label="支撑">
            <template #default="scope">
              {{ scope.row.support_enabled ? '是' : '否' }}
            </template>
          </el-table-column>
          <el-table-column prop="price" label="价格">
            <template #default="scope">
              ¥{{ scope.row.price }}
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'

export default {
  name: 'Orders',
  setup() {
    const orders = ref([])
    const statusFilter = ref('')
    const currentPage = ref(1)
    const pageSize = ref(10)
    const totalOrders = ref(0)
    const detailDialogVisible = ref(false)
    const selectedOrder = ref(null)
    
    const filteredOrders = computed(() => {
      if (!statusFilter.value) return orders.value
      return orders.value.filter(order => order.status === statusFilter.value)
    })
    
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
    
    const getQualityText = (quality) => {
      const texts = {
        fast: '快速',
        standard: '标准',
        high: '精细'
      }
      return texts[quality] || quality
    }
    
    const formatDate = (dateString) => {
      return new Date(dateString).toLocaleString()
    }
    
    const loadOrders = async () => {
      try {
        // TODO: 实际API调用
        // const response = await axios.get('/api/orders')
        // orders.value = response.data
        
        // 模拟数据
        orders.value = [
          {
            id: '550e8400-e29b-41d4-a716-446655440000',
            created_at: '2024-01-20T10:30:00Z',
            total_amount: 25.50,
            shipping_cost: 5.00,
            status: 'printing',
            tracking_number: null,
            items: []
          }
        ]
        totalOrders.value = orders.value.length
      } catch (error) {
        ElMessage.error('加载订单列表失败')
      }
    }
    
    const viewOrder = (order) => {
      selectedOrder.value = order
      detailDialogVisible.value = true
    }
    
    const cancelOrder = async (order) => {
      try {
        await ElMessageBox.confirm(
          '确定要取消这个订单吗？',
          '取消订单',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        
        // TODO: 实际API调用
        // await axios.put(`/api/orders/${order.id}/status`, { status: 'cancelled' })
        
        ElMessage.success('订单已取消')
        loadOrders()
      } catch (error) {
        if (error !== 'cancel') {
          ElMessage.error('取消订单失败')
        }
      }
    }
    
    const handleSizeChange = (val) => {
      pageSize.value = val
      loadOrders()
    }
    
    const handleCurrentChange = (val) => {
      currentPage.value = val
      loadOrders()
    }
    
    onMounted(() => {
      loadOrders()
    })
    
    return {
      orders,
      filteredOrders,
      statusFilter,
      currentPage,
      pageSize,
      totalOrders,
      detailDialogVisible,
      selectedOrder,
      getStatusType,
      getStatusText,
      getQualityText,
      formatDate,
      viewOrder,
      cancelOrder,
      handleSizeChange,
      handleCurrentChange
    }
  }
}
</script>

<style scoped>
.orders-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.orders-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination-container {
  margin-top: 20px;
  text-align: center;
}

.text-muted {
  color: #999;
}
</style>