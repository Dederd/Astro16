<template>
  <div class="finish-page fade-in">
    <!-- Success state -->
    <div v-if="status === 'success'" class="status-card success">
      <div class="status-icon">🎉</div>
      <h1>Pembayaran Berhasil!</h1>
      <p>Terima kasih, <strong>{{ order?.customer_name }}</strong>! Pesananmu sedang kami proses.</p>
      <div class="order-id-box">
        <span>ID Pesanan</span>
        <strong>{{ orderId }}</strong>
      </div>
    </div>

    <!-- Pending state -->
    <div v-else-if="status === 'pending'" class="status-card pending">
      <div class="status-icon">⏳</div>
      <h1>Menunggu Pembayaran</h1>
      <p>Pesananmu sudah dibuat. Selesaikan pembayaran agar pesanan diproses.</p>
      <div class="order-id-box">
        <span>ID Pesanan</span>
        <strong>{{ orderId }}</strong>
      </div>
    </div>

    <!-- Loading order -->
    <div v-if="loadingOrder" class="loading-center">
      <div class="spinner-lg"></div>
      <p>Mengambil detail pesanan...</p>
    </div>

    <!-- Order detail -->
    <div v-if="order && !loadingOrder" class="order-detail card">
      <h2>📋 Detail Pesanan</h2>
      <div class="detail-rows">
        <div class="detail-row">
          <span>Nama</span>
          <span>{{ order.customer_name }}</span>
        </div>
        <div class="detail-row">
          <span>Desain</span>
          <span>{{ order.design_name || order.catalog_item_id || '—' }}</span>
        </div>
        <div class="detail-row">
          <span>Ukuran</span>
          <span>{{ order.size }}</span>
        </div>
        <div class="detail-row">
          <span>Alamat Kirim</span>
          <span>{{ order.shipping_address }}, {{ order.shipping_city }}</span>
        </div>
        <div class="detail-row">
          <span>Kurir</span>
          <span>{{ order.courier_service?.toUpperCase() || '—' }}</span>
        </div>
        <div class="detail-row total">
          <span>Total Dibayar</span>
          <span class="total-val">Rp{{ formatPrice(order.total_amount) }}</span>
        </div>
      </div>
    </div>

    <!-- Tracking section -->
    <div v-if="order && status === 'success'" class="tracking-card card">
      <h2>🚚 Status Pengiriman</h2>

      <div v-if="!order.tracking_number" class="tracking-pending-msg">
        <span class="track-icon">📦</span>
        <div>
          <p><strong>Pesananmu sedang disiapkan</strong></p>
          <p>Nomor resi akan muncul di sini setelah bouquet dikirim (1-2 hari kerja).</p>
        </div>
      </div>

      <div v-else class="tracking-info">
        <div class="tracking-header">
          <div class="tracking-badge">{{ order.courier_service?.toUpperCase() }}</div>
          <div class="tracking-resi">
            <span>No. Resi:</span>
            <strong>{{ order.tracking_number }}</strong>
            <button class="btn-copy" @click="copyResi">{{ copied ? '✓ Disalin' : '📋 Salin' }}</button>
          </div>
        </div>

        <!-- Courier tracking data -->
        <div v-if="loadingTracking" class="tracking-loading">
          <span class="spinner-sm"></span> Mengambil info pengiriman...
        </div>

        <div v-else-if="trackingData" class="tracking-timeline">
          <!-- Status summary -->
          <div class="track-status-summary">
            <div class="track-dot active"></div>
            <div>
              <strong>{{ trackingStatus }}</strong>
              <p>{{ order.shipping_city }}</p>
            </div>
          </div>

          <!-- History -->
          <div v-if="trackingHistory.length > 0" class="track-history">
            <div v-for="(h, i) in trackingHistory" :key="i" class="track-history-item">
              <div class="track-history-dot" :class="{ first: i === 0 }"></div>
              <div class="track-history-body">
                <span class="track-history-date">{{ h.date }}</span>
                <span class="track-history-desc">{{ h.desc }}</span>
                <span v-if="h.location" class="track-history-location">📍 {{ h.location }}</span>
              </div>
            </div>
          </div>
        </div>

        <a
          :href="courierTrackUrl"
          target="_blank"
          class="btn btn-outline btn-sm track-external-btn"
          v-if="order.tracking_number"
        >
          🔗 Lacak di situs {{ order.courier_service?.toUpperCase() }}
        </a>
      </div>
    </div>

    <!-- Actions -->
    <div class="finish-actions">
      <button class="btn btn-primary" @click="$router.push('/')">🏠 Kembali ke Beranda</button>
      <button class="btn btn-outline" @click="$router.push('/catalog')">🛍️ Lihat Katalog Lagi</button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getOrder, getOrderTracking } from '@/services/api'

