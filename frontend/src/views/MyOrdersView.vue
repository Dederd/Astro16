<template>
  <div class="my-orders-page fade-in">
    <div class="page-header">
      <h1>📦 Pesanan Saya</h1>
      <p>Pantau status semua pesanan bouquetmu di sini</p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="orders-list">
      <div v-for="i in 3" :key="i" class="skeleton" style="height: 140px; border-radius: 16px;"></div>
    </div>

    <!-- Error -->
    <div v-else-if="loadError" class="empty-state">
      <span>⚠️</span>
      <p>Gagal memuat pesanan. Coba refresh halaman.</p>
    </div>

    <!-- Empty -->
    <div v-else-if="orders.length === 0" class="empty-state">
      <span>🌸</span>
      <p>Belum ada pesanan.</p>
      <button class="btn btn-primary" @click="router.push('/order')">Buat Bouquet Sekarang</button>
    </div>

    <!-- Orders list -->
    <div v-else class="orders-list">
      <div
        v-for="order in orders"
        :key="order.id"
        class="order-card"
        @click="toggleDetail(order.id)"
      >
        <!-- Card header -->
        <div class="order-header">
          <div class="order-id-wrap">
            <span class="order-id">#{{ order.id }}</span>
            <span class="order-date">{{ formatDate(order.created_at) }}</span>
          </div>
          <div class="order-status-wrap">
            <span class="status-badge" :class="statusClass(order.status)">
              {{ statusLabel(order.status) }}
            </span>
            <span class="toggle-icon">{{ expandedId === order.id ? '▲' : '▼' }}</span>
          </div>
        </div>

        <!-- Summary row -->
        <div class="order-summary">
          <span class="order-design">{{ order.design_name || order.catalog_name || 'Custom Bouquet' }}</span>
          <span class="order-price">Rp{{ formatPrice(order.total_amount) }}</span>
        </div>

        <!-- Expanded detail -->
        <transition name="expand">
          <div v-if="expandedId === order.id" class="order-detail" @click.stop>
            <!-- Tracking steps -->
            <div class="tracking-steps">
              <div
                v-for="(step, i) in trackingSteps(order.status)"
                :key="i"
                class="tracking-step"
                :class="{ done: step.done, active: step.active }"
              >
                <div class="step-dot">
                  <span v-if="step.done">✓</span>
                  <span v-else-if="step.active" class="pulse-dot"></span>
                  <span v-else>{{ i + 1 }}</span>
                </div>
                <div class="step-info">
                  <span class="step-label">{{ step.label }}</span>
                  <span v-if="step.time" class="step-time">{{ step.time }}</span>
                </div>
                <div v-if="i < trackingSteps(order.status).length - 1" class="step-line" :class="{ done: step.done }"></div>
              </div>
            </div>

            <!-- Order info -->
            <div class="detail-grid">
              <div class="detail-item">
                <span class="detail-label">Penerima</span>
                <span class="detail-val">{{ order.customer_name }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">No. HP</span>
                <span class="detail-val">{{ order.customer_phone }}</span>
              </div>
              <div class="detail-item" style="grid-column: 1/-1">
                <span class="detail-label">Alamat Pengiriman</span>
                <span class="detail-val">{{ order.shipping_address }}, {{ order.shipping_city }} {{ order.shipping_postcode }}</span>
              </div>
              <div v-if="order.courier_service" class="detail-item">
                <span class="detail-label">Kurir</span>
                <span class="detail-val">{{ order.courier_service?.toUpperCase() }}</span>
              </div>
              <div v-if="order.tracking_number && order.status !== 'pending' && order.status !== 'paid'" class="detail-item">
                <span class="detail-label">No. Resi</span>
                <span class="detail-val monospace" style="font-weight: 600; color: var(--deep-rose);">{{ order.tracking_number }}</span>
              </div>
              <div v-if="order.payment_id" class="detail-item">
                <span class="detail-label">ID Pembayaran</span>
                <span class="detail-val monospace">{{ order.payment_id }}</span>
              </div>
            </div>

            <!-- Notes -->
            <div v-if="order.notes" class="order-notes">
              📝 {{ order.notes }}
            </div>

            <!-- Action buttons -->
            <div class="order-actions">
              <!-- Lanjut Pembayaran (for pending orders) -->
              <button
                v-if="order.status === 'pending'"
                class="btn btn-primary btn-sm"
                :disabled="paymentLoadingId === order.id"
                @click.stop="continuePayment(order)"
              >
                <span v-if="paymentLoadingId === order.id">⏳ Memproses...</span>
                <span v-else>💳 Lanjut Pembayaran</span>
              </button>

              <!-- Hubungi Penjual -->
              <div class="contact-buttons">
                <a
                  :href="`https://wa.me/+62823235092101?text=Halo%2C%20saya%20ingin%20menanyakan%20tentang%20pesanan%20saya%20nomor%20${order.id}`"
                  target="_blank"
                  class="btn btn-outline btn-sm whatsapp"
                  title="Chat via WhatsApp"
                >
                  💬 WhatsApp
                </a>
                <a
                  href="mailto:astroreborn441@gmail.com?subject=Pertanyaan%20Pesanan%20Bouquet"
                  class="btn btn-outline btn-sm email"
                  title="Email"
                >
                  ✉️ Email
                </a>
              </div>
            </div>
          </div>
        </transition>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getUserOrders, getOrder, createPaymentToken, notifyAdminNewOrder } from '@/services/api'

