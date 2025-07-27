<template>
  <div class="admin-layout">
    <!-- Admin Sidebar -->
    <div class="admin-sidebar" :class="{ 'sidebar-collapsed': collapsed }">
      <div class="sidebar-header">
        <div class="logo">
          <h3 v-if="!collapsed">管理后台</h3>
          <h3 v-else>Admin</h3>
        </div>
        <el-button @click="toggleSidebar" text class="collapse-btn">
          <el-icon><Fold v-if="!collapsed" /><Expand v-else /></el-icon>
        </el-button>
      </div>
      
      <el-menu
        :default-active="$route.path"
        :collapse="collapsed"
        :unique-opened="true"
        router
        class="admin-menu"
      >
        <el-menu-item index="/admin/dashboard">
          <el-icon><Monitor /></el-icon>
          <template #title>仪表板</template>
        </el-menu-item>
        
        <el-sub-menu index="users">
          <template #title>
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </template>
          <el-menu-item index="/admin/users">用户列表</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="orders">
          <template #title>
            <el-icon><Document /></el-icon>
            <span>订单管理</span>
          </template>
          <el-menu-item index="/admin/orders">订单列表</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="materials">
          <template #title>
            <el-icon><Box /></el-icon>
            <span>材料管理</span>
          </template>
          <el-menu-item index="/admin/materials">材料列表</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="system">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>系统管理</span>
          </template>
          <el-menu-item index="/admin/settings">系统设置</el-menu-item>
          <el-menu-item index="/admin/logs">操作日志</el-menu-item>
        </el-sub-menu>
        
        <el-menu-item index="/" class="back-to-site">
          <el-icon><Back /></el-icon>
          <template #title>返回网站</template>
        </el-menu-item>
      </el-menu>
    </div>
    
    <!-- Admin Content -->
    <div class="admin-content" :class="{ 'content-expanded': collapsed }">
      <!-- Top Bar -->
      <div class="admin-topbar">
        <div class="topbar-left">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/admin/dashboard' }">管理后台</el-breadcrumb-item>
            <el-breadcrumb-item v-if="breadcrumbs.length > 0" v-for="breadcrumb in breadcrumbs" :key="breadcrumb">
              {{ breadcrumb }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        
        <div class="topbar-right">
          <el-dropdown>
            <span class="user-dropdown">
              <el-avatar :size="32" :src="userAvatar">
                {{ userStore.userInfo.username?.charAt(0).toUpperCase() }}
              </el-avatar>
              <span class="username">{{ userStore.userInfo.username }}</span>
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="$router.push('/profile')">个人中心</el-dropdown-item>
                <el-dropdown-item @click="logout" divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
      
      <!-- Page Content -->
      <div class="page-content">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { 
  Monitor, User, Document, Box, Setting, Back, 
  Fold, Expand, ArrowDown 
} from '@element-plus/icons-vue'

export default {
  name: 'AdminLayout',
  components: {
    Monitor, User, Document, Box, Setting, Back,
    Fold, Expand, ArrowDown
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const userStore = useUserStore()
    
    const collapsed = ref(false)
    const userAvatar = ref('')
    
    const breadcrumbs = computed(() => {
      const routeMap = {
        '/admin/dashboard': [],
        '/admin/users': ['用户管理'],
        '/admin/orders': ['订单管理'],
        '/admin/materials': ['材料管理'],
        '/admin/settings': ['系统管理', '系统设置'],
        '/admin/logs': ['系统管理', '操作日志']
      }
      return routeMap[route.path] || []
    })
    
    const toggleSidebar = () => {
      collapsed.value = !collapsed.value
      // Save to localStorage
      localStorage.setItem('admin_sidebar_collapsed', collapsed.value.toString())
    }
    
    const logout = () => {
      userStore.logout()
      router.push('/admin/login')
    }
    
    // Load sidebar state from localStorage
    const savedCollapsed = localStorage.getItem('admin_sidebar_collapsed')
    if (savedCollapsed !== null) {
      collapsed.value = savedCollapsed === 'true'
    }
    
    return {
      collapsed,
      breadcrumbs,
      userAvatar,
      userStore,
      toggleSidebar,
      logout
    }
  }
}
</script>

<style scoped>
.admin-layout {
  display: flex;
  height: 100vh;
  background: #f5f7fa;
}

.admin-sidebar {
  width: 250px;
  background: #fff;
  border-right: 1px solid #e6e6e6;
  transition: width 0.3s;
  overflow: hidden;
}

.admin-sidebar.sidebar-collapsed {
  width: 64px;
}

.sidebar-header {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  border-bottom: 1px solid #e6e6e6;
}

.sidebar-collapsed .sidebar-header {
  padding: 0 16px;
  justify-content: center;
}

.logo h3 {
  margin: 0;
  color: #409eff;
  font-weight: 600;
}

.collapse-btn {
  color: #666;
}

.sidebar-collapsed .collapse-btn {
  position: absolute;
  right: 8px;
}

.admin-menu {
  border: none;
  height: calc(100vh - 60px);
}

.admin-menu .el-menu-item,
.admin-menu .el-sub-menu__title {
  height: 50px;
  line-height: 50px;
}

.back-to-site {
  position: absolute;
  bottom: 0;
  width: 100%;
  border-top: 1px solid #e6e6e6;
}

.admin-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  width: calc(100vw - 250px);
  transition: width 0.3s;
}

.admin-content.content-expanded {
  width: calc(100vw - 64px);
}

.admin-topbar {
  height: 60px;
  background: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
}

.topbar-left {
  flex: 1;
}

.topbar-right {
  display: flex;
  align-items: center;
}

.user-dropdown {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 6px;
  transition: background-color 0.2s;
}

.user-dropdown:hover {
  background-color: #f5f7fa;
}

.username {
  margin: 0 8px;
  color: #333;
  font-weight: 500;
}

.page-content {
  flex: 1;
  overflow-y: auto;
  background: #f5f7fa;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .admin-sidebar {
    position: fixed;
    left: 0;
    top: 0;
    bottom: 0;
    z-index: 1000;
    transform: translateX(-100%);
    transition: transform 0.3s;
  }
  
  .admin-sidebar.show {
    transform: translateX(0);
  }
  
  .admin-content {
    width: 100vw;
    margin-left: 0;
  }
  
  .admin-topbar {
    padding: 0 16px;
  }
  
  .username {
    display: none;
  }
}

@media (max-width: 480px) {
  .admin-topbar {
    padding: 0 12px;
  }
  
  .topbar-left {
    font-size: 0.9rem;
  }
}
</style>