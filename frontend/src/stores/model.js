import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { FileUploadUtil, ModelAnalyzer } from '@/utils/fileUtils'

export const useModelStore = defineStore('model', () => {
  const uploadedModels = ref([])
  const currentModel = ref(null)
  const modelStats = ref({
    volume: 0,
    vertices: 0,
    faces: 0,
    complexity: null,
    estimatedPrintTime: null
  })

  const hasModel = computed(() => !!currentModel.value)
  const modelCount = computed(() => uploadedModels.value.length)

  const setCurrentModel = (model) => {
    currentModel.value = model
  }

  const updateModelStats = (stats) => {
    modelStats.value = { ...modelStats.value, ...stats }
  }

  const addUploadedModel = (modelData) => {
    const model = {
      id: modelData.id || Date.now().toString(),
      filename: modelData.filename,
      originalName: modelData.original_name,
      size: modelData.size,
      url: modelData.url,
      uploadTime: new Date().toISOString(),
      stats: modelData.stats || {}
    }
    
    uploadedModels.value.push(model)
    setCurrentModel(model)
    
    return model
  }

  const removeModel = (modelId) => {
    const index = uploadedModels.value.findIndex(m => m.id === modelId)
    if (index > -1) {
      uploadedModels.value.splice(index, 1)
      
      // 如果删除的是当前模型，清空当前模型
      if (currentModel.value?.id === modelId) {
        currentModel.value = null
        modelStats.value = {
          volume: 0,
          vertices: 0,
          faces: 0,
          complexity: null,
          estimatedPrintTime: null
        }
      }
    }
  }

  const calculatePrintSettings = (materialId, quality, supportEnabled) => {
    if (!modelStats.value.volume || !materialId) {
      return {
        price: 0,
        weight: 0,
        printTime: '未知'
      }
    }

    // 材料密度映射（g/cm³）
    const materialDensity = {
      1: 1.24, // PLA
      2: 1.04, // ABS  
      3: 1.27, // PETG
      4: 1.20, // TPU
      5: 1.28  // WOOD
    }

    // 材料价格映射（元/克）
    const materialPrice = {
      1: 0.8,  // PLA
      2: 1.0,  // ABS
      3: 1.2,  // PETG
      4: 2.0,  // TPU
      5: 1.5   // WOOD
    }

    const density = materialDensity[materialId] || 1.2
    const pricePerGram = materialPrice[materialId] || 1.0
    
    const weight = modelStats.value.volume * density
    const materialCost = weight * pricePerGram
    
    // 质量系数
    const qualityMultiplier = {
      fast: 0.8,
      standard: 1.0,
      high: 1.5
    }
    
    const baseCost = materialCost * (qualityMultiplier[quality] || 1.0)
    const supportCost = supportEnabled ? Math.max(5, weight * 0.1) : 0
    const totalPrice = baseCost + supportCost

    // 估算打印时间
    const printTime = ModelAnalyzer.estimatePrintTime(modelStats.value.volume, quality)

    return {
      price: parseFloat(totalPrice.toFixed(2)),
      weight: parseFloat(weight.toFixed(2)),
      printTime: printTime.formatted,
      breakdown: {
        materialCost: parseFloat(materialCost.toFixed(2)),
        qualityMultiplier: qualityMultiplier[quality] || 1.0,
        supportCost: parseFloat(supportCost.toFixed(2)),
        totalCost: parseFloat(totalPrice.toFixed(2))
      }
    }
  }

  const clearModels = () => {
    uploadedModels.value = []
    currentModel.value = null
    modelStats.value = {
      volume: 0,
      vertices: 0,
      faces: 0,
      complexity: null,
      estimatedPrintTime: null
    }
  }

  return {
    uploadedModels,
    currentModel,
    modelStats,
    hasModel,
    modelCount,
    setCurrentModel,
    updateModelStats,
    addUploadedModel,
    removeModel,
    calculatePrintSettings,
    clearModels
  }
})