const route = useRoute()

const orderId = computed(() => route.query.order_id)
const status = computed(() => route.query.status || 'success')

const order = ref(null)
const loadingOrder = ref(true)
const loadingTracking = ref(false)
const trackingData = ref(null)
const copied = ref(false)

onMounted(async () => {
  if (!orderId.value) { loadingOrder.value = false; return }
  try {
    const res = await getOrder(orderId.value)
    order.value = res.data.data
  } catch (e) {
    console.error(e)
  } finally {
    loadingOrder.value = false
  }

  // Fetch tracking jika ada resi
  if (order.value?.tracking_number) {
    loadingTracking.value = true
    try {
      const res = await getOrderTracking(orderId.value)
      trackingData.value = res.data.data
    } catch { /* ignore */ } finally {
      loadingTracking.value = false
    }
  }
})

const trackingStatus = computed(() => {
  return trackingData.value?.courier_data?.summary?.description
    || order.value?.shipping_status
    || 'Dalam Perjalanan'
})

const trackingHistory = computed(() => {
  return trackingData.value?.courier_data?.history || []
})

const courierTrackUrl = computed(() => {
  const resi = order.value?.tracking_number
  const courier = order.value?.courier_service
  if (!resi) return '#'
  if (courier === 'jne') return `https://www.jne.co.id/id/tracking/trace?awb=${resi}`
  if (courier === 'jnt') return `https://www.jet.co.id/track?resi=${resi}`
  if (courier === 'sicepat') return `https://www.sicepat.com/checkAwb?awb=${resi}`
  return '#'
})

async function copyResi() {
  if (order.value?.tracking_number) {
    await navigator.clipboard.writeText(order.value.tracking_number)
    copied.value = true
    setTimeout(() => copied.value = false, 2000)
  }
}

function formatPrice(p) { return (p || 0).toLocaleString('id-ID') }
</script>

<style scoped>
.finish-page { max-width: 680px; margin: 0 auto; padding: 40px 20px 80px; display: flex; flex-direction: column; gap: 24px; }

