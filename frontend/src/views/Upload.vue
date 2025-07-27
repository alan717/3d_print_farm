<template>
  <div class="upload-page">
    <el-row :gutter="20">
      <el-col :span="24">
        <h1>上传3D模型</h1>
      </el-col>
    </el-row>
    
    <el-row :gutter="20">
      <el-col :span="12">
        <el-card class="upload-card">
          <template #header>
            <span>文件上传</span>
          </template>
          
          <el-upload
            class="upload-dragger"
            drag
            action="/api/files/upload"
            :headers="uploadHeaders"
            accept=".stl,.obj"
            :before-upload="beforeUpload"
            :on-success="onUploadSuccess"
            :on-error="onUploadError"
          >
            <el-icon class="el-icon--upload">
              <UploadFilled />
            </el-icon>
            <div class="el-upload__text">
              将文件拖到此处，或<em>点击上传</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                支持STL、OBJ格式，文件大小不超过50MB
              </div>
            </template>
          </el-upload>
          
          <div class="demo-section">
            <el-divider>或</el-divider>
            <el-button @click="loadDemoModel" type="success" plain>
              加载演示模型
            </el-button>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card class="preview-card">
          <template #header>
            <span>3D预览</span>
          </template>
          
          <div class="three-preview">
            <three-viewer
              v-if="uploadedFile"
              :model-url="getModelUrl()"
              :model-type="getModelType()"
              :width="500"
              :height="300"
              @model-loaded="onModelLoaded"
              @model-error="onModelError"
              @model-stats="onModelStats"
            />
            <div v-else class="preview-placeholder">
              <el-icon size="60">
                <View />
              </el-icon>
              <p>上传模型后将显示3D预览</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="24">
        <el-card class="config-card">
          <template #header>
            <span>打印配置</span>
          </template>
          
          <el-form :model="printConfig" label-width="120px">
            <el-row :gutter="20">
              <el-col :span="8">
                <el-form-item label="材料选择">
                  <el-select v-model="printConfig.materialId" placeholder="请选择材料">
                    <el-option
                      v-for="material in materials"
                      :key="material.id"
                      :label="`${material.name} - ¥${material.price_per_gram}/克`"
                      :value="material.id"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
              
              <el-col :span="8">
                <el-form-item label="打印质量">
                  <el-select v-model="printConfig.quality" placeholder="请选择质量">
                    <el-option label="快速 (0.3mm)" value="fast" />
                    <el-option label="标准 (0.2mm)" value="standard" />
                    <el-option label="精细 (0.1mm)" value="high" />
                  </el-select>
                </el-form-item>
              </el-col>
              
              <el-col :span="8">
                <el-form-item label="支撑结构">
                  <el-switch v-model="printConfig.supportEnabled" />
                </el-form-item>
              </el-col>
            </el-row>
            
            <el-form-item label="数量">
              <el-input-number v-model="printConfig.quantity" :min="1" :max="10" />
            </el-form-item>
            
            <el-form-item>
              <div class="price-display">
                <span class="price-label">预估价格: </span>
                <span class="price-value">¥{{ estimatedPrice }}</span>
              </div>
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="addToCart" :disabled="!canAddToCart">
                加入购物车
              </el-button>
              <el-button @click="previewOrder">预览订单</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { useModelStore } from '@/stores/model'
import { ElMessage } from 'element-plus'
import { UploadFilled, View } from '@element-plus/icons-vue'
import axios from 'axios'
import ThreeViewer from '@/components/ThreeViewer.vue'
import { ModelGenerator } from '@/utils/modelGenerator.js'

export default {
  name: 'Upload',
  components: {
    ThreeViewer,
    UploadFilled,
    View
  },
  setup() {
    const userStore = useUserStore()
    const modelStore = useModelStore()
    
    const materials = ref([])
    const uploadedFile = ref(null)
    
    const printConfig = reactive({
      materialId: null,
      quality: 'standard',
      quantity: 1,
      supportEnabled: false
    })
    
    const uploadHeaders = computed(() => ({
      Authorization: `Bearer ${userStore.token}`
    }))
    
    const estimatedPrice = computed(() => {
      if (!uploadedFile.value || !printConfig.materialId) return '0.00'
      
      const pricing = modelStore.calculatePrintSettings(
        printConfig.materialId,
        printConfig.quality,
        printConfig.supportEnabled
      )
      
      return (pricing.price * printConfig.quantity).toFixed(2)
    })
    
    const canAddToCart = computed(() => {
      return uploadedFile.value && printConfig.materialId && printConfig.quality
    })
    
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
      uploadedFile.value = response
      modelStore.addUploadedModel(response)
      ElMessage.success('文件上传成功!')
    }
    
    const onUploadError = () => {
      ElMessage.error('文件上传失败!')
    }
    
    const getModelUrl = () => {
      if (!uploadedFile.value) return ''
      
      // 如果是演示模型，直接返回blob URL
      if (uploadedFile.value.url && uploadedFile.value.url.startsWith('blob:')) {
        return uploadedFile.value.url
      }
      
      // 对于真实上传的文件，使用API端点
      return uploadedFile.value.url || `http://localhost:8080/api/files/${uploadedFile.value.id}`
    }
    
    const getModelType = () => {
      if (!uploadedFile.value) return 'stl'
      const filename = uploadedFile.value.filename || uploadedFile.value.original_name || ''
      return filename.toLowerCase().endsWith('.obj') ? 'obj' : 'stl'
    }
    
    const onModelLoaded = (model) => {
      console.log('3D模型加载成功:', model)
    }
    
    const onModelError = (error) => {
      console.error('3D模型加载失败:', error)
      ElMessage.error('3D模型预览失败')
    }
    
    const onModelStats = (stats) => {
      console.log('模型统计信息:', stats)
      modelStore.updateModelStats(stats)
    }
    
    const updatePriceCalculation = (volume) => {
      // 价格计算现在由modelStore处理
      // 这个方法可以用于触发界面更新
      console.log('模型体积更新:', volume)
    }
    
    const addToCart = () => {
      // TODO: 实现添加到购物车功能
      ElMessage.success('已添加到购物车!')
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
          created_at: new Date().toISOString()
        }
        
        uploadedFile.value = mockResponse
        modelStore.addUploadedModel(mockResponse)
        
        ElMessage.success('演示模型加载成功！')
      } catch (error) {
        console.error('加载演示模型失败:', error)
        ElMessage.error('加载演示模型失败: ' + error.message)
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
      uploadedFile,
      printConfig,
      uploadHeaders,
      estimatedPrice,
      canAddToCart,
      beforeUpload,
      onUploadSuccess,
      onUploadError,
      getModelUrl,
      getModelType,
      onModelLoaded,
      onModelError,
      onModelStats,
      loadDemoModel,
      addToCart,
      previewOrder
    }
  }
}
</script>

<style scoped>
.upload-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.upload-card, .preview-card, .config-card {
  height: auto;
  margin-bottom: 20px;
}

.upload-dragger {
  width: 100%;
}

.demo-section {
  margin-top: 20px;
  text-align: center;
}

.three-preview {
  height: 300px;
  background: #f5f5f5;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.preview-placeholder {
  text-align: center;
  color: #999;
}

.preview-placeholder p {
  margin: 10px 0 0 0;
}

.price-display {
  font-size: 18px;
  margin: 20px 0;
}

.price-label {
  color: #666;
}

.price-value {
  color: #e6a23c;
  font-weight: bold;
  font-size: 20px;
}
</style>