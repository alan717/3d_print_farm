import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'
import Login from '@/views/Login.vue'
import Register from '@/views/Register.vue'
import Upload from '@/views/Upload.vue'
import Cart from '@/views/Cart.vue'
import Orders from '@/views/Orders.vue'
import Profile from '@/views/Profile.vue'
import AdminLogin from '@/views/AdminLogin.vue'

// Admin components
import AdminLayout from '@/components/AdminLayout.vue'
import AdminDashboard from '@/views/admin/Dashboard.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/register',
    name: 'Register',
    component: Register
  },
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: AdminLogin
  },
  {
    path: '/upload',
    name: 'Upload',
    component: Upload,
    meta: { requiresAuth: true }
  },
  {
    path: '/cart',
    name: 'Cart',
    component: Cart
  },
  {
    path: '/orders',
    name: 'Orders',
    component: Orders,
    meta: { requiresAuth: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: Profile,
    meta: { requiresAuth: true }
  },
  {
    path: '/admin',
    component: AdminLayout,
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      {
        path: '',
        redirect: '/admin/dashboard'
      },
      {
        path: 'dashboard',
        name: 'AdminDashboard',
        component: AdminDashboard
      },
      {
        path: 'users',
        name: 'AdminUsers',
        component: () => import('@/views/admin/Users.vue')
      },
      {
        path: 'orders',
        name: 'AdminOrders',
        component: () => import('@/views/admin/Orders.vue')
      },
      {
        path: 'materials',
        name: 'AdminMaterials',
        component: () => import('@/views/admin/Materials.vue')
      },
      {
        path: 'settings',
        name: 'AdminSettings',
        component: () => import('@/views/admin/Settings.vue')
      },
      {
        path: 'logs',
        name: 'AdminLogs',
        component: () => import('@/views/admin/Logs.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const isAuthenticated = localStorage.getItem('token')
  const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
  
  // 检查是否需要登录
  if (to.meta.requiresAuth && !isAuthenticated) {
    // 如果访问的是管理员页面，重定向到管理员登录页
    if (to.path.startsWith('/admin')) {
      next('/admin/login')
    } else {
      next('/login')
    }
    return
  }
  
  // 检查是否需要管理员权限
  if (to.meta.requiresAdmin) {
    if (!userInfo.role || (userInfo.role !== 'admin' && userInfo.role !== 'super_admin')) {
      // 没有管理员权限，重定向到管理员登录页
      next('/admin/login')
      return
    }
  }
  
  // 如果已登录的管理员访问管理员登录页，重定向到仪表板
  if (to.path === '/admin/login' && isAuthenticated && 
      userInfo.role && (userInfo.role === 'admin' || userInfo.role === 'super_admin')) {
    next('/admin/dashboard')
    return
  }
  
  next()
})

export default router