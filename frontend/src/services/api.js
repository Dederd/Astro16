import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1',
  timeout: 30000,
})

// Ambil session ID dari localStorage untuk rate limiting
function getSessionId() {
  return localStorage.getItem('bouquet_session') || 'anonymous'
}

// Inject session header ke setiap request agent
api.interceptors.request.use(config => {
  config.headers['X-Session-ID'] = getSessionId()
  return config
})

// Response interceptor - log error detail
api.interceptors.response.use(
  res => res,
  err => {
    const detail = err.response?.data?.details || err.response?.data?.error || err.message
    console.error(`[API Error] ${err.config?.url} — ${detail}`)
    return Promise.reject(err)
  }
)

export default api

// Auth
export const authRegister = (data) => api.post('/auth/register', data)
export const authLogin = (data) => api.post('/auth/login', data)
export const authMe = () => api.get('/auth/me')

// User orders
export const getUserOrders = () => api.get('/user/orders')

// Bouquet types
export const getBouquetTypes = () => api.get('/bouquet-types')

// Flowers
export const getFlowers = (occasion) =>
  api.get('/flowers', { params: occasion ? { occasion } : {} })

// Catalog
export const getCatalog = (occasion) =>
  api.get('/catalog', { params: occasion ? { occasion } : {} })

// AI Agents
export const agentVerifySelection = (data) =>
  api.post('/agent/verify-selection', data)

export const agentGenerateBouquet = (data) =>
  api.post('/agent/generate-bouquet', data)

export const getGenerateStatus = () =>
  api.get('/agent/generate-status')

export const buyGenerateQuota = () =>
  api.post('/agent/buy-quota')

export const confirmQuotaPayment = (data) =>
  api.post('/agent/confirm-quota', data)

// Orders
export const createOrder = (data) => api.post('/orders', data)
export const getOrder = (id) => api.get(`/orders/${id}`)
export const getOrderTracking = (id) => api.get(`/orders/${id}/tracking`)

// Payment
export const createPaymentToken = (orderId) =>
  api.post('/payment/token', { order_id: orderId })

// Admin
const adminKey = import.meta.env.VITE_ADMIN_KEY || 'admin-bouquet-2024'
const adminApi = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1',
  timeout: 30000,
  headers: { 'X-Admin-Key': adminKey },
})

export const adminGetStats = () => adminApi.get('/admin/stats')
export const adminGetOrders = () => adminApi.get('/admin/orders')
export const adminUpdateOrder = (id, data) => adminApi.put(`/admin/orders/${id}`, data)
export const adminGetFlowers = () => adminApi.get('/admin/flowers')
export const adminCreateFlower = (data) => adminApi.post('/admin/flowers', data)
export const adminUpdateFlower = (id, data) => adminApi.put(`/admin/flowers/${id}`, data)
export const adminDeleteFlower = (id) => adminApi.delete(`/admin/flowers/${id}`)
export const adminGetCatalog = () => adminApi.get('/admin/catalog')
export const adminCreateCatalog = (data) => adminApi.post('/admin/catalog', data)
export const adminUpdateCatalog = (id, data) => adminApi.put(`/admin/catalog/${id}`, data)
export const adminDeleteCatalog = (id) => adminApi.delete(`/admin/catalog/${id}`)

