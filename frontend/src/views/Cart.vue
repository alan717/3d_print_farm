<template>
  <div class="cart-page">
    <div class="cart-container">
      <el-row>
        <el-col :span="24">
          <div class="page-header">
            <h1>购物车</h1>
            <p v-if="cartStore.totalItems > 0">共 {{ cartStore.totalItems }} 件商品</p>
            <p v-else>购物车为空</p>
          </div>
        </el-col>
      </el-row>

      <div v-if="cartStore.items.length === 0" class="empty-cart">
        <el-empty description="购物车是空的">
          <el-button type="primary" @click="$router.push('/upload')">
            去上传模型
          </el-button>
        </el-empty>
      </div>

      <div v-else>
        <!-- 购物车商品列表 -->
        <el-card class="cart-items">
          <template #header>
            <div class="card-header">
              <span>商品列表</span>
              <el-button type="danger" size="small" @click="handleClearCart">
                清空购物车
              </el-button>
            </div>
          </template>

          <div class="cart-item" v-for="item in cartStore.items" :key="item.id">
            <el-row :gutter="20" align="middle">
              <el-col :span="2">
                <el-checkbox v-model="selectedItems" :label="item.id" />
              </el-col>
              
              <el-col :span="6">
                <div class="item-info">
                  <h4>{{ item.fileName }}</h4>
                  <p class="item-details">
                    材料: {{ item.materialName }}
                  </p>
                  <p class="item-details">
                    质量: {{ cartStore.getQualityDisplayName(item.quality) }}
                  </p>
                  <p class="item-details">
                    支撑: {{ item.supportEnabled ? '启用' : '禁用' }}
                  </p>
                </div>
              </el-col>

              <el-col :span="4">
                <div class="item-price">
                  <span class="price">¥{{ item.price.toFixed(2) }}</span>
                  <span class="unit">/件</span>
                </div>
              </el-col>

              <el-col :span="4">
                <el-input-number
                  v-model="item.quantity"
                  :min="1"
                  :max="10"
                  @change="(value) => updateQuantity(item.id, value)"
                  size="small"
                />
              </el-col>

              <el-col :span="4">
                <div class="item-total">
                  ¥{{ (item.price * item.quantity).toFixed(2) }}
                </div>
              </el-col>

              <el-col :span="4">
                <div class="item-actions">
                  <el-button
                    type="danger"
                    size="small"
                    @click="removeItem(item.id)"
                    text
                  >
                    删除
                  </el-button>
                </div>
              </el-col>
            </el-row>
          </div>
        </el-card>

        <!-- 购物车统计 -->
        <el-card class="cart-summary">
          <template #header>
            <span>订单统计</span>
          </template>

          <el-row>
            <el-col :span="12">
              <div class="summary-item">
                <span>商品总数:</span>
                <span class="value">{{ cartStore.totalItems }} 件</span>
              </div>
              <div class="summary-item">
                <span>商品金额:</span>
                <span class="value">¥{{ selectedTotalPrice.toFixed(2) }}</span>
              </div>
              <div class="summary-item">
                <span>运费:</span>
                <span class="value">¥{{ shippingCost.toFixed(2) }}</span>
              </div>
              <div class="summary-item total">
                <span>总计:</span>
                <span class="value">¥{{ (selectedTotalPrice + shippingCost).toFixed(2) }}</span>
              </div>
            </el-col>

            <el-col :span="12">
              <div class="checkout-section">
                <el-button
                  type="primary"
                  size="large"
                  :disabled="selectedItems.length === 0"
                  @click="handleCheckout"
                  class="checkout-button"
                >
                  结算 ({{ selectedItems.length }})
                </el-button>
                <p class="checkout-note">
                  请选择要结算的商品
                </p>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '@/stores/cart'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'

