<template>
  <div class="admin-logs">
    <div class="page-header">
      <h2>操作日志</h2>
      <p>查看管理员操作记录</p>
    </div>

    <!-- 搜索和筛选 -->
    <el-card class="filter-card" shadow="never">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="操作类型">
          <el-select v-model="searchForm.action" placeholder="全部操作" clearable @change="loadLogs">
            <el-option label="创建" value="create" />
            <el-option label="更新" value="update" />
            <el-option label="删除" value="delete" />
            <el-option label="状态变更" value="status_change" />
          </el-select>
        </el-form-item>
        <el-form-item label="资源类型">
          <el-select v-model="searchForm.resource_type" placeholder="全部资源" clearable @change="loadLogs">
            <el-option label="用户" value="user" />
            <el-option label="订单" value="order" />
            <el-option label="材料" value="material" />
            <el-option label="系统设置" value="system_setting" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadLogs">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 日志列表 -->
    <el-card shadow="never">
      <el-table
        :data="logs"
        v-loading="loading"
        style="width: 100%"
        row-key="id"
        :expand-row-keys="expandedRows"
        @expand-change="handleExpandChange"
      >
        <el-table-column type="expand">
          <template #default="{ row }">
            <div class="log-details">
              <el-descriptions :column="2" border>
                <el-descriptions-item label="操作ID">{{ row.id }}</el-descriptions-item>
                <el-descriptions-item label="管理员">{{ row.admin?.username || '未知' }}</el-descriptions-item>
                <el-descriptions-item label="IP地址">{{ row.ip_address || '-' }}</el-descriptions-item>
                <el-descriptions-item label="用户代理">
                  <div class="user-agent">{{ row.user_agent || '-' }}</div>
                </el-descriptions-item>
              </el-descriptions>
              
              <h4 style="margin-top: 20px;">操作详情</h4>
              <div class="details-content">
                <pre v-if="row.details">{{ formatDetails(row.details) }}</pre>
                <span v-else>无详细信息</span>
              </div>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="action" label="操作" width="120">
          <template #default="{ row }">
            <el-tag :type="getActionType(row.action)" size="small">
              {{ getActionText(row.action) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource_type" label="资源类型" width="100">
          <template #default="{ row }">
            <el-tag type="info" size="small">
              {{ getResourceText(row.resource_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource_id" label="资源ID" width="120">
          <template #default="{ row }">
            <span class="resource-id">{{ row.resource_id || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="admin.username" label="操作者" width="100">
          <template #default="{ row }">
            {{ row.admin?.username || '系统' }}
          </template>
        </el-table-column>
        <el-table-column prop="ip_address" label="IP地址" width="120">
          <template #default="{ row }">
            {{ row.ip_address || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="操作时间" width="150">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
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
          @size-change="loadLogs"
          @current-change="loadLogs"
        />
      </div>
    </el-card>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

export default {
  name: 'AdminLogs',
  setup() {
    const loading = ref(false)
    const logs = ref([])
    const expandedRows = ref([])
    
    const searchForm = reactive({
      action: '',
      resource_type: ''
    })
    
    const pagination = reactive({
      page: 1,
      limit: 50,
      total: 0
    })

    const loadLogs = async () => {
      try {
        loading.value = true
        const params = {
          page: pagination.page,
          limit: pagination.limit,
          ...searchForm
        }
        
        const response = await axios.get('/api/admin/logs', { params })
        logs.value = response.data.logs
        pagination.total = response.data.total
      } catch (error) {
        console.error('加载操作日志失败:', error)
        ElMessage.error('加载操作日志失败')
      } finally {
        loading.value = false
      }
    }

    const resetSearch = () => {
      searchForm.action = ''
      searchForm.resource_type = ''
      pagination.page = 1
      loadLogs()
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

    const getActionType = (action) => {
      const actionTypeMap = {
        'create': 'success',
        'create_material': 'success',
        'update': 'warning',
        'update_material': 'warning',
        'update_user_status': 'warning',
        'update_order_status': 'warning',
        'update_setting': 'warning',
        'delete': 'danger',
        'delete_material': 'danger'
      }
      return actionTypeMap[action] || 'info'
    }

    const getActionText = (action) => {
      const actionTextMap = {
        'create': '创建',
        'create_material': '创建材料',
        'update': '更新',
        'update_material': '更新材料',
        'update_user_status': '更新用户状态',
        'update_order_status': '更新订单状态',
        'update_setting': '更新设置',
        'delete': '删除',
        'delete_material': '删除材料'
      }
      return actionTextMap[action] || action
    }

    const getResourceText = (resourceType) => {
      const resourceTextMap = {
        'user': '用户',
        'order': '订单',
        'material': '材料',
        'system_setting': '设置'
      }
      return resourceTextMap[resourceType] || resourceType
    }

    const formatDetails = (details) => {
      if (!details) return ''
      try {
        // 如果是 JSON 字符串，尝试格式化
        const parsed = JSON.parse(details)
        return JSON.stringify(parsed, null, 2)
      } catch {
        // 如果不是 JSON，直接返回
        return details
      }
    }

    const formatDateTime = (dateString) => {
      const date = new Date(dateString)
      return date.toLocaleString('zh-CN')
    }

    onMounted(() => {
      loadLogs()
    })

    return {
      loading,
      logs,
      expandedRows,
      searchForm,
      pagination,
      loadLogs,
      resetSearch,
      handleExpandChange,
      getActionType,
      getActionText,
      getResourceText,
      formatDetails,
      formatDateTime
    }
  }
}
</script>

<style scoped>
.admin-logs {
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

.resource-id {
  font-family: monospace;
  font-size: 0.85rem;
  color: #666;
}

.log-details {
  padding: 20px;
  background: #f8f9fa;
}

.log-details h4 {
  margin: 0 0 16px 0;
  color: #333;
}

.user-agent {
  max-width: 300px;
  word-break: break-all;
  font-size: 0.8rem;
  color: #666;
}

.details-content {
  background: #fff;
  border: 1px solid #e6e6e6;
  border-radius: 4px;
  padding: 12px;
  max-height: 200px;
  overflow-y: auto;
}

.details-content pre {
  margin: 0;
  font-size: 0.8rem;
  color: #333;
  white-space: pre-wrap;
  word-break: break-all;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .admin-logs {
    padding: 16px;
  }
  
  .search-form {
    display: block;
  }
  
  .search-form .el-form-item {
    display: block;
    margin-bottom: 16px;
  }
  
  .search-form .el-select {
    width: 100%;
  }
  
  .log-details {
    padding: 15px;
  }
  
  .user-agent {
    max-width: 200px;
  }
}
</style>