const router = useRouter()
const route = useRoute()
const orders = ref([])
const loading = ref(true)
const loadError = ref(false)
const expandedId = ref(null)

async function loadOrders() {
  loading.value = true
  loadError.value = false
  try {
    const res = await getUserOrders()
    orders.value = (res.data.data || []).sort((a, b) =>
      new Date(b.created_at) - new Date(a.created_at)
    )
  } catch (e) {
    console.error('[MyOrders] gagal load:', e)
    loadError.value = true
  } finally {
    loading.value = false
  }
}

// Jika ada order_id dari query (setelah payment), pastikan status sudah terupdate
async function checkAndRefreshPaidOrder(orderId) {
  if (!orderId) return
  try {
    const res = await getOrder(orderId)
    const freshOrder = res.data.data
    if (freshOrder) {
      const idx = orders.value.findIndex(o => o.id === orderId)
      if (idx !== -1) {
        orders.value[idx] = freshOrder
      } else {
        await loadOrders()
      }
    }
  } catch (e) {
    console.error('[MyOrders] gagal refresh order:', e)
  }
}

onMounted(async () => {
  await loadOrders()
  // Jika navigasi dari halaman payment finish, refresh order terkait
  if (route.query.order_id) {
    await checkAndRefreshPaidOrder(route.query.order_id)
  }
})

// Watch jika query berubah (misal user kembali ke halaman ini)
watch(() => route.query.order_id, async (newId) => {
  if (newId) {
    await checkAndRefreshPaidOrder(newId)
  }
})

function toggleDetail(id) {
  expandedId.value = expandedId.value === id ? null : id
}

function formatDate(dt) {
  if (!dt) return '-'
  return new Date(dt).toLocaleDateString('id-ID', {
    day: 'numeric', month: 'long', year: 'numeric',
    hour: '2-digit', minute: '2-digit'
  })
}

function formatPrice(p) {
  return (p || 0).toLocaleString('id-ID')
}

const STATUS_MAP = {
  pending:    { label: '⏳ Menunggu Pembayaran', cls: 'pending' },
  paid:       { label: '✅ Dibayar',              cls: 'paid' },
  processing: { label: '🔧 Sedang Diproses',      cls: 'processing' },
  shipped:    { label: '🚚 Dalam Pengiriman',      cls: 'shipped' },
  delivered:  { label: '🎉 Terkirim',              cls: 'delivered' },
  cancelled:  { label: '❌ Dibatalkan',             cls: 'cancelled' },
}

function statusLabel(s) { return STATUS_MAP[s]?.label || s }
function statusClass(s) { return STATUS_MAP[s]?.cls || 'pending' }

const STEPS = [
  { key: 'pending',    label: 'Pesanan Diterima' },
  { key: 'paid',       label: 'Pembayaran Dikonfirmasi' },
  { key: 'processing', label: 'Sedang Dirangkai' },
  { key: 'shipped',    label: 'Dalam Pengiriman' },
  { key: 'delivered',  label: 'Terkirim 🎉' },
]

function trackingSteps(status) {
  if (status === 'cancelled') {
    return [{ label: 'Pesanan Dibatalkan', done: true, active: false }]
  }
  const currentIdx = STEPS.findIndex(s => s.key === status)
  return STEPS.map((s, i) => ({
    label: s.label,
    done: i < currentIdx,
    active: i === currentIdx,
    time: null,
  }))
}

