<template>
  <div class="admin-users">
    <div class="page-header">
      <h2>用户管理</h2>
      <p>管理系统中的所有用户账户</p>
    </div>

    <!-- 搜索和筛选 -->
    <el-card class="filter-card" shadow="never">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="搜索">
          <el-input
            v-model="searchForm.search"
            placeholder="用户名或邮箱"
            style="width: 200px;"
            clearable
            @change="loadUsers"
          />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="searchForm.role" placeholder="全部角色" clearable @change="loadUsers">
            <el-option label="普通用户" value="user" />
            <el-option label="管理员" value="admin" />
            <el-option label="超级管理员" value="super_admin" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部状态" clearable @change="loadUsers">
            <el-option label="活跃" value="active" />
            <el-option label="未激活" value="inactive" />
            <el-option label="已封禁" value="banned" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadUsers">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 用户列表 -->
    <el-card shadow="never">
      <el-table
        :data="users"
        v-loading="loading"
        style="width: 100%"
        row-key="id"
      >
        <el-table-column prop="username" label="用户名" min-width="120" />
        <el-table-column prop="email" label="邮箱" min-width="180" />
        <el-table-column prop="phone" label="手机号" min-width="120">
          <template #default="{ row }">
            {{ row.phone || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="role" label="角色" width="100">
          <template #default="{ row }">
            <el-tag :type="getRoleType(row.role)" size="small">
              {{ getRoleText(row.role) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="120">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-dropdown @command="handleCommand" trigger="click">
              <el-button type="primary" size="small">
                操作 <el-icon><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item :command="{ action: 'activate', user: row }" v-if="row.status !== 'active'">
                    激活用户
                  </el-dropdown-item>
                  <el-dropdown-item :command="{ action: 'deactivate', user: row }" v-if="row.status === 'active'">
                    停用用户
                  </el-dropdown-item>
                  <el-dropdown-item :command="{ action: 'ban', user: row }" v-if="row.status !== 'banned'" divided>
                    封禁用户
                  </el-dropdown-item>
                  <el-dropdown-item :command="{ action: 'unban', user: row }" v-if="row.status === 'banned'">
                    解封用户
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
          @size-change="loadUsers"
          @current-change="loadUsers"
        />
      </div>
    </el-card>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'
import axios from 'axios'

export default {
  name: 'AdminUsers',
  components: {
    ArrowDown
  },
  setup() {
    const loading = ref(false)
    const users = ref([])
    
    const searchForm = reactive({
      search: '',
      role: '',
      status: ''
    })
    
    const pagination = reactive({
      page: 1,
      limit: 20,
      total: 0
    })

    const loadUsers = async () => {
      try {
        loading.value = true
        const params = {
          page: pagination.page,
          limit: pagination.limit,
          ...searchForm
        }
        
        const response = await axios.get('/api/admin/users', { params })
        users.value = response.data.users
        pagination.total = response.data.total
      } catch (error) {
        console.error('加载用户列表失败:', error)
        ElMessage.error('加载用户列表失败')
      } finally {
        loading.value = false
      }
    }

    const resetSearch = () => {
      searchForm.search = ''
      searchForm.role = ''
      searchForm.status = ''
      pagination.page = 1
      loadUsers()
    }

    const handleCommand = async ({ action, user }) => {
      const actionMap = {
        activate: { status: 'active', text: '激活' },
        deactivate: { status: 'inactive', text: '停用' },
        ban: { status: 'banned', text: '封禁' },
        unban: { status: 'active', text: '解封' }
      }

      const { status, text } = actionMap[action]
      
      try {
        await ElMessageBox.confirm(
          `确定要${text}用户 "${user.username}" 吗？`,
          `${text}用户`,
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )

        await axios.put(`/api/admin/users/${user.id}/status`, { status })
        ElMessage.success(`用户${text}成功`)
        loadUsers()
      } catch (error) {
        if (error !== 'cancel') {
          console.error(`${text}用户失败:`, error)
          ElMessage.error(`${text}用户失败`)
        }
      }
    }

    const getRoleType = (role) => {
      const roleMap = {
        'user': '',
        'admin': 'warning',
        'super_admin': 'danger'
      }
      return roleMap[role] || ''
    }

    const getRoleText = (role) => {
      const roleMap = {
        'user': '用户',
        'admin': '管理员',
        'super_admin': '超管'
      }
      return roleMap[role] || role
    }

    const getStatusType = (status) => {
      const statusMap = {
        'active': 'success',
        'inactive': 'warning',
        'banned': 'danger'
      }
      return statusMap[status] || ''
    }

    const getStatusText = (status) => {
      const statusMap = {
        'active': '活跃',
        'inactive': '未激活',
        'banned': '已封禁'
      }
      return statusMap[status] || status
    }

    const formatDate = (dateString) => {
      const date = new Date(dateString)
      return date.toLocaleDateString('zh-CN')
    }

    onMounted(() => {
      loadUsers()
    })

    return {
      loading,
      users,
      searchForm,
      pagination,
      loadUsers,
      resetSearch,
      handleCommand,
      getRoleType,
      getRoleText,
      getStatusType,
      getStatusText,
      formatDate
    }
  }
}
</script>

<style scoped>
.admin-users {
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

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .admin-users {
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
}
</style>