.status-card {
  border-radius: var(--radius);
  padding: 48px 36px;
  text-align: center;
  border: 1.5px solid;
}
.status-card.success { background: #F1F8E9; border-color: #A5D6A7; }
.status-card.pending { background: #FFF8E1; border-color: #FFD54F; }
.status-icon { font-size: 4rem; margin-bottom: 16px; }
.status-card h1 { font-size: 2rem; margin-bottom: 10px; }
.status-card p { color: var(--warm-gray); font-size: 0.95rem; }

.order-id-box {
  display: inline-flex;
  flex-direction: column;
  gap: 4px;
  background: rgba(255,255,255,0.7);
  border: 1px solid rgba(0,0,0,0.08);
  padding: 12px 24px;
  border-radius: var(--radius);
  margin-top: 20px;
}
.order-id-box span { font-size: 0.75rem; color: var(--warm-gray); }
.order-id-box strong { font-size: 1.1rem; color: var(--charcoal); letter-spacing: 0.05em; }

.loading-center { text-align: center; padding: 40px; }
.spinner-lg { width: 40px; height: 40px; border: 3px solid var(--blush); border-top-color: var(--deep-rose); border-radius: 50%; animation: spin 0.8s linear infinite; margin: 0 auto 16px; }
@keyframes spin { to { transform: rotate(360deg); } }

.order-detail, .tracking-card { padding: 28px; }
.order-detail h2, .tracking-card h2 { font-size: 1.1rem; margin-bottom: 20px; }
.detail-rows { display: flex; flex-direction: column; gap: 0; }
.detail-row { display: flex; justify-content: space-between; padding: 12px 0; font-size: 0.88rem; border-bottom: 1px solid var(--light-gray); color: var(--warm-gray); }
.detail-row span:last-child { color: var(--charcoal); font-weight: 500; text-align: right; max-width: 60%; }
.detail-row.total span { font-weight: 600; font-size: 0.95rem; }
.total-val { color: var(--deep-rose) !important; font-size: 1.1rem !important; }

/* Tracking */
.tracking-pending-msg {
  display: flex;
  gap: 16px;
  align-items: flex-start;
  padding: 16px;
  background: var(--cream);
  border-radius: var(--radius-sm);
  font-size: 0.88rem;
  color: var(--warm-gray);
}
.track-icon { font-size: 1.8rem; flex-shrink: 0; }
.tracking-pending-msg p:first-child { color: var(--charcoal); margin-bottom: 4px; }

.tracking-header { margin-bottom: 16px; }
.tracking-badge { display: inline-block; background: var(--deep-rose); color: white; font-size: 0.75rem; font-weight: 700; padding: 3px 12px; border-radius: var(--radius-pill); margin-bottom: 10px; }
.tracking-resi { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; font-size: 0.9rem; }
.tracking-resi span { color: var(--warm-gray); }
.tracking-resi strong { color: var(--charcoal); letter-spacing: 0.05em; }
.btn-copy { background: none; border: 1px solid var(--light-gray); padding: 3px 10px; border-radius: 6px; font-size: 0.75rem; cursor: pointer; color: var(--warm-gray); transition: var(--transition); }
.btn-copy:hover { background: var(--cream); }

.tracking-loading { font-size: 0.85rem; color: var(--warm-gray); padding: 12px 0; }
.spinner-sm { display: inline-block; width: 14px; height: 14px; border: 2px solid var(--blush); border-top-color: var(--rose); border-radius: 50%; animation: spin 0.8s linear infinite; vertical-align: middle; margin-right: 6px; }

.track-status-summary { display: flex; gap: 14px; align-items: center; padding: 14px 0; border-bottom: 1px solid var(--light-gray); margin-bottom: 16px; }
.track-dot { width: 14px; height: 14px; background: var(--rose); border-radius: 50%; border: 2px solid white; box-shadow: 0 0 0 2px var(--rose); flex-shrink: 0; }
.track-dot.active { background: #43A047; box-shadow: 0 0 0 2px #43A047; }
.track-status-summary strong { font-size: 0.9rem; }
.track-status-summary p { font-size: 0.78rem; color: var(--warm-gray); }

.track-history { display: flex; flex-direction: column; gap: 12px; }
.track-history-item { display: flex; gap: 12px; }
.track-history-dot { width: 10px; height: 10px; background: var(--blush); border-radius: 50%; border: 2px solid var(--rose); flex-shrink: 0; margin-top: 3px; }
.track-history-dot.first { background: var(--rose); }
.track-history-body { display: flex; flex-direction: column; gap: 1px; }
.track-history-date { font-size: 0.72rem; color: var(--warm-gray); }
.track-history-desc { font-size: 0.84rem; color: var(--charcoal); }
.track-history-location { font-size: 0.72rem; color: var(--warm-gray); }

.track-external-btn { margin-top: 16px; }
.btn-sm { padding: 8px 16px; font-size: 0.82rem; }

.finish-actions { display: flex; gap: 12px; justify-content: center; flex-wrap: wrap; padding-top: 8px; }
</style>
