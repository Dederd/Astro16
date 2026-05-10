<template>
  <div class="step-container fade-in">
    <div class="step-header">
      <h1>Checkout 💳</h1>
      <p>Lengkapi data pemesan, alamat pengiriman, dan lakukan pembayaran</p>
    </div>

    <!-- Login notice -->
    <div v-if="!authStore.isLoggedIn" class="login-notice">
      <span>🔒 Kamu perlu <strong>masuk</strong> atau <strong>daftar</strong> untuk menyelesaikan pembelian.</span>
      <button class="btn btn-primary btn-sm" @click="openAuth && openAuth('login')">Masuk / Daftar</button>
    </div>

    <div class="checkout-layout">
      <!-- Left: Form -->
      <div class="checkout-form-col">
        <!-- Data Pemesan -->
        <div class="card form-section">
          <h2>📋 Data Pemesan</h2>
          <div class="form-group">
            <label class="form-label">Nama Lengkap *</label>
            <input v-model="form.customer_name" type="text" class="form-input" placeholder="Nama lengkapmu" />
          </div>
          <div class="form-group">
            <label class="form-label">Email *</label>
            <input v-model="form.customer_email" type="email" class="form-input" placeholder="email@contoh.com" />
          </div>
          <div class="form-group">
            <label class="form-label">Nomor WhatsApp *</label>
            <input v-model="form.customer_phone" type="tel" class="form-input" placeholder="08xxxxxxxxxx" />
          </div>
          <div class="form-group">
            <label class="form-label">Catatan (opsional)</label>
            <textarea v-model="form.notes" class="form-input" placeholder="Pesan spesial..." rows="3" style="resize:vertical;"></textarea>
          </div>
        </div>

        <!-- Alamat Pengiriman -->
        <div class="card form-section">
          <h2>📦 Alamat Pengiriman</h2>
          <div class="form-group">
            <label class="form-label">Nama Penerima *</label>
            <input v-model="form.shipping_recipient" type="text" class="form-input" placeholder="Nama penerima bouquet" />
          </div>
          <div class="form-group">
            <label class="form-label">No. HP Penerima *</label>
            <input v-model="form.shipping_phone" type="tel" class="form-input" placeholder="08xxxxxxxxxx" />
          </div>
          <div class="form-group">
            <label class="form-label">Alamat Lengkap *</label>
            <textarea v-model="form.shipping_address" class="form-input" placeholder="Jl. contoh No. 1, RT/RW, Kelurahan, Kecamatan" rows="3" style="resize:vertical;"></textarea>
          </div>
          <div class="form-row">
            <div class="form-group">
              <label class="form-label">Kota *</label>
              <input v-model="form.shipping_city" type="text" class="form-input" placeholder="Jakarta Selatan" />
            </div>
            <div class="form-group">
              <label class="form-label">Kode Pos *</label>
              <input v-model="form.shipping_postcode" type="text" class="form-input" placeholder="12345" maxlength="5" />
            </div>
          </div>
        </div>

        <!-- Pilih Kurir -->
        <div class="card form-section">
          <h2>🚚 Pilih Kurir</h2>
          <div class="courier-grid">
            <div
              v-for="courier in couriers"
              :key="courier.value"
              class="courier-option"
              :class="{ active: form.courier_service === courier.value }"
              @click="selectCourier(courier)"
            >
              <div class="courier-logo">{{ courier.logo }}</div>
              <div class="courier-info">
                <span class="courier-name">{{ courier.name }}</span>
                <span class="courier-desc">{{ courier.desc }}</span>
              </div>
              <span class="courier-price">Rp{{ formatPrice(courier.fee) }}</span>
              <div v-if="form.courier_service === courier.value" class="courier-check">✓</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right: Order summary -->
      <div class="order-summary-col">
        <div class="summary-card card">
          <h2>Ringkasan Pesanan</h2>

          <!-- Design preview -->
          <div class="design-preview">
            <div class="preview-visual" :style="`background: ${gradientFor(designStyle)}`">
              <span class="preview-emoji">{{ emojiFor(designStyle) }}</span>
            </div>
            <div class="preview-info">
              <span class="preview-name">{{ designName }}</span>
              <span class="preview-style">{{ designStyle }}</span>
              <span class="preview-occasion">{{ store.selectedBouquetType?.name }}</span>
            </div>
          </div>

          <!-- Bunga breakdown -->
          <div v-if="store.orderMode !== 'catalog' && store.selectedFlowers.length" class="summary-flowers-section">
            <span class="summary-flowers-label">Komposisi Bunga</span>
            <div class="summary-flowers-list">
              <div v-for="f in store.selectedFlowers" :key="f.flower_id" class="summary-flower-row">
                <span>🌸 {{ f.name }}</span>
              </div>
            </div>
          </div>

          <div class="summary-divider"></div>

          <!-- Price breakdown -->
          <div class="price-breakdown">
            <div class="breakdown-row">
              <span>🌸 Harga Bunga</span>
              <span>Rp{{ formatPrice(store.priceBreakdown?.flower_cost || 0) }}</span>
            </div>
            <div class="breakdown-row">
              <span>🎀 Biaya Pembuatan</span>
              <span>Rp{{ formatPrice(store.MAKING_FEE) }}</span>
            </div>
            <div v-if="store.aiFeePaid" class="breakdown-row ai-fee">
              <span>🤖 Biaya AI Generate</span>
              <span>Rp{{ formatPrice(store.AI_FEE) }}</span>
            </div>
            <div v-if="store.extraQuotaFee > 0" class="breakdown-row ai-fee">
              <span>✨ Kuota Generate Tambahan</span>
              <span>Rp{{ formatPrice(store.extraQuotaFee) }}</span>
            </div>
            <div class="breakdown-row">
              <span>🚚 Ongkos Kirim ({{ form.courier_service ? form.courier_service.toUpperCase() : '—' }})</span>
              <span>{{ form.courier_service ? 'Rp' + formatPrice(store.shippingCost) : '—' }}</span>
            </div>
          </div>

          <div class="summary-total">
            <span>Total</span>
            <span class="total-price">Rp{{ formatPrice(store.totalPrice) }}</span>
          </div>
        </div>

        <!-- Error -->
        <div v-if="error" class="error-msg">⚠️ {{ error }}</div>

        <!-- Pay button -->
        <button
          class="btn btn-gold pay-btn"
          :disabled="!isFormValid || paymentLoading"
          @click="handlePayment"
        >
          <span v-if="paymentLoading">
            <span class="spinner-sm"></span> Memproses...
          </span>
          <span v-else>💳 Bayar Sekarang — Rp{{ formatPrice(store.totalPrice) }}</span>
        </button>

        <p class="payment-note">
          🔒 Pembayaran aman diproses oleh Midtrans.<br>
          Tersedia: Transfer Bank, GoPay, OVO, QRIS, dll.
        </p>
        <button class="btn btn-ghost back-btn" @click="store.setStep(3)">← Ubah Desain</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, reactive, inject, watch } from 'vue'