const paymentLoadingId = ref(null)

async function pollOrderStatus(orderId, maxRetries = 15, delayMs = 500) {
  for (let i = 0; i < maxRetries; i++) {
    try {
      const res = await getOrder(orderId)
      const o = res.data.data
      if (o?.status === 'paid') return o
      if (i < maxRetries - 1) await new Promise(r => setTimeout(r, delayMs))
    } catch (e) { console.error('[pollOrderStatus] error:', e) }
  }
  return null
}

async function continuePayment(order) {
  if (paymentLoadingId.value) return
  paymentLoadingId.value = order.id
  try {
    const res = await createPaymentToken(order.id)
    const data = res.data.data

    // Sudah terbayar (backend detect dari Midtrans)
    if (data?.status === 'paid') {
      const idx = orders.value.findIndex(o => o.id === order.id)
      if (idx !== -1) orders.value[idx] = { ...orders.value[idx], status: 'paid' }
      alert('Pembayaran sudah berhasil dikonfirmasi.')
      return
    }

    const { token, redirect_url } = data || {}

    if (!token && !redirect_url) {
      alert('Gagal mendapatkan token pembayaran.')
      return
    }

    if (window.snap && token) {
      window.snap.pay(token, {
        onSuccess: async () => {
          const paidOrder = await pollOrderStatus(order.id)
          const finalOrder = paidOrder || (await getOrder(order.id)).data.data
          if (finalOrder) {
            const idx = orders.value.findIndex(o => o.id === order.id)
            if (idx !== -1) orders.value[idx] = finalOrder
          }
          if (finalOrder?.status === 'paid') {
            try {
              await notifyAdminNewOrder({
                order_id: order.id,
                customer_name: finalOrder.customer_name,
                customer_phone: finalOrder.customer_phone,
                total_amount: finalOrder.total_amount,
                design_name: finalOrder.design_name || finalOrder.catalog_item_id,
              })
            } catch (e) { console.error('[notifyAdmin] error:', e) }
          }
          router.push({ path: '/payment/finish', query: { order_id: order.id, status: 'success' } })
        },
        onPending: () => {
          alert('Pembayaran sedang diproses. Harap tunggu konfirmasi.')
        },
        onError: () => {
          alert('Pembayaran gagal. Silakan coba lagi.')
        },
        onClose: () => { paymentLoadingId.value = null },
      })
    } else if (redirect_url) {
      window.location.href = redirect_url
    } else {
      alert('Gagal mendapatkan token pembayaran.')
    }
  } catch (e) {
    console.error('[continuePayment] error:', e)
    const status = e?.response?.status
    const errData = e?.response?.data
    if (status === 410 || errData?.status === 'expire' || errData?.status === 'cancel') {
      alert('Sesi pembayaran sudah expired. Silakan hubungi kami untuk membuat pesanan baru.')
    } else {
      alert('Gagal memproses pembayaran: ' + (errData?.details || errData?.error || e.message))
    }
  } finally {
    paymentLoadingId.value = null
  }
}
</script>

<style scoped>
.my-orders-page {
  max-width: 720px;
  margin: 0 auto;
  padding-bottom: 80px;
}

.page-header {
  text-align: center;
  padding: 48px 0 32px;
}
.page-header h1 { font-size: clamp(1.8rem, 4vw, 2.4rem); margin-bottom: 8px; }
.page-header p { color: var(--warm-gray); font-size: 0.95rem; }

.orders-list { display: flex; flex-direction: column; gap: 16px; }

.order-card {
  background: var(--white);
  border: 1.5px solid var(--light-gray);
  border-radius: var(--radius);
  padding: 20px 24px;
  cursor: pointer;
  transition: var(--transition);
}
.order-card:hover { border-color: var(--rose); box-shadow: 0 4px 20px rgba(198,40,40,0.08); }

.order-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 10px;
  flex-wrap: wrap;
}

.order-id-wrap { display: flex; flex-direction: column; gap: 2px; }
.order-id { font-size: 0.8rem; font-weight: 600; color: var(--charcoal); font-family: monospace; }
.order-date { font-size: 0.75rem; color: var(--warm-gray); }

.order-status-wrap { display: flex; align-items: center; gap: 10px; }

