<template>
  <div class="nav-header">
    <el-row justify="space-between" align="middle" style="height: 100%;">
      <el-col :span="6">
        <router-link to="/" class="logo">
          <h2>3D打印农场</h2>
        </router-link>
      </el-col>
      <el-col :span="12">
        <el-menu
          mode="horizontal"
          :default-active="$route.path"
          router
          class="nav-menu"
        >
          <el-menu-item index="/">首页</el-menu-item>
          <el-menu-item index="/upload">上传模型</el-menu-item>
          <el-menu-item index="/orders">我的订单</el-menu-item>
        </el-menu>
      </el-col>
      <el-col :span="6" style="text-align: right;">
        <div v-if="isLoggedIn" class="user-actions">
          <el-badge :value="cartStore.totalItems" :hidden="cartStore.totalItems === 0" class="cart-badge">
            <el-button @click="$router.push('/cart')" text>
              <el-icon size="20"><ShoppingCart /></el-icon>
            </el-button>
          </el-badge>
          <el-dropdown>
            <span class="el-dropdown-link">
              {{ userInfo.username }}
              <el-icon><arrow-down /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="$router.push('/profile')">个人中心</el-dropdown-item>
                <el-dropdown-item @click="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
        <div v-else class="auth-buttons">
          <el-badge :value="cartStore.totalItems" :hidden="cartStore.totalItems === 0" class="cart-badge">
            <el-button @click="$router.push('/cart')" text>
              <el-icon size="20"><ShoppingCart /></el-icon>
            </el-button>
          </el-badge>
          <el-button @click="$router.push('/login')">登录</el-button>
          <el-button type="primary" @click="$router.push('/register')">注册</el-button>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useCartStore } from '@/stores/cart'
import { ShoppingCart, ArrowDown } from '@element-plus/icons-vue'

export default {
  name: 'NavHeader',
  components: {
    ShoppingCart,
    ArrowDown
  },
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    const cartStore = useCartStore()
    
    const isLoggedIn = computed(() => userStore.isLoggedIn)
    const userInfo = computed(() => userStore.userInfo)
    
    const logout = () => {
      userStore.logout()
      router.push('/')
    }
    
    return {
      isLoggedIn,
      userInfo,
      cartStore,
      logout
    }
  }
}
</script>

<style scoped>
.nav-header {
  height: 60px;
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  padding: 0 20px;
}

.logo {
  text-decoration: none;
  color: #409eff;
}

.logo h2 {
  margin: 0;
  font-weight: 600;
}

.nav-menu {
  border-bottom: none;
  justify-content: center;
}

.user-actions {
  padding: 0 20px;
}

.el-dropdown-link {
  cursor: pointer;
  color: #409eff;
  display: flex;
  align-items: center;
}

.auth-buttons {
  padding: 0 20px;
}

.auth-buttons .el-button {
  margin-left: 10px;
}

.cart-badge {
  margin-right: 15px;
}

.cart-badge .el-button {
  color: #666;
  padding: 8px;
}

.cart-badge .el-button:hover {
  color: #409eff;
}
</style>