export default {
  name: 'Cart',
  setup() {
    const router = useRouter()
    const cartStore = useCartStore()
    const userStore = useUserStore()

    const selectedItems = ref([])
    const shippingCost = ref(15.00) // 固定运费，后续可以做成动态计算

    // 计算选中商品的总价
    const selectedTotalPrice = computed(() => {
      return cartStore.items
        .filter(item => selectedItems.value.includes(item.id))
        .reduce((total, item) => total + (item.price * item.quantity), 0)
    })

    // 更新商品数量
    const updateQuantity = (itemId, quantity) => {
      cartStore.updateQuantity(itemId, quantity)
    }

    // 移除商品
    const removeItem = async (itemId) => {
      try {
        await ElMessageBox.confirm(
          '确定要从购物车中移除这个商品吗？',
          '确认删除',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        cartStore.removeItem(itemId)
        // 如果移除的商品在选中列表中，也要移除
        const index = selectedItems.value.indexOf(itemId)
        if (index > -1) {
          selectedItems.value.splice(index, 1)
        }
      } catch {
        // 用户取消删除
      }
    }

    // 清空购物车
    const handleClearCart = async () => {
      try {
        await ElMessageBox.confirm(
          '确定要清空购物车吗？这个操作不能撤销。',
          '确认清空',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        cartStore.clearCart()
        selectedItems.value = []
      } catch {
        // 用户取消清空
      }
    }

    // 结算
    const handleCheckout = async () => {
      if (!userStore.isLoggedIn) {
        ElMessage.warning('请先登录')
        router.push('/login')
        return
      }

      if (selectedItems.value.length === 0) {
        ElMessage.warning('请选择要结算的商品')
        return
      }

      try {
        // 准备订单数据
        const orderItems = cartStore.items
          .filter(item => selectedItems.value.includes(item.id))
          .map(item => ({
            file_id: item.fileId,
            material_id: item.materialId,
            quality: item.quality,
            quantity: item.quantity,
            support_enabled: item.supportEnabled,
            price: item.price
          }))

        const orderData = {
          items: orderItems,
          total_amount: selectedTotalPrice.value + shippingCost.value,
          shipping_cost: shippingCost.value
        }

        console.log('订单数据:', orderData)
        
        // 这里应该调用创建订单的API
        // const response = await axios.post('/api/orders', orderData)
        
        // 模拟订单创建成功
        ElMessage.success('订单创建成功！')
        
        // 从购物车中移除已结算的商品
        selectedItems.value.forEach(itemId => {
          cartStore.removeItem(itemId)
        })
        selectedItems.value = []
        
        // 跳转到订单页面
        router.push('/orders')
        
      } catch (error) {
        console.error('创建订单失败:', error)
        ElMessage.error('创建订单失败，请重试')
      }
    }

    // 全选/取消全选
    const handleSelectAll = (checked) => {
      if (checked) {
        selectedItems.value = cartStore.items.map(item => item.id)
      } else {
        selectedItems.value = []
      }
    }

    onMounted(() => {
      // 默认选中所有商品
      selectedItems.value = cartStore.items.map(item => item.id)
    })

    return {
      cartStore,
      selectedItems,
      selectedTotalPrice,
      shippingCost,
      updateQuantity,
      removeItem,
      handleClearCart,
      handleCheckout,
      handleSelectAll
    }
  }
}
</script>

<style scoped>
.cart-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.page-header {
  text-align: center;
  margin-bottom: 30px;
}

.page-header h1 {
  margin: 0 0 10px 0;
  color: #333;
}

.page-header p {
  margin: 0;
  color: #666;
  font-size: 14px;
}

.empty-cart {
  text-align: center;
  padding: 60px 0;
}

.cart-items {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.cart-item {
  padding: 20px 0;
  border-bottom: 1px solid #f0f0f0;
}

.cart-item:last-child {
  border-bottom: none;
}

.item-info h4 {
  margin: 0 0 8px 0;
  color: #333;
  font-size: 16px;
}

.item-details {
  margin: 4px 0;
  color: #666;
  font-size: 14px;
}

.item-price {
  text-align: center;
}

.price {
  font-size: 18px;
  font-weight: 600;
  color: #e6a23c;
}

.unit {
  font-size: 12px;
  color: #999;
  margin-left: 4px;
}

.item-total {
  text-align: center;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.item-actions {
  text-align: center;
}

.cart-summary {
  background: #f8f9fa;
}

.summary-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  font-size: 16px;
}

.summary-item.total {
  border-top: 1px solid #ddd;
  margin-top: 10px;
  padding-top: 15px;
  font-size: 18px;
  font-weight: 600;
}

.summary-item .value {
  color: #e6a23c;
  font-weight: 600;
}

.checkout-section {
  text-align: center;
  padding: 20px;
}

.checkout-button {
  width: 200px;
  height: 50px;
  font-size: 18px;
  font-weight: 600;
}

.checkout-note {
  margin: 10px 0 0 0;
  color: #999;
  font-size: 14px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .cart-page {
    padding: 10px;
  }

  .cart-item {
    padding: 15px 0;
  }

  .item-info h4 {
    font-size: 14px;
  }

  .item-details {
    font-size: 12px;
  }

  .price {
    font-size: 16px;
  }

  .item-total {
    font-size: 16px;
  }

  .summary-item {
    font-size: 14px;
  }

  .checkout-button {
    width: 100%;
    height: 45px;
    font-size: 16px;
  }
}

@media (max-width: 480px) {
  .cart-item {
    padding: 10px 0;
  }

  .summary-item {
    font-size: 13px;
  }

  .checkout-button {
    height: 40px;
    font-size: 15px;
  }
}
</style>