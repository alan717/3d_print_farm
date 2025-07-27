<template>
  <div class="nav-header">
    <!-- Desktop Layout -->
    <div class="nav-desktop">
      <div class="nav-left">
        <router-link to="/" class="logo">
          <h2>3D打印农场</h2>
        </router-link>
      </div>
      
      <div class="nav-center">
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
      </div>
      
      <div class="nav-right">
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
      </div>
    </div>
    
    <!-- Mobile Layout -->
    <div class="nav-mobile">
      <div class="mobile-header">
        <router-link to="/" class="logo">
          <h2>3D打印农场</h2>
        </router-link>
        
        <div class="mobile-actions">
          <!-- Shopping Cart -->
          <el-badge :value="cartStore.totalItems" :hidden="cartStore.totalItems === 0" class="cart-badge mobile-cart">
            <el-button @click="$router.push('/cart')" text>
              <el-icon size="18"><ShoppingCart /></el-icon>
            </el-button>
          </el-badge>
          
          <!-- Hamburger Menu -->
          <el-button @click="toggleMobileMenu" text class="hamburger-btn">
            <el-icon size="20"><Menu /></el-icon>
          </el-button>
        </div>
      </div>
      
      <!-- Mobile Menu Drawer -->
      <transition name="slide-down">
        <div v-show="showMobileMenu" class="mobile-menu">
          <div class="menu-items">
            <router-link @click="closeMobileMenu" to="/" class="menu-item">
              <el-icon><HomeFilled /></el-icon>
              <span>首页</span>
            </router-link>
            <router-link @click="closeMobileMenu" to="/upload" class="menu-item">
              <el-icon><Upload /></el-icon>
              <span>上传模型</span>
            </router-link>
            <router-link @click="closeMobileMenu" to="/orders" class="menu-item">
              <el-icon><Document /></el-icon>
              <span>我的订单</span>
            </router-link>
          </div>
          
          <div class="menu-user-section">
            <div v-if="isLoggedIn" class="mobile-user-actions">
              <div class="user-info">
                <el-icon><User /></el-icon>
                <span>{{ userInfo.username }}</span>
              </div>
              <div class="user-buttons">
                <el-button @click="goToProfile" plain>个人中心</el-button>
                <el-button @click="logout" type="danger" plain>退出登录</el-button>
              </div>
            </div>
            <div v-else class="mobile-auth-buttons">
              <el-button @click="goToLogin" style="width: 100%; margin-bottom: 10px;">登录</el-button>
              <el-button @click="goToRegister" type="primary" style="width: 100%;">注册</el-button>
            </div>
          </div>
        </div>
      </transition>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useCartStore } from '@/stores/cart'
import { ShoppingCart, ArrowDown, Menu, HomeFilled, Upload, Document, User } from '@element-plus/icons-vue'

export default {
  name: 'NavHeader',
  components: {
    ShoppingCart,
    ArrowDown,
    Menu,
    HomeFilled,
    Upload,
    Document,
    User
  },
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    const cartStore = useCartStore()
    
    const showMobileMenu = ref(false)
    
    const isLoggedIn = computed(() => userStore.isLoggedIn)
    const userInfo = computed(() => userStore.userInfo)
    
    const toggleMobileMenu = () => {
      showMobileMenu.value = !showMobileMenu.value
    }
    
    const closeMobileMenu = () => {
      showMobileMenu.value = false
    }
    
    const logout = () => {
      userStore.logout()
      router.push('/')
      closeMobileMenu()
    }
    
    const goToProfile = () => {
      router.push('/profile')
      closeMobileMenu()
    }
    
    const goToLogin = () => {
      router.push('/login')
      closeMobileMenu()
    }
    
    const goToRegister = () => {
      router.push('/register')
      closeMobileMenu()
    }
    
    return {
      isLoggedIn,
      userInfo,
      cartStore,
      showMobileMenu,
      toggleMobileMenu,
      closeMobileMenu,
      logout,
      goToProfile,
      goToLogin,
      goToRegister
    }
  }
}
</script>

<style scoped>
.nav-header {
  height: 60px;
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  position: relative;
  z-index: 1000;
}