import { useOrderStore } from '@/stores/order'
import { useAuthStore } from '@/stores/auth'
import { createOrder, createPaymentToken } from '@/services/api'
import { useRouter } from 'vue-router'

const store = useOrderStore()
const authStore = useAuthStore()
const router = useRouter()
const openAuth = inject('openAuth', null)

const form = reactive({
  customer_name: '',
  customer_email: '',
  customer_phone: '',
  notes: '',
  shipping_recipient: '',
  shipping_phone: '',
  shipping_address: '',
  shipping_city: '',
  shipping_postcode: '',
  courier_service: '',
})

// Auto-fill dari data akun user
watch(() => authStore.user, (u) => {
  if (u) {
    if (!form.customer_name) form.customer_name = u.name || ''
    if (!form.customer_email) form.customer_email = u.email || ''
    if (!form.customer_phone) form.customer_phone = u.phone || ''
  }
}, { immediate: true })

const couriers = [
  { value: 'jne',     name: 'JNE',          logo: '📦', desc: 'Reguler 2-3 hari',    fee: 15000 },
  { value: 'jnt',     name: 'J&T Express',  logo: '🚚', desc: 'Reguler 1-2 hari',    fee: 12000 },
  { value: 'sicepat', name: 'SiCepat',      logo: '⚡', desc: 'Same day tersedia',    fee: 18000 },
]

function selectCourier(courier) {
  form.courier_service = courier.value
  store.setShippingCost(courier.fee)
}

const paymentLoading = ref(false)
const error = ref('')

const designName = computed(() => {
  if (store.orderMode === 'catalog') return store.selectedCatalogItem?.name || '—'
  return store.selectedDesign?.name || '—'
})

const designStyle = computed(() => {
  if (store.orderMode === 'catalog') return store.selectedCatalogItem?.style || 'Classic'
  return store.selectedDesign?.style || 'Romantic'
})

const stemCount = computed(() => {
  if (store.orderMode === 'catalog') return store.selectedCatalogItem?.stem_count || 0
  return store.selectedSize === 'small'
    ? store.selectedDesign?.small?.stem_count
    : store.selectedDesign?.large?.stem_count
})

const isFormValid = computed(() =>
  form.customer_name.trim() &&
  form.customer_email.trim() &&
  form.customer_phone.trim() &&
  form.shipping_address.trim() &&
  form.shipping_city.trim() &&
  form.shipping_postcode.trim() &&
  form.shipping_phone.trim() &&
  form.courier_service
)

