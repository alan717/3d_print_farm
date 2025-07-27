<template>
  <div class="upload-page">
    <el-row :gutter="20">
      <el-col :span="24">
        <h1>上传3D模型</h1>
      </el-col>
    </el-row>
    
    <!-- File Upload Section -->
    <el-row :gutter="20">
      <el-col :span="24">
        <el-card class="upload-card">
          <template #header>
            <div class="upload-header">
              <span>文件上传</span>
              <div class="header-actions">
                <el-button @click="loadDemoModel" type="success" size="small" plain>
                  加载演示模型
                </el-button>
                <el-button @click="clearAllFiles" size="small" v-if="uploadedFiles.length > 0">
                  清空所有文件
                </el-button>
              </div>
            </div>
          </template>
          
          <el-upload
            class="upload-dragger"
            drag
            action="/api/files/upload"
            :headers="uploadHeaders"
            accept=".stl,.obj"
            multiple
            :before-upload="beforeUpload"
            :on-success="onUploadSuccess"
            :on-error="onUploadError"
            :show-file-list="false"
          >
            <el-icon class="el-icon--upload">
              <UploadFilled />
            </el-icon>
            <div class="el-upload__text">
              将文件拖到此处，或<em>点击上传</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                支持STL、OBJ格式，文件大小不超过50MB，支持多文件选择
              </div>
            </template>
          </el-upload>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- File List Section -->
    <el-row :gutter="20" style="margin-top: 20px;" v-if="uploadedFiles.length > 0">
      <el-col :span="24">
        <el-card class="file-list-card">
          <template #header>
            <div class="file-list-header">
              <span>已上传文件 ({{ uploadedFiles.length }})</span>
              <div class="batch-actions">
                <el-button @click="addAllToCart" type="primary" size="small" :disabled="!canAddAllToCart">
                  全部加入购物车
                </el-button>
              </div>
            </div>
          </template>
          
          <div class="file-tabs">
            <el-tabs v-model="activeFileTab" type="card">
              <el-tab-pane 
                v-for="(file, index) in uploadedFiles" 
                :key="file.id" 
                :label="getTabLabel(file, index)"
                :name="file.id"
              >
                <div class="file-content">
                  <el-row :gutter="20">
                    <!-- 3D Preview -->
                    <el-col :xs="24" :sm="24" :md="12">
                      <div class="preview-section">
                        <div class="preview-header">
                          <h3>{{ file.filename || file.original_name }}</h3>
                          <el-button @click="removeFile(file.id)" type="danger" size="small" text>
                            <el-icon><Delete /></el-icon>
                          </el-button>
                        </div>
                        
                        <div class="three-preview">
                          <three-viewer
                            :model-url="getFileModelUrl(file)"
                            :model-type="getFileModelType(file)"
                            :width="400"
                            :height="300"
                            @model-loaded="(model) => onModelLoaded(file.id, model)"
                            @model-error="(error) => onModelError(file.id, error)"
                            @model-stats="(stats) => onModelStats(file.id, stats)"
                          />
                        </div>
                        
                        <div class="file-info" v-if="file.stats">
                          <el-descriptions :column="2" size="small" border>
                            <el-descriptions-item label="顶点数">{{ file.stats.vertices }}</el-descriptions-item>
                            <el-descriptions-item label="面数">{{ file.stats.faces }}</el-descriptions-item>
                            <el-descriptions-item label="体积">{{ file.stats.volume }} cm³</el-descriptions-item>
                            <el-descriptions-item label="文件大小">{{ formatFileSize(file.size) }}</el-descriptions-item>
                          </el-descriptions>
                        </div>
                      </div>
                    </el-col>
                    
                    <!-- Print Configuration -->
                    <el-col :xs="24" :sm="24" :md="12">
                      <div class="config-section">
                        <h3>打印配置</h3>
                        
                        <el-form :model="file.config" label-width="100px">
                          <el-form-item label="材料选择">
                            <el-select v-model="file.config.materialId" placeholder="请选择材料" style="width: 100%;">
                              <el-option
                                v-for="material in materials"
                                :key="material.id"
                                :label="`${material.name} - ¥${material.price_per_gram}/克`"
                                :value="material.id"
                              />
                            </el-select>
                          </el-form-item>
                          
                          <el-form-item label="打印质量">
                            <el-select v-model="file.config.quality" placeholder="请选择质量" style="width: 100%;">
                              <el-option label="快速 (0.3mm)" value="fast" />
                              <el-option label="标准 (0.2mm)" value="standard" />
                              <el-option label="精细 (0.1mm)" value="high" />
                            </el-select>
                          </el-form-item>
                          
                          <el-form-item label="支撑结构">
                            <el-switch v-model="file.config.supportEnabled" />
                          </el-form-item>
                          
                          <el-form-item label="数量">
                            <el-input-number v-model="file.config.quantity" :min="1" :max="10" style="width: 100%;" />
                          </el-form-item>
                          
                          <el-form-item>
                            <div class="price-display">
                              <span class="price-label">预估价格: </span>
                              <span class="price-value">¥{{ getFileEstimatedPrice(file) }}</span>
                            </div>
                          </el-form-item>
                          
                          <el-form-item>
                            <el-button 
                              type="primary" 
                              @click="addFileToCart(file)" 
                              :disabled="!canAddFileToCart(file)"
                              style="width: 100%;"
                            >
                              加入购物车
                            </el-button>
                          </el-form-item>
                        </el-form>
                      </div>
                    </el-col>
                  </el-row>
                </div>
              </el-tab-pane>
            </el-tabs>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- Empty State -->
    <el-row :gutter="20" style="margin-top: 20px;" v-if="uploadedFiles.length === 0">
      <el-col :span="24">
        <el-empty description="暂无上传文件" :image-size="200">
          <template #description>
            <p>请先上传3D模型文件或加载演示模型</p>
          </template>
        </el-empty>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { useUserStore } from '@/stores/user'