/* Desktop Layout */
.nav-desktop {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.nav-left {
  flex: 0 0 auto;
}

.nav-center {
  flex: 1 1 auto;
  display: flex;
  justify-content: center;
  max-width: 600px;
}

.nav-right {
  flex: 0 0 auto;
}

.logo {
  text-decoration: none;
  color: #409eff;
}

.logo h2 {
  margin: 0;
  font-weight: 600;
  font-size: 1.4rem;
  white-space: nowrap;
}

.nav-menu {
  border-bottom: none;
  width: 100%;
  justify-content: center;
}

.user-actions {
  display: flex;
  align-items: center;
  gap: 15px;
}

.el-dropdown-link {
  cursor: pointer;
  color: #409eff;
  display: flex;
  align-items: center;
  white-space: nowrap;
}

.auth-buttons {
  display: flex;
  align-items: center;
  gap: 10px;
}

.cart-badge {
  display: flex;
  align-items: center;
}

.cart-badge .el-button {
  color: #666;
  padding: 8px;
}

.cart-badge .el-button:hover {
  color: #409eff;
}

/* Mobile Layout */
.nav-mobile {
  display: none;
  height: 100%;
}

.mobile-header {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 15px;
}

.mobile-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.mobile-cart .el-button {
  padding: 6px;
  color: #666;
}

.hamburger-btn {
  padding: 6px;
  color: #666;
}

.mobile-menu {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: #fff;
  border-bottom: 1px solid #e6e6e6;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  z-index: 999;
}

.menu-items {
  padding: 15px 0;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 12px 20px;
  color: #333;
  text-decoration: none;
  border-bottom: 1px solid #f5f5f5;
  transition: background-color 0.3s;
}

.menu-item:hover {
  background-color: #f8f9fa;
}

.menu-item:last-child {
  border-bottom: none;
}

.menu-item .el-icon {
  margin-right: 10px;
  color: #666;
}

.menu-item span {
  font-size: 14px;
}

.menu-user-section {
  border-top: 1px solid #e6e6e6;
  padding: 15px 20px;
  background-color: #f8f9fa;
}

.mobile-user-actions .user-info {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
  font-weight: 500;
  color: #333;
}

.mobile-user-actions .user-info .el-icon {
  margin-right: 8px;
  color: #409eff;
}

.user-buttons {
  display: flex;
  gap: 10px;
}

.user-buttons .el-button {
  flex: 1;
}

.mobile-auth-buttons {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

/* Transitions */
.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.3s ease;
  transform-origin: top;
}

.slide-down-enter-from {
  opacity: 0;
  transform: translateY(-10px);
}

.slide-down-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* Responsive Breakpoints */
@media (min-width: 769px) {
  .nav-desktop {
    display: flex;
  }
  
  .nav-mobile {
    display: none;
  }
}

@media (max-width: 768px) {
  .nav-desktop {
    display: none;
  }
  
  .nav-mobile {
    display: block;
  }
  
  .logo h2 {
    font-size: 1.2rem;
  }
}

@media (max-width: 480px) {
  .mobile-header {
    padding: 0 10px;
  }
  
  .logo h2 {
    font-size: 1.1rem;
  }
  
  .menu-item {
    padding: 10px 15px;
  }
  
  .menu-user-section {
    padding: 12px 15px;
  }
}

@media (max-width: 360px) {
  .mobile-header {
    padding: 0 8px;
  }
  
  .logo h2 {
    font-size: 1rem;
  }
  
  .mobile-actions {
    gap: 5px;
  }
}

/* Desktop fine-tuning for different screen sizes */
@media (min-width: 769px) and (max-width: 1024px) {
  .nav-desktop {
    padding: 0 15px;
  }
  
  .nav-center {
    max-width: 500px;
  }
  
  .logo h2 {
    font-size: 1.3rem;
  }
}

@media (min-width: 1200px) {
  .nav-desktop {
    padding: 0 30px;
  }
  
  .nav-center {
    max-width: 700px;
  }
  
  .logo h2 {
    font-size: 1.5rem;
  }
}
</style>