async function handlePayment() {
  if (!authStore.isLoggedIn) {
    if (openAuth) openAuth('login')
    return
  }
  if (!isFormValid.value) {
    error.value = 'Lengkapi semua field yang wajib diisi (*).'
    return
  }
  // Bug fix 1: prevent payment if total is 0 (design not selected)
  if (!store.totalPrice || store.totalPrice <= 0) {
    error.value = 'Terjadi kesalahan: total harga tidak valid. Kembali dan pilih desain bouquet.'
    return
  }
  error.value = ''
  paymentLoading.value = true

  try {
    const isCatalog = store.orderMode === 'catalog'
    const breakdown = store.priceBreakdown

    const orderPayload = {
      customer_name:     form.customer_name,
      customer_email:    form.customer_email,
      customer_phone:    form.customer_phone,
      bouquet_type_id:   store.selectedBouquetType?.id || '',
      selected_flowers:  isCatalog ? [] : store.selectedFlowers,
      design_id:         isCatalog ? store.selectedCatalogItem?.id : store.selectedDesign?.id,
      design_name:       designName.value,
      size:              store.selectedSize || 'catalog',
      total_amount:      store.totalPrice,
      notes:             form.notes,
      shipping_address:  form.shipping_address,
      shipping_city:     form.shipping_city,
      shipping_postcode: form.shipping_postcode,
      shipping_phone:    form.shipping_phone,
      courier_service:   form.courier_service,
      order_source:      isCatalog ? 'catalog' : 'ai_generated',
      catalog_item_id:   isCatalog ? store.selectedCatalogItem?.id : '',
      // Cost breakdown
      flower_cost:        breakdown?.flower_cost || 0,
      making_fee:         breakdown?.making_fee || 0,
      ai_fee:             breakdown?.ai_fee || 0,
      extra_quota_fee:    breakdown?.extra_quota_fee || 0,
      shipping_cost:      breakdown?.shipping_cost || 0,
    }

    const orderRes = await createOrder(orderPayload)
    const order = orderRes.data.data
    store.setCreatedOrder(order)

    const tokenRes = await createPaymentToken(order.id)
    const { token, redirect_url } = tokenRes.data.data

    if (window.snap && token) {
      window.snap.pay(token, {
        onSuccess: () => {
          router.push({ path: '/payment/finish', query: { order_id: order.id, status: 'success' } })
        },
        onPending: () => {
          router.push({ path: '/payment/finish', query: { order_id: order.id, status: 'pending' } })
        },
        onError: (result) => {
          console.error('[Midtrans Snap] onError:', result)
          error.value = 'Pembayaran gagal. Silakan coba lagi.'
          paymentLoading.value = false
        },
        onClose: () => { paymentLoading.value = false },
      })
    } else if (redirect_url) {
      window.location.href = redirect_url
    } else {
      error.value = 'Gagal mendapat token pembayaran.'
      paymentLoading.value = false
    }
  } catch (e) {
    const backendError = e?.response?.data?.details || e?.response?.data?.error || e.message
    error.value = backendError || 'Terjadi kesalahan. Silakan coba lagi.'
    console.error('[Checkout] Error:', e?.response?.data || e)
    paymentLoading.value = false
  }
}

function formatPrice(p) { return (p || 0).toLocaleString('id-ID') }

function gradientFor(style) {
  const g = {
    Romantic: 'linear-gradient(135deg,#FFECD2,#FCB69F)', Elegant: 'linear-gradient(135deg,#E8D5C4,#C4A882)',
    Modern: 'linear-gradient(135deg,#D4E8E0,#8ABFAD)', Natural: 'linear-gradient(135deg,#E8F5E8,#8BC888)',
    Classic: 'linear-gradient(135deg,#FFF8F0,#F5D5B0)', Premium: 'linear-gradient(135deg,#F5E6FF,#C485FF)',
    Warm: 'linear-gradient(135deg,#FFE0E0,#FF9F9F)', Playful: 'linear-gradient(135deg,#FFE0F0,#FF8CC8)',
  }
  return g[style] || g.Romantic
}

function emojiFor(style) {
  const e = { Romantic:'💐',Elegant:'🌹',Modern:'🌿',Natural:'🌻',Classic:'🌸',Premium:'👑',Warm:'❤️',Playful:'🌈' }
  return e[style] || '💐'
}
</script>

<style scoped>
.step-container { max-width: 1020px; margin: 0 auto; }
.step-header { text-align: center; margin-bottom: 32px; }
.step-header h1 { font-size: clamp(1.8rem, 4vw, 2.5rem); margin-bottom: 10px; }
.step-header p { color: var(--warm-gray); }

