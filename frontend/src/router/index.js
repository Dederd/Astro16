import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'
import Order from '@/views/Order.vue'
import CatalogView from '@/views/CatalogView.vue'
import PaymentFinishView from '@/views/PaymentFinishView.vue'
import AdminView from '@/views/AdminView.vue'

const routes = [
  { path: '/', component: Home },
  { path: '/order', component: Order },
  { path: '/catalog', component: CatalogView },
  { path: '/payment/finish', component: PaymentFinishView },
  { path: '/admin', component: AdminView },
]

export default createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})
