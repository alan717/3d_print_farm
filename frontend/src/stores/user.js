import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref({
    id: null,
    username: '',
    email: ''
  })
  const token = ref(localStorage.getItem('token') || '')

  const isLoggedIn = computed(() => !!token.value)

  const login = async (credentials) => {
    try {
      const response = await axios.post('/api/auth/login', credentials)
      const { token: authToken, user } = response.data
      
      token.value = authToken
      userInfo.value = user
      localStorage.setItem('token', authToken)
      
      // 设置axios默认headers
      axios.defaults.headers.common['Authorization'] = `Bearer ${authToken}`
      
      return { success: true }
    } catch (error) {
      return { 
        success: false, 
        message: error.response?.data?.message || '登录失败' 
      }
    }
  }

  const register = async (userData) => {
    try {
      await axios.post('/api/auth/register', userData)
      return { success: true, message: '注册成功' }
    } catch (error) {
      return { 
        success: false, 
        message: error.response?.data?.message || '注册失败' 
      }
    }
  }

  const logout = () => {
    token.value = ''
    userInfo.value = { id: null, username: '', email: '' }
    localStorage.removeItem('token')
    delete axios.defaults.headers.common['Authorization']
  }

  const fetchUserInfo = async () => {
    if (!token.value) return
    
    try {
      const response = await axios.get('/api/user/profile')
      userInfo.value = response.data
    } catch (error) {
      console.error('获取用户信息失败:', error)
      logout()
    }
  }

  // 初始化时设置token
  if (token.value) {
    axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
    fetchUserInfo()
  }

  const setToken = (authToken) => {
    token.value = authToken
    localStorage.setItem('token', authToken)
    axios.defaults.headers.common['Authorization'] = `Bearer ${authToken}`
  }

  const setUserInfo = (user) => {
    userInfo.value = user
    localStorage.setItem('userInfo', JSON.stringify(user))
  }

  return {
    userInfo,
    token,
    isLoggedIn,
    login,
    register,
    logout,
    fetchUserInfo,
    setToken,
    setUserInfo
  }
})