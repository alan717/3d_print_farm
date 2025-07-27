<template>
  <div class="admin-settings">
    <div class="page-header">
      <h2>系统设置</h2>
      <p>管理系统配置和参数</p>
    </div>

    <el-card shadow="never">
      <el-form
        ref="settingsFormRef"
        :model="settingsForm"
        label-width="150px"
        v-loading="loading"
      >
        <el-divider content-position="left">
          <h3>网站基本信息</h3>
        </el-divider>
        
        <el-form-item label="网站名称">
          <el-input
            v-model="settingsForm.site_name"
            placeholder="请输入网站名称"
            style="width: 400px;"
          />
        </el-form-item>
        
        <el-form-item label="网站描述">
          <el-input
            v-model="settingsForm.site_description"
            type="textarea"
            :rows="3"
            placeholder="请输入网站描述"
            style="width: 400px;"
          />
        </el-form-item>

        <el-divider content-position="left">
          <h3>文件上传设置</h3>
        </el-divider>
        
        <el-form-item label="最大文件大小">
          <el-input-number
            v-model="fileSizeMB"
            :min="1"
            :max="500"
            style="width: 200px;"
          />
          <span style="margin-left: 10px; color: #666;">MB</span>
        </el-form-item>
        
        <el-form-item label="支持的格式">
          <el-input
            v-model="settingsForm.supported_formats"
            placeholder="例如: STL,OBJ,PLY"
            style="width: 400px;"
          />
          <div class="form-tip">格式之间用逗号分隔</div>
        </el-form-item>

        <el-divider content-position="left">
          <h3>业务设置</h3>
        </el-divider>
        
        <el-form-item label="默认材料">
          <el-select
            v-model="settingsForm.default_material"
            placeholder="请选择默认材料"
            style="width: 200px;"
          >
            <el-option
              v-for="material in materials"
              :key="material.id"
              :label="material.name"
              :value="material.id.toString()"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="自动审核订单">
          <el-switch
            v-model="autoApproveOrders"
            active-text="启用"
            inactive-text="禁用"
          />
          <div class="form-tip">启用后新订单将自动进入打印流程</div>
        </el-form-item>

        <el-divider content-position="left">
          <h3>系统维护</h3>
        </el-divider>
        
        <el-form-item label="维护模式">
          <el-switch
            v-model="maintenanceMode"
            active-text="启用"
            inactive-text="禁用"
          />
          <div class="form-tip">启用后网站将显示维护页面</div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="saveSettings" :loading="saving">
            保存设置
          </el-button>
          <el-button @click="loadSettings">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 系统统计信息 -->
    <el-card shadow="never" style="margin-top: 24px;">
      <template #header>
        <h3>系统信息</h3>
      </template>
      
      <el-descriptions :column="2" border>
        <el-descriptions-item label="系统版本">v1.1.0-beta</el-descriptions-item>
        <el-descriptions-item label="数据库版本">PostgreSQL 14+</el-descriptions-item>
        <el-descriptions-item label="后端框架">Go + Gin</el-descriptions-item>
        <el-descriptions-item label="前端框架">Vue 3 + Element Plus</el-descriptions-item>
        <el-descriptions-item label="当前时间">{{ currentTime }}</el-descriptions-item>
        <el-descriptions-item label="服务器状态">
          <el-tag type="success">运行中</el-tag>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

export default {
  name: 'AdminSettings',
  setup() {
    const loading = ref(false)
    const saving = ref(false)
    const settingsFormRef = ref()
    const materials = ref([])
    const currentTime = ref('')
    
    const settingsForm = reactive({
      site_name: '',
      site_description: '',
      max_file_size: '',
      supported_formats: '',
      default_material: '',
      maintenance_mode: '',
      auto_approve_orders: ''
    })

    // 计算属性用于处理特殊格式
    const fileSizeMB = computed({
      get: () => {
        const bytes = parseInt(settingsForm.max_file_size) || 0
        return Math.round(bytes / (1024 * 1024))
      },
      set: (value) => {
        settingsForm.max_file_size = (value * 1024 * 1024).toString()
      }
    })

    const maintenanceMode = computed({
      get: () => settingsForm.maintenance_mode === 'true',
      set: (value) => {
        settingsForm.maintenance_mode = value.toString()
      }
    })

    const autoApproveOrders = computed({
      get: () => settingsForm.auto_approve_orders === 'true',
      set: (value) => {
        settingsForm.auto_approve_orders = value.toString()
      }
    })

    const loadSettings = async () => {
      try {
        loading.value = true
        const response = await axios.get('/api/admin/settings')
        
        // 将设置数组转换为对象
        response.data.forEach(setting => {
          settingsForm[setting.setting_key] = setting.setting_value
        })
      } catch (error) {
        console.error('加载系统设置失败:', error)
        ElMessage.error('加载系统设置失败')
      } finally {
        loading.value = false
      }
    }

    const loadMaterials = async () => {
      try {
        const response = await axios.get('/api/materials')
        materials.value = response.data
      } catch (error) {
        console.error('加载材料列表失败:', error)
      }
    }

    const saveSettings = async () => {
      try {
        saving.value = true
        
        // 逐个保存设置
        const savePromises = Object.entries(settingsForm).map(([key, value]) => {
          return axios.put(`/api/admin/settings/${key}`, { value })
        })
        
        await Promise.all(savePromises)
        ElMessage.success('设置保存成功')
      } catch (error) {
        console.error('保存设置失败:', error)
        ElMessage.error('保存设置失败')
      } finally {
        saving.value = false
      }
    }

    const updateCurrentTime = () => {
      currentTime.value = new Date().toLocaleString('zh-CN')
    }

    onMounted(() => {
      loadSettings()
      loadMaterials()
      updateCurrentTime()
      
      // 每秒更新时间
      setInterval(updateCurrentTime, 1000)
    })

    return {
      loading,
      saving,
      settingsFormRef,
      settingsForm,
      materials,
      currentTime,
      fileSizeMB,
      maintenanceMode,
      autoApproveOrders,
      loadSettings,
      saveSettings
    }
  }
}
</script>

<style scoped>
.admin-settings {
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

.form-tip {
  font-size: 0.8rem;
  color: #999;
  margin-top: 4px;
}

h3 {
  margin: 0;
  color: #333;
  font-size: 1.1rem;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .admin-settings {
    padding: 16px;
  }
  
  .el-form-item .el-input,
  .el-form-item .el-select,
  .el-form-item .el-input-number {
    width: 100% !important;
  }
}
</style>