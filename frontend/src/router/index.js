import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'
import Order from '@/views/Order.vue'
import CatalogView from '@/views/CatalogView.vue'
import PaymentFinishView from '@/views/PaymentFinishView.vue'
import AdminView from '@/views/AdminView.vue'
import MyOrdersView from '@/views/MyOrdersView.vue'

const routes = [
  { path: '/', component: Home },
  { path: '/order', component: Order },
  { path: '/catalog', component: CatalogView },
  { path: '/payment/finish', component: PaymentFinishView },
  { path: '/admin', component: AdminView },
  { path: '/my-orders', component: MyOrdersView, meta: { requiresAuth: true } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

// Navigation guard: halaman yang butuh login
router.beforeEach((to, from, next) => {
  if (to.meta.requiresAuth) {
    const token = localStorage.getItem('bouquet_token')
    if (!token) {
      return next({ path: '/', query: { requireLogin: '1', redirect: to.fullPath } })
    }
  }
  next()
})

export default router