import { useModelStore } from '@/stores/model'
import { useCartStore } from '@/stores/cart'
import { ElMessage, ElMessageBox } from 'element-plus'
import { UploadFilled, View, Delete } from '@element-plus/icons-vue'
import axios from 'axios'
import ThreeViewer from '@/components/ThreeViewer.vue'
import { ModelGenerator } from '@/utils/modelGenerator.js'

export default {
  name: 'Upload',
  components: {
    ThreeViewer,
    UploadFilled,
    View,
    Delete
  },
  setup() {
    const userStore = useUserStore()
    const modelStore = useModelStore()
    const cartStore = useCartStore()
    
    const materials = ref([])
    const uploadedFiles = ref([])
    const activeFileTab = ref('')
    
    // Default configuration for new files
    const getDefaultConfig = () => ({
      materialId: null,
      quality: 'standard',
      quantity: 1,
      supportEnabled: false
    })
    
    const uploadHeaders = computed(() => ({
      Authorization: `Bearer ${userStore.token}`
    }))
    
    const canAddAllToCart = computed(() => {
      return uploadedFiles.value.length > 0 && uploadedFiles.value.every(file => 
        file.config.materialId && file.config.quality
      )
    })
    
    const getFileEstimatedPrice = (file) => {
      if (!file.config.materialId) return '0.00'
      
      const selectedMaterial = materials.value.find(m => m.id === file.config.materialId)
      if (!selectedMaterial) return '0.00'
      
      const basePrice = selectedMaterial.price_per_gram * 100
      const qualityMultiplier = {
        'fast': 1.0,
        'standard': 1.2,
        'high': 1.5
      }
      const supportMultiplier = file.config.supportEnabled ? 1.1 : 1.0
      const finalPrice = basePrice * qualityMultiplier[file.config.quality] * supportMultiplier
      
      return (finalPrice * file.config.quantity).toFixed(2)
    }
    
    const canAddFileToCart = (file) => {
      return file && file.config.materialId && file.config.quality
    }
    
    const beforeUpload = (file) => {
      const isValidType = file.type === 'application/octet-stream' || file.name.endsWith('.stl') || file.name.endsWith('.obj')
      const isLt50M = file.size / 1024 / 1024 < 50
      
      if (!isValidType) {
        ElMessage.error('只支持STL和OBJ格式的文件!')
        return false
      }
      if (!isLt50M) {
        ElMessage.error('文件大小不能超过50MB!')
        return false
      }
      return true
    }
    
    const onUploadSuccess = (response) => {
      const fileWithConfig = {
        ...response,
        config: getDefaultConfig(),
        stats: null
      }
      
      uploadedFiles.value.push(fileWithConfig)
      
      // Set active tab to the newly uploaded file
      activeFileTab.value = response.id
      
      modelStore.addUploadedModel(response)
      ElMessage.success(`文件 ${response.filename || response.original_name} 上传成功!`)
    }
    
    const onUploadError = () => {
      ElMessage.error('文件上传失败!')
    }
    
    const getFileModelUrl = (file) => {
      if (!file) return ''
      
      // 如果是演示模型，直接返回blob URL
      if (file.url && file.url.startsWith('blob:')) {
        return file.url
      }
      
      // 对于真实上传的文件，使用API端点
      return file.url || `http://localhost:8080/api/files/${file.id}`
    }
    
    const getFileModelType = (file) => {
      if (!file) return 'stl'
      const filename = file.filename || file.original_name || ''
      return filename.toLowerCase().endsWith('.obj') ? 'obj' : 'stl'
    }
    
    const getTabLabel = (file, index) => {
      const filename = file.filename || file.original_name || `模型${index + 1}`
      return filename.length > 15 ? filename.substring(0, 15) + '...' : filename
    }
    
    const formatFileSize = (bytes) => {
      if (bytes === 0) return '0 B'
      const k = 1024
      const sizes = ['B', 'KB', 'MB', 'GB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
    }
    
    const onModelLoaded = (fileId, model) => {
      console.log('3D模型加载成功:', fileId, model)
    }
    
    const onModelError = (fileId, error) => {
      console.error('3D模型加载失败:', fileId, error)
      ElMessage.error('3D模型预览失败')
    }
    
    const onModelStats = (fileId, stats) => {
      console.log('模型统计信息:', fileId, stats)
      const file = uploadedFiles.value.find(f => f.id === fileId)
      if (file) {
        file.stats = stats
      }
      modelStore.updateModelStats(stats)
    }
    
    const updatePriceCalculation = (volume) => {
      // 价格计算现在由modelStore处理
      // 这个方法可以用于触发界面更新
      console.log('模型体积更新:', volume)
    }
    
    const addFileToCart = (file) => {
      // 验证必要信息
      if (!file.config.materialId) {
        ElMessage.warning('请选择打印材料')
        return
      }
      
      if (!file.config.quality) {
        ElMessage.warning('请选择打印质量')
        return
      }

      // 获取材料信息
      const selectedMaterial = materials.value.find(m => m.id === file.config.materialId)
      if (!selectedMaterial) {
        ElMessage.error('所选材料不存在')
        return
      }

      // 计算价格
      const basePrice = selectedMaterial.price_per_gram * 100
      const qualityMultiplier = {
        'fast': 1.0,
        'standard': 1.2,
        'high': 1.5
      }
      const supportMultiplier = file.config.supportEnabled ? 1.1 : 1.0
      const finalPrice = basePrice * qualityMultiplier[file.config.quality] * supportMultiplier

      // 准备购物车商品数据
      const cartItem = {
        fileId: file.id,
        fileName: file.filename || file.original_name || '3D模型',
        materialId: file.config.materialId,
        materialName: selectedMaterial.name,
        quality: file.config.quality,
        quantity: file.config.quantity,
        supportEnabled: file.config.supportEnabled,
        price: finalPrice
      }

      // 添加到购物车
      const success = cartStore.addItem(cartItem)
      
      if (success) {
        console.log('商品已添加到购物车:', cartItem)
      }
    }
    
    const addAllToCart = async () => {
      const validFiles = uploadedFiles.value.filter(file => canAddFileToCart(file))
      
      if (validFiles.length === 0) {
        ElMessage.warning('没有可添加的文件，请先配置打印参数')
        return
      }
      
      try {
        await ElMessageBox.confirm(
          `确定将 ${validFiles.length} 个文件全部添加到购物车吗？`,
          '批量添加确认',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'info'
          }
        )
        
        let successCount = 0
        validFiles.forEach(file => {
          addFileToCart(file)
          successCount++
        })
        
        ElMessage.success(`成功添加 ${successCount} 个文件到购物车`)
      } catch {
        // 用户取消操作
      }
    }
    
    const loadDemoModel = () => {
      try {
        ElMessage.info('正在生成演示模型...')
        
        // 生成演示立方体模型
        const demoFile = ModelGenerator.createDemoFile('cube', 20)
        const demoUrl = ModelGenerator.createModelURL(demoFile)
        
        console.log('演示文件生成成功:', {
          file: demoFile,
          url: demoUrl,
          size: demoFile.size
        })
        
        // 模拟上传成功的响应
        const mockResponse = {
          id: 'demo-' + Date.now(),
          filename: 'demo_cube.stl',
          original_name: 'demo_cube.stl',
          size: demoFile.size,
          url: demoUrl,
          created_at: new Date().toISOString(),
          config: getDefaultConfig(),
          stats: null
        }
        
        uploadedFiles.value.push(mockResponse)
        activeFileTab.value = mockResponse.id
        modelStore.addUploadedModel(mockResponse)
        
        ElMessage.success('演示模型加载成功！')
      } catch (error) {
        console.error('加载演示模型失败:', error)
        ElMessage.error('加载演示模型失败: ' + error.message)
      }
    }
    
    const removeFile = async (fileId) => {
      try {
        await ElMessageBox.confirm(
          '确定要删除这个文件吗？',
          '删除确认',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        
        const fileIndex = uploadedFiles.value.findIndex(f => f.id === fileId)
        if (fileIndex > -1) {
          uploadedFiles.value.splice(fileIndex, 1)
          
          // 如果删除的是当前活动标签页，切换到第一个
          if (activeFileTab.value === fileId && uploadedFiles.value.length > 0) {
            activeFileTab.value = uploadedFiles.value[0].id
          }
          
          ElMessage.success('文件删除成功')
        }
      } catch {
        // 用户取消操作
      }
    }
    
    const clearAllFiles = async () => {
      try {
        await ElMessageBox.confirm(
          '确定要清空所有文件吗？',
          '清空确认',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        
        uploadedFiles.value = []
        activeFileTab.value = ''
        ElMessage.success('所有文件已清空')
      } catch {
        // 用户取消操作
      }
    }
    
    const previewOrder = () => {
      // TODO: 实现订单预览功能
      console.log('预览订单')
    }
    
    const loadMaterials = async () => {
      try {
        console.log('正在加载材料列表...')
        const response = await axios.get('/api/materials')
        console.log('材料数据:', response.data)
        materials.value = response.data
        console.log('materials.value:', materials.value)
      } catch (error) {
        console.error('加载材料失败:', error)
        ElMessage.error('加载材料列表失败: ' + (error.response?.data?.error || error.message))
      }
    }
    
    onMounted(() => {
      loadMaterials()
    })
    
    return {
      materials,
      uploadedFiles,
      activeFileTab,
      uploadHeaders,
      canAddAllToCart,
      beforeUpload,
      onUploadSuccess,
      onUploadError,
      getFileModelUrl,
      getFileModelType,
      getTabLabel,
      formatFileSize,
      getFileEstimatedPrice,
      canAddFileToCart,
      onModelLoaded,
      onModelError,
      onModelStats,
      loadDemoModel,
      addFileToCart,
      addAllToCart,
      removeFile,
      clearAllFiles,
      previewOrder
    }
  }
}
</script>

<style scoped>
.upload-page {
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px;
}

.upload-card, .file-list-card {
  margin-bottom: 20px;
}

.upload-header, .file-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions, .batch-actions {
  display: flex;
  gap: 10px;
}

.upload-dragger {
  width: 100%;
}

.file-tabs {
  margin-top: 10px;
}

.file-content {
  padding: 20px 0;
}

.preview-section, .config-section {
  height: 100%;
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.preview-header h3 {
  margin: 0;
  color: #333;
  font-size: 1.1rem;
  word-break: break-all;
}

.three-preview {
  height: 300px;
  background: #f5f5f5;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 15px;
}

.file-info {
  margin-top: 15px;
}

.config-section h3 {
  margin: 0 0 20px 0;
  color: #333;
  font-size: 1.1rem;
}

.price-display {
  font-size: 16px;
  margin: 15px 0;
  text-align: center;
  padding: 10px;
  background: #f8f9fa;
  border-radius: 4px;
}

.price-label {
  color: #666;
}

.price-value {
  color: #e6a23c;
  font-weight: bold;
  font-size: 18px;
}

/* Responsive Design */
@media (max-width: 768px) {
  .upload-page {
    padding: 15px;
  }
  
  .file-content {
    padding: 15px 0;
  }
  
  .three-preview {
    height: 250px;
  }
  
  .preview-header h3 {
    font-size: 1rem;
  }
  
  .config-section h3 {
    font-size: 1rem;
  }
}

@media (max-width: 480px) {
  .upload-page {
    padding: 10px;
  }
  
  .header-actions, .batch-actions {
    flex-direction: column;
    gap: 5px;
  }
  
  .three-preview {
    height: 200px;
  }
  
  .price-display {
    font-size: 14px;
  }
  
  .price-value {
    font-size: 16px;
  }
}
</style>