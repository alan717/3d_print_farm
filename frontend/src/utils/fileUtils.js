import axios from 'axios'

// 文件上传工具类
export class FileUploadUtil {
  static async uploadModel(file, token) {
    const formData = new FormData()
    formData.append('file', file)
    
    try {
      const response = await axios.post('/api/files/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
          'Authorization': `Bearer ${token}`
        },
        onUploadProgress: (progressEvent) => {
          const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
          console.log(`上传进度: ${progress}%`)
        }
      })
      
      return {
        success: true,
        data: response.data
      }
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.message || '上传失败'
      }
    }
  }
  
  static validateModelFile(file) {
    const allowedTypes = ['.stl', '.obj']
    const maxSize = 50 * 1024 * 1024 // 50MB
    
    // 检查文件扩展名
    const fileName = file.name.toLowerCase()
    const isValidType = allowedTypes.some(type => fileName.endsWith(type))
    
    if (!isValidType) {
      return {
        valid: false,
        message: '只支持STL和OBJ格式的文件'
      }
    }
    
    // 检查文件大小
    if (file.size > maxSize) {
      return {
        valid: false,
        message: '文件大小不能超过50MB'
      }
    }
    
    return {
      valid: true,
      message: '文件验证通过'
    }
  }
  
  static getFileInfo(file) {
    return {
      name: file.name,
      size: file.size,
      type: file.type,
      lastModified: file.lastModified,
      extension: file.name.split('.').pop().toLowerCase()
    }
  }
}

// 3D模型分析工具
export class ModelAnalyzer {
  static calculateVolume(geometry) {
    // 简化的体积计算方法
    // 实际应用中可能需要更复杂的算法
    if (!geometry.boundingBox) {
      geometry.computeBoundingBox()
    }
    
    const box = geometry.boundingBox
    const size = {
      x: box.max.x - box.min.x,
      y: box.max.y - box.min.y,
      z: box.max.z - box.min.z
    }
    
    // 返回边界盒体积（实际体积会小于这个值）
    return size.x * size.y * size.z
  }
  
  static getModelComplexity(geometry) {
    const vertexCount = geometry.attributes.position.count
    const faceCount = vertexCount / 3
    
    if (vertexCount < 1000) {
      return { level: 'low', description: '简单模型' }
    } else if (vertexCount < 10000) {
      return { level: 'medium', description: '中等复杂度' }
    } else {
      return { level: 'high', description: '高复杂度模型' }
    }
  }
  
  static estimatePrintTime(volume, quality) {
    // 简化的打印时间估算
    const baseTimePerCm3 = {
      fast: 10,     // 分钟/cm³
      standard: 15,
      high: 25
    }
    
    const timeMinutes = volume * (baseTimePerCm3[quality] || 15)
    const hours = Math.floor(timeMinutes / 60)
    const minutes = Math.round(timeMinutes % 60)
    
    return {
      totalMinutes: timeMinutes,
      formatted: `${hours}小时${minutes}分钟`
    }
  }
}