.login-notice {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  background: linear-gradient(135deg,#FFF8F0,#FFF0E8);
  border: 1.5px solid var(--blush);
  border-radius: var(--radius-sm);
  padding: 14px 20px;
  margin-bottom: 20px;
  font-size: 0.88rem;
  color: var(--warm-gray);
  flex-wrap: wrap;
}
.btn-sm { padding: 8px 18px; font-size: 0.82rem; }

.checkout-layout {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 28px;
  align-items: start;
}

.checkout-form-col { display: flex; flex-direction: column; gap: 20px; }
.form-section { padding: 28px; }
.form-section h2 { font-size: 1.1rem; margin-bottom: 20px; color: var(--charcoal); }
.form-row { display: grid; grid-template-columns: 1fr 120px; gap: 12px; }

/* Courier */
.courier-grid { display: flex; flex-direction: column; gap: 10px; }
.courier-option {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border: 1.5px solid var(--light-gray);
  border-radius: var(--radius);
  cursor: pointer;
  transition: var(--transition);
}
.courier-option:hover { border-color: var(--rose); }
.courier-option.active { border-color: var(--deep-rose); background: #FFF8F8; }
.courier-logo { font-size: 1.5rem; }
.courier-info { flex: 1; }
.courier-name { font-weight: 500; font-size: 0.9rem; display: block; }
.courier-desc { font-size: 0.78rem; color: var(--warm-gray); }
.courier-price { font-size: 0.85rem; font-weight: 500; color: var(--deep-rose); white-space: nowrap; }
.courier-check {
  width: 22px; height: 22px;
  background: var(--deep-rose); color: white;
  border-radius: 50%; display: flex; align-items: center; justify-content: center;
  font-size: 0.75rem; font-weight: 700;
}

/* Summary */
.order-summary-col { position: sticky; top: 80px; display: flex; flex-direction: column; gap: 16px; }
.summary-card { padding: 24px; }
.summary-card h2 { font-size: 1.1rem; margin-bottom: 20px; }

.design-preview {
  display: flex; gap: 14px; align-items: center;
  margin-bottom: 16px; padding-bottom: 16px; border-bottom: 1px solid var(--light-gray);
}
.preview-visual { width: 60px; height: 60px; border-radius: 12px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.preview-emoji { font-size: 1.8rem; }
.preview-info { display: flex; flex-direction: column; gap: 2px; }
.preview-name { font-weight: 500; font-size: 0.95rem; color: var(--charcoal); }
.preview-style, .preview-occasion { font-size: 0.78rem; color: var(--warm-gray); }

.summary-flowers-section { padding: 10px 0; }
.summary-flowers-label { display: block; font-size: 0.72rem; font-weight: 600; text-transform: uppercase; letter-spacing: 0.08em; color: var(--warm-gray); margin-bottom: 6px; }
.summary-flowers-list { display: flex; flex-direction: column; gap: 4px; }
.summary-flower-row { display: flex; justify-content: space-between; font-size: 0.82rem; color: var(--charcoal); }
.flower-subtotal { color: var(--warm-gray); }

.summary-divider { height: 1px; background: var(--light-gray); margin: 12px 0; }

/* Price breakdown */
.price-breakdown { display: flex; flex-direction: column; gap: 8px; margin-bottom: 16px; }
.breakdown-row {
  display: flex; justify-content: space-between;
  font-size: 0.84rem; color: var(--warm-gray);
}
.breakdown-row span:last-child { color: var(--charcoal); font-weight: 500; }
.ai-fee span { color: var(--rose) !important; }

.summary-total { display: flex; justify-content: space-between; align-items: center; font-weight: 500; padding-top: 12px; border-top: 1.5px solid var(--light-gray); }
.total-price { font-family: var(--font-display); font-size: 1.4rem; color: var(--deep-rose); }

.error-msg { background: #FFEBEE; color: #C62828; padding: 12px 16px; border-radius: var(--radius-sm); font-size: 0.85rem; border: 1px solid #FFCDD2; }
.pay-btn { width: 100%; justify-content: center; padding: 16px; font-size: 1rem; }
.payment-note { font-size: 0.78rem; color: var(--warm-gray); text-align: center; line-height: 1.6; }
.back-btn { width: 100%; justify-content: center; font-size: 0.85rem; }

.spinner-sm { display: inline-block; width: 16px; height: 16px; border: 2px solid rgba(255,255,255,0.4); border-top-color: white; border-radius: 50%; animation: spin 0.8s linear infinite; vertical-align: middle; margin-right: 6px; }
@keyframes spin { to { transform: rotate(360deg); } }

/* Mobile */
@media (max-width: 820px) {
  .checkout-layout { grid-template-columns: 1fr; }
  .order-summary-col { position: static; }
  .form-row { grid-template-columns: 1fr; }
  .login-notice { flex-direction: column; align-items: flex-start; }
}
</style>