.status-badge {
  font-size: 0.75rem;
  font-weight: 600;
  padding: 4px 12px;
  border-radius: var(--radius-pill);
}
.status-badge.pending    { background: #FFF8E1; color: #F57F17; }
.status-badge.paid       { background: #E8F5E9; color: #2E7D32; }
.status-badge.processing { background: #E3F2FD; color: #1565C0; }
.status-badge.shipped    { background: #FFF3E0; color: #E65100; }
.status-badge.delivered  { background: #F3E5F5; color: #6A1B9A; }
.status-badge.cancelled  { background: #FFEBEE; color: #C62828; }

.toggle-icon { font-size: 0.65rem; color: var(--warm-gray); }

.order-summary {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}
.order-design { font-size: 0.9rem; font-weight: 500; color: var(--charcoal); }
.order-price { font-size: 1rem; font-weight: 600; color: var(--deep-rose); }

/* Expanded detail */
.order-detail {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid var(--light-gray);
}

/* Tracking */
.tracking-steps {
  display: flex;
  gap: 0;
  margin-bottom: 24px;
  overflow-x: auto;
  padding-bottom: 8px;
}

.tracking-step {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
  min-width: 80px;
  position: relative;
  text-align: center;
}

.step-dot {
  width: 32px; height: 32px;
  border-radius: 50%;
  background: var(--light-gray);
  color: var(--warm-gray);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  font-weight: 700;
  margin-bottom: 8px;
  position: relative;
  z-index: 1;
  transition: var(--transition);
}

.tracking-step.done .step-dot {
  background: var(--deep-rose);
  color: white;
}

.tracking-step.active .step-dot {
  background: var(--rose);
  color: white;
  box-shadow: 0 0 0 4px rgba(198,40,40,0.2);
}

.pulse-dot {
  width: 10px; height: 10px;
  background: white;
  border-radius: 50%;
  animation: pulse 1.2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.6; transform: scale(0.8); }
}

.step-line {
  position: absolute;
  top: 16px;
  left: calc(50% + 16px);
  width: calc(100% - 32px);
  height: 2px;
  background: var(--light-gray);
  z-index: 0;
}
.step-line.done { background: var(--deep-rose); }

.step-info { display: flex; flex-direction: column; gap: 2px; }
.step-label { font-size: 0.68rem; color: var(--warm-gray); line-height: 1.3; }
.tracking-step.done .step-label,
.tracking-step.active .step-label { color: var(--charcoal); font-weight: 500; }
.step-time { font-size: 0.62rem; color: var(--warm-gray); }

/* Detail grid */
.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px 20px;
  margin-bottom: 16px;
}

.detail-item { display: flex; flex-direction: column; gap: 3px; }
.detail-label { font-size: 0.72rem; color: var(--warm-gray); font-weight: 500; text-transform: uppercase; letter-spacing: 0.05em; }
.detail-val { font-size: 0.85rem; color: var(--charcoal); }
.monospace { font-family: monospace; font-size: 0.8rem; }

.order-notes {
  background: var(--cream);
  border-radius: var(--radius-sm);
  padding: 10px 14px;
  font-size: 0.82rem;
  color: var(--warm-gray);
  font-style: italic;
}

/* Expand transition */
.expand-enter-active, .expand-leave-active { transition: all 0.28s ease; overflow: hidden; }
.expand-enter-from, .expand-leave-to { opacity: 0; max-height: 0; padding-top: 0; }
.expand-enter-to, .expand-leave-from { opacity: 1; max-height: 600px; }

.empty-state { text-align: center; padding: 80px 20px; color: var(--warm-gray); }
.empty-state span { font-size: 3rem; display: block; margin-bottom: 16px; }
.empty-state p { margin-bottom: 20px; }

/* Order actions */
.order-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid var(--light-gray);
  align-items: center;
}

.contact-buttons {
  display: flex;
  gap: 8px;
  margin-left: auto;
}

.btn-sm {
  padding: 8px 16px;
  font-size: 0.82rem;
}

.whatsapp {
  background: #25D366 !important;
  color: white !important;
  border-color: #25D366 !important;
}
.whatsapp:hover {
  background: #1EAD52 !important;
  border-color: #1EAD52 !important;
}

.email {
  background: #EA4335 !important;
  color: white !important;
  border-color: #EA4335 !important;
}
.email:hover {
  background: #D33425 !important;
  border-color: #D33425 !important;
}

@media (max-width: 640px) {
  .order-actions {
    flex-direction: column;
  }
  .contact-buttons {
    margin-left: 0;
    width: 100%;
  }
  .contact-buttons .btn {
    flex: 1;
  }
}
</style>
