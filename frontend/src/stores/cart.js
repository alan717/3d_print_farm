import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'

export const useCartStore = defineStore('cart', () => {
  // 购物车商品列表
  const items = ref([])

  // 计算购物车总数量
  const totalItems = computed(() => {
    return items.value.reduce((total, item) => total + item.quantity, 0)
  })

  // 计算购物车总价格
  const totalPrice = computed(() => {
    return items.value.reduce((total, item) => total + (item.price * item.quantity), 0)
  })

  // 添加商品到购物车
  const addItem = (item) => {
    try {
      // 检查是否已存在相同配置的商品
      const existingItemIndex = items.value.findIndex(cartItem => 
        cartItem.fileId === item.fileId &&
        cartItem.materialId === item.materialId &&
        cartItem.quality === item.quality &&
        cartItem.supportEnabled === item.supportEnabled
      )

      if (existingItemIndex > -1) {
        // 如果已存在，增加数量
        items.value[existingItemIndex].quantity += item.quantity
        ElMessage.success(`已增加商品数量，当前数量：${items.value[existingItemIndex].quantity}`)
      } else {
        // 如果不存在，添加新商品
        const cartItem = {
          id: generateCartItemId(),
          fileId: item.fileId,
          fileName: item.fileName,
          materialId: item.materialId,
          materialName: item.materialName,
          quality: item.quality,
          quantity: item.quantity,
          supportEnabled: item.supportEnabled,
          price: item.price,
          addedAt: new Date().toISOString()
        }
        items.value.push(cartItem)
        ElMessage.success('商品已添加到购物车')
      }

      // 保存到本地存储
      saveToLocalStorage()
      return true
    } catch (error) {
      console.error('添加商品到购物车失败:', error)
      ElMessage.error('添加到购物车失败')
      return false
    }
  }

  // 更新商品数量
  const updateQuantity = (itemId, quantity) => {
    const item = items.value.find(item => item.id === itemId)
    if (item) {
      if (quantity <= 0) {
        removeItem(itemId)
      } else {
        item.quantity = quantity
        saveToLocalStorage()
        ElMessage.success('数量已更新')
      }
    }
  }

  // 移除商品
  const removeItem = (itemId) => {
    const index = items.value.findIndex(item => item.id === itemId)
    if (index > -1) {
      items.value.splice(index, 1)
      saveToLocalStorage()
      ElMessage.success('商品已移除')
    }
  }

  // 清空购物车
  const clearCart = () => {
    items.value = []
    saveToLocalStorage()
    ElMessage.success('购物车已清空')
  }

  // 生成购物车商品ID
  const generateCartItemId = () => {
    return 'cart_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9)
  }

  // 保存到本地存储
  const saveToLocalStorage = () => {
    try {
      localStorage.setItem('shopping_cart', JSON.stringify(items.value))
    } catch (error) {
      console.error('保存购物车到本地存储失败:', error)
    }
  }

  // 从本地存储加载
  const loadFromLocalStorage = () => {
    try {
      const saved = localStorage.getItem('shopping_cart')
      if (saved) {
        items.value = JSON.parse(saved)
      }
    } catch (error) {
      console.error('从本地存储加载购物车失败:', error)
      items.value = []
    }
  }

  // 获取质量显示名称
  const getQualityDisplayName = (quality) => {
    const qualityMap = {
      'fast': '快速 (0.3mm)',
      'standard': '标准 (0.2mm)', 
      'high': '精细 (0.1mm)'
    }
    return qualityMap[quality] || quality
  }

  // 初始化时从本地存储加载
  loadFromLocalStorage()

  return {
    items,
    totalItems,
    totalPrice,
    addItem,
    updateQuantity,
    removeItem,
    clearCart,
    getQualityDisplayName,
    loadFromLocalStorage
  }
})