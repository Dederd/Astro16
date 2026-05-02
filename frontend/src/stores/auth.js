import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/services/api'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref(JSON.parse(localStorage.getItem('bouquet_user') || 'null'))
  const token = ref(localStorage.getItem('bouquet_token') || null)
  const loading = ref(false)
  const error = ref('')

  // Computed
  const isLoggedIn = computed(() => !!token.value && !!user.value)

  // Inject token ke setiap request API jika login
  function setupInterceptor() {
    api.interceptors.request.use(config => {
      if (token.value) {
        config.headers['Authorization'] = `Bearer ${token.value}`
      }
      return config
    })
  }
  setupInterceptor()

  async function register(name, email, password, phone) {
    loading.value = true
    error.value = ''
    try {
      const res = await api.post('/auth/register', { name, email, password, phone })
      _saveSession(res.data.data)
      return { success: true }
    } catch (e) {
      error.value = e.response?.data?.error || 'Registrasi gagal'
      return { success: false, message: error.value }
    } finally {
      loading.value = false
    }
  }

  async function login(email, password) {
    loading.value = true
    error.value = ''
    try {
      const res = await api.post('/auth/login', { email, password })
      _saveSession(res.data.data)
      return { success: true }
    } catch (e) {
      error.value = e.response?.data?.error || 'Email atau password salah'
      return { success: false, message: error.value }
    } finally {
      loading.value = false
    }
  }

  function logout() {
    user.value = null
    token.value = null
    localStorage.removeItem('bouquet_user')
    localStorage.removeItem('bouquet_token')
  }

  function _saveSession(data) {
    token.value = data.token
    user.value = data.user
    localStorage.setItem('bouquet_token', data.token)
    localStorage.setItem('bouquet_user', JSON.stringify(data.user))
  }

  return {
    user, token, loading, error, isLoggedIn,
    register, login, logout,
  }
})
