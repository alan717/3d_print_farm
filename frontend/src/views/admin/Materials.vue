<template>
  <div class="admin-materials">
    <div class="page-header">
      <h2>材料管理</h2>
      <p>管理3D打印材料和价格</p>
    </div>

    <!-- 操作按钮 -->
    <div class="action-bar">
      <el-button type="primary" @click="showCreateDialog">
        <el-icon><Plus /></el-icon>
        添加材料
      </el-button>
    </div>

    <!-- 材料列表 -->
    <el-card shadow="never">
      <el-table
        :data="materials"
        v-loading="loading"
        style="width: 100%"
        row-key="id"
      >
        <el-table-column prop="name" label="材料名称" min-width="120" />
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column prop="price_per_gram" label="价格(¥/克)" width="120">
          <template #default="{ row }">
            ¥{{ row.price_per_gram.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="properties" label="属性" min-width="200">
          <template #default="{ row }">
            <div v-if="getProperties(row.properties)" class="properties">
              <el-tag 
                v-for="(value, key) in getProperties(row.properties)" 
                :key="key" 
                size="small" 
                class="property-tag"
              >
                {{ key }}: {{ value }}
              </el-tag>
            </div>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="available" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.available ? 'success' : 'danger'" size="small">
              {{ row.available ? '可用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="120">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button size="small" @click="showEditDialog(row)">
                <el-icon><Edit /></el-icon>
              </el-button>
              <el-button size="small" type="danger" @click="deleteMaterial(row)">
                <el-icon><Delete /></el-icon>
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑材料对话框 -->
    <el-dialog
      v-model="materialDialog.visible"
      :title="materialDialog.isEdit ? '编辑材料' : '添加材料'"
      width="600px"
    >
      <el-form
        ref="materialFormRef"
        :model="materialDialog.form"
        :rules="materialDialog.rules"
        label-width="100px"
      >
        <el-form-item label="材料名称" prop="name">
          <el-input v-model="materialDialog.form.name" placeholder="请输入材料名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="materialDialog.form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入材料描述"
          />
        </el-form-item>
        <el-form-item label="价格(¥/克)" prop="price_per_gram">
          <el-input-number
            v-model="materialDialog.form.price_per_gram"
            :min="0"
            :step="0.01"
            :precision="2"
            style="width: 100%;"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch
            v-model="materialDialog.form.available"
            active-text="可用"
            inactive-text="停用"
          />
        </el-form-item>
        <el-form-item label="材料属性">
          <div class="properties-editor">
            <div 
              v-for="(prop, index) in materialDialog.form.properties" 
              :key="index"
              class="property-item"
            >
              <el-input
                v-model="prop.key"
                placeholder="属性名"
                style="width: 40%;"
              />
              <el-input
                v-model="prop.value"
                placeholder="属性值"
                style="width: 40%; margin-left: 10px;"
              />
              <el-button
                @click="removeProperty(index)"
                type="danger"
                size="small"
                style="margin-left: 10px;"
              >
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
            <el-button @click="addProperty" type="primary" size="small" plain>
              <el-icon><Plus /></el-icon>
              添加属性
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="materialDialog.visible = false">取消</el-button>
          <el-button type="primary" @click="saveMaterial">确定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import axios from 'axios'

export default {
  name: 'AdminMaterials',
  components: {
    Plus, Edit, Delete
  },
  setup() {
    const loading = ref(false)
    const materials = ref([])
    const materialFormRef = ref()
    
    const materialDialog = reactive({
      visible: false,
      isEdit: false,
      currentId: null,
      form: {
        name: '',
        description: '',
        price_per_gram: 0,
        available: true,
        properties: []
      },
      rules: {
        name: [
          { required: true, message: '请输入材料名称', trigger: 'blur' }
        ],
        price_per_gram: [
          { required: true, message: '请输入价格', trigger: 'blur' },
          { type: 'number', min: 0, message: '价格不能小于0', trigger: 'blur' }
        ]
      }
    })

    const loadMaterials = async () => {
      try {
        loading.value = true
        const response = await axios.get('/api/materials')
        materials.value = response.data
      } catch (error) {
        console.error('加载材料列表失败:', error)
        ElMessage.error('加载材料列表失败')
      } finally {
        loading.value = false
      }
    }

    const showCreateDialog = () => {
      resetForm()
      materialDialog.isEdit = false
      materialDialog.visible = true
    }

    const showEditDialog = (material) => {
      resetForm()
      materialDialog.isEdit = true
      materialDialog.currentId = material.id
      materialDialog.form.name = material.name
      materialDialog.form.description = material.description
      materialDialog.form.price_per_gram = material.price_per_gram
      materialDialog.form.available = material.available
      
      // Parse properties
      const properties = getProperties(material.properties)
      if (properties) {
        materialDialog.form.properties = Object.entries(properties).map(([key, value]) => ({
          key, value
        }))
      }
      
      materialDialog.visible = true
    }

    const resetForm = () => {
      materialDialog.form = {
        name: '',
        description: '',
        price_per_gram: 0,
        available: true,
        properties: []
      }
      materialDialog.currentId = null
    }

    const addProperty = () => {
      materialDialog.form.properties.push({ key: '', value: '' })
    }

    const removeProperty = (index) => {
      materialDialog.form.properties.splice(index, 1)
    }

    const saveMaterial = async () => {
      try {
        await materialFormRef.value.validate()
        
        // Convert properties array to object
        const propertiesObj = {}
        materialDialog.form.properties.forEach(prop => {
          if (prop.key && prop.value) {
            propertiesObj[prop.key] = prop.value
          }
        })

        const data = {
          ...materialDialog.form,
          properties: propertiesObj
        }

        if (materialDialog.isEdit) {
          await axios.put(`/api/admin/materials/${materialDialog.currentId}`, data)
          ElMessage.success('材料更新成功')
        } else {
          await axios.post('/api/admin/materials', data)
          ElMessage.success('材料创建成功')
        }

        materialDialog.visible = false
        loadMaterials()
      } catch (error) {
        console.error('保存材料失败:', error)
        ElMessage.error('保存材料失败')
      }
    }

    const deleteMaterial = async (material) => {
      try {
        await ElMessageBox.confirm(
          `确定要删除材料 "${material.name}" 吗？`,
          '删除材料',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )

        await axios.delete(`/api/admin/materials/${material.id}`)
        ElMessage.success('材料删除成功')
        loadMaterials()
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除材料失败:', error)
          ElMessage.error('删除材料失败')
        }
      }
    }

    const getProperties = (propertiesString) => {
      if (!propertiesString) return null
      try {
        return JSON.parse(propertiesString)
      } catch {
        return null
      }
    }

    const formatDate = (dateString) => {
      const date = new Date(dateString)
      return date.toLocaleDateString('zh-CN')
    }

    onMounted(() => {
      loadMaterials()
    })

    return {
      loading,
      materials,
      materialDialog,
      materialFormRef,
      loadMaterials,
      showCreateDialog,
      showEditDialog,
      addProperty,
      removeProperty,
      saveMaterial,
      deleteMaterial,
      getProperties,
      formatDate
    }
  }
}
</script>

<style scoped>
.admin-materials {
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

.action-bar {
  margin-bottom: 16px;
}

.properties {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.property-tag {
  margin: 0;
}

.properties-editor {
  width: 100%;
}

.property-item {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.dialog-footer {
  text-align: right;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .admin-materials {
    padding: 16px;
  }
  
  .property-item {
    flex-direction: column;
    align-items: stretch;
  }
  
  .property-item .el-input {
    width: 100% !important;
    margin: 5px 0 !important;
  }
  
  .property-item .el-button {
    margin: 5px 0 !important;
    width: 100%;
  }
}
</style>