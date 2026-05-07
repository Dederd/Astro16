<template>
  <div class="step-container fade-in">
    <div class="step-header">
      <h1>Desain Bouquetmu ✨</h1>
      <p>AI kami akan merancang pilihan bouquet berdasarkan bunga-bunga yang kamu pilih</p>
    </div>

    <!-- Selected flowers recap -->
    <div class="flowers-recap card">
      <div class="recap-header">
        <span class="recap-title">Bunga Pilihanmu</span>
        <button class="btn btn-ghost btn-sm" @click="store.setStep(2)">Edit</button>
      </div>
      <div class="recap-flowers">
        <span v-for="f in store.selectedFlowers" :key="f.flower_id" class="recap-tag">
          🌸 {{ f.name }} × {{ f.quantity }}
        </span>
        <span class="recap-total">Total: {{ store.flowerCount }} tangkai</span>
      </div>
    </div>

    <!-- Generate limit warning + Buy Quota -->
    <div v-if="store.isGenerateLimited && !generated" class="limit-banner">
      <span>🔒</span>
      <div class="limit-content">
        <strong>Kuota generate habis ({{ store.generateCount }}/{{ store.generateLimit }}x)</strong>
        <p>Beli paket kuota tambahan untuk terus merancang bouquet impianmu.</p>
        <div class="buy-quota-box">
          <div class="quota-pack-info">
            <span class="quota-pack-label">✨ Paket 3 Generate</span>
            <span class="quota-pack-price">Rp5.000</span>
          </div>
          <button class="btn btn-primary btn-sm buy-quota-btn" @click="buyQuota" :disabled="buyingQuota">
            <span v-if="buyingQuota">⏳ Memproses...</span>
            <span v-else>🛒 Beli Sekarang</span>
          </button>
        </div>
        <p class="quota-note">* Biaya akan ditambahkan ke total pesananmu</p>
      </div>
    </div>

    <!-- Generate quota info -->
    <div v-else-if="!generated && !loading" class="quota-info">
      <span class="quota-badge">{{ store.generateCount }}/{{ store.generateLimit }} generate gratis terpakai</span>
    </div>

    <!-- Generate zone (before generation) -->
    <div v-if="!generated && !loading && !store.isGenerateLimited" class="generate-zone">
      <div class="generate-illustration">
        <div class="gen-flower">💐</div>
        <div class="gen-sparkles">✨ 🤖 ✨</div>
      </div>
      <h2>Siap Generate Desain?</h2>
      <p>Klik tombol di bawah untuk membiarkan AI kami menciptakan desain bouquet yang sempurna untukmu</p>
      <button class="btn btn-primary btn-xl" @click="generate">
        <span>🪄 Generate Bouquet Sekarang</span>
      </button>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="ai-loading">
      <div class="ai-thinking">
        <div class="thinking-flower">💐</div>
        <div class="thinking-dots"><span></span><span></span><span></span></div>
      </div>
      <h3>AI sedang merancang bouquetmu...</h3>
      <p>{{ loadingMsg }}</p>
    </div>

    <!-- Error state -->
    <div v-if="error" class="error-state">
      <span>⚠️</span>
      <p>{{ error }}</p>
      <button class="btn btn-outline" @click="generate" v-if="!store.isGenerateLimited">Coba Lagi</button>
    </div>

    <!-- Generated designs -->
    <div v-if="generated && !loading">
      <div v-if="store.designMessage" class="agent-bubble" style="margin-bottom: 32px;">
        <div class="agent-label">Agent AI 🤖</div>
        <p>{{ store.designMessage }}</p>
      </div>

      <div class="designs-label">
        <h2>Pilih Desain Favoritmu</h2>
        <p>Setiap desain tersedia dalam 2 ukuran — pilih yang paling sesuai</p>
      </div>

      <div class="designs-grid">
        <DesignCard
          v-for="design in store.generatedDesigns"
          :key="design.id"
          :design="design"
          :is-selected="store.selectedDesign?.id === design.id"
          :selected-size="store.selectedSize"
          @select="(d, size) => store.selectDesign(d, size)"
        />
      </div>

      <!-- Regenerate button (hanya jika belum limit) -->
      <div class="regen-wrap">
        <button
          class="btn btn-ghost"
          @click="generate"
          :disabled="store.isGenerateLimited"
          :title="store.isGenerateLimited ? 'Batas generate gratis tercapai' : ''"
        >
          🔄 Generate Ulang
          <span v-if="!store.isGenerateLimited" class="quota-small">
            (sisa {{ store.generateLimit - store.generateCount }}x)
          </span>
          <span v-else class="quota-small quota-exhausted">🔒 Terbatas</span>
        </button>
      </div>
    </div>

    <!-- Navigation -->
    <div class="step-actions">
      <button class="btn btn-ghost" @click="store.setStep(2)">← Kembali</button>
      <button
        v-if="store.selectedDesign"
        class="btn btn-primary btn-lg"
        @click="store.setStep(4)"
      >
        Lanjut ke Pembayaran →
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useOrderStore } from '@/stores/order'
import { agentGenerateBouquet, getGenerateStatus, buyGenerateQuota, confirmQuotaPayment } from '@/services/api'
import DesignCard from '@/components/DesignCard.vue'

const store = useOrderStore()
const loading = ref(false)
const generated = ref(false)
const error = ref('')
const buyingQuota = ref(false)
const loadingMsgs = [
  'Menganalisis kombinasi bungamu...',
  'Menciptakan palet warna yang harmonis...',
  'Merancang komposisi bouquet yang sempurna...',
  'Menyentuh setiap detail dengan penuh cinta...',
]
const loadingMsg = ref(loadingMsgs[0])
let msgInterval = null

onMounted(async () => {
  // Cek sisa kuota generate
  try {
    const res = await getGenerateStatus()
    store.setGeneratedDesigns(
      store.generatedDesigns,
      store.designMessage,
      res.data.generate_count,
      res.data.limit,
    )
    if (res.data.extra_quota !== undefined) {
      store.setExtraQuota(res.data.extra_quota, 0)
    }
  } catch { /* ignore */ }

  if (store.generatedDesigns.length > 0) {
    generated.value = true
  }
})

async function buyQuota() {
  buyingQuota.value = true
  try {
    const res = await buyGenerateQuota()
    const { snap_token, order_id, session_id } = res.data

    if (!snap_token) {
      alert('Gagal mendapatkan token pembayaran.')
      return
    }

    // Buka Midtrans Snap popup
    window.snap.pay(snap_token, {
      onSuccess: async () => {
        try {
          const confirmRes = await confirmQuotaPayment({ order_id, session_id })
          const d = confirmRes.data
          store.setExtraQuota(d.extra_quota, d.extra_quota_fee)
          store.setGeneratedDesigns(
            store.generatedDesigns,
            store.designMessage,
            d.generate_count,
            d.limit,
          )
          // Auto-generate setelah kuota berhasil dibeli
          await generate()
        } catch (e) {
          alert('Pembayaran berhasil, namun gagal menambah kuota. Hubungi support.')
        }
      },
      onPending: () => {
        alert('Pembayaran sedang diproses. Kuota akan ditambahkan setelah pembayaran dikonfirmasi.')
      },
      onError: () => {
        alert('Pembayaran gagal. Silakan coba lagi.')
      },
      onClose: () => {
        // User menutup popup tanpa bayar
        buyingQuota.value = false
      },
    })
  } catch (e) {
    alert('Gagal membuat token pembayaran: ' + (e?.response?.data?.error || e.message))
  } finally {
    buyingQuota.value = false
  }
}

async function generate() {
  if (store.isGenerateLimited) return

  loading.value = true
  generated.value = false
  error.value = ''
  store.selectDesign(null, null)

  let msgIdx = 0
  msgInterval = setInterval(() => {
    msgIdx = (msgIdx + 1) % loadingMsgs.length
    loadingMsg.value = loadingMsgs[msgIdx]
  }, 2000)

  try {
    // BUG FIX #2: kirim total_stem_count supaya AI output sinkron dengan pilihan user
    const totalStemCount = store.selectedFlowers.reduce((sum, f) => sum + f.quantity, 0)

    const res = await agentGenerateBouquet({
      bouquet_type_id: store.selectedBouquetType.id,
      selected_flowers: store.selectedFlowers,
      total_stem_count: totalStemCount,
      style_hint: store.aiStyleHint || undefined,
      description_hint: store.aiDescriptionHint || undefined,
    })

    const { designs, message } = res.data.data
    const generateCount = res.data.generate_count
    const limit = res.data.limit
    const aiFee = res.data.ai_fee || 0

    store.setGeneratedDesigns(designs, message, generateCount, limit)
    // Tandai apakah generate ini berbayar (AI fee dikenakan)
    store.setAIFeePaid(aiFee > 0)
    generated.value = true
  } catch (e) {
    const errData = e?.response?.data
    if (e?.response?.status === 429) {
      error.value = errData?.error || 'Batas generate gratis tercapai.'
      store.setGeneratedDesigns(store.generatedDesigns, store.designMessage, errData?.generate_count, errData?.limit)
    } else {
      error.value = 'Maaf, terjadi kesalahan saat generate desain. Silakan coba lagi.'
    }
    console.error(e)
  } finally {
    loading.value = false
    clearInterval(msgInterval)
  }
}
</script>

<style scoped>
.step-container { max-width: 960px; margin: 0 auto; }
.step-header { text-align: center; margin-bottom: 32px; }
.step-header h1 { font-size: clamp(1.8rem, 4vw, 2.5rem); margin-bottom: 10px; }
.step-header p { color: var(--warm-gray); }

.flowers-recap { padding: 20px 24px; margin-bottom: 24px; }
.recap-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px; }
.recap-title { font-size: 0.82rem; font-weight: 600; letter-spacing: 0.08em; text-transform: uppercase; color: var(--warm-gray); }
.recap-flowers { display: flex; flex-wrap: wrap; gap: 8px; align-items: center; }
.recap-tag { background: var(--cream); border: 1px solid var(--blush); padding: 5px 12px; border-radius: var(--radius-pill); font-size: 0.83rem; color: var(--charcoal); }
.recap-total { font-weight: 600; font-size: 0.85rem; color: var(--deep-rose); margin-left: 4px; }

.limit-banner {
  display: flex;
  gap: 16px;
  align-items: flex-start;
  background: #FFF3E0;
  border: 1px solid #FFB74D;
  border-radius: var(--radius);
  padding: 16px 20px;
  margin-bottom: 24px;
  font-size: 0.9rem;
}
.limit-banner span { font-size: 1.5rem; flex-shrink: 0; }
.limit-banner p { color: var(--warm-gray); margin: 4px 0 0; font-size: 0.85rem; }
.limit-content { flex: 1; }
.limit-content strong { display: block; margin-bottom: 4px; color: var(--charcoal); }

.buy-quota-box {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: white;
  border: 1.5px solid #FFB74D;
  border-radius: var(--radius-sm);
  padding: 10px 14px;
  margin: 12px 0 6px;
  gap: 12px;
}
.quota-pack-info { display: flex; flex-direction: column; gap: 2px; }
.quota-pack-label { font-size: 0.82rem; font-weight: 600; color: var(--charcoal); }
.quota-pack-price { font-size: 1rem; font-weight: 700; color: var(--deep-rose); }
.buy-quota-btn { white-space: nowrap; flex-shrink: 0; }
.quota-note { font-size: 0.75rem; color: var(--warm-gray); margin-top: 4px !important; }

.quota-info { text-align: center; margin-bottom: 16px; }
.quota-badge {
  display: inline-block;
  background: var(--cream);
  border: 1px solid var(--blush);
  color: var(--warm-gray);
  padding: 4px 14px;
  border-radius: var(--radius-pill);
  font-size: 0.8rem;
}

.generate-zone {
  text-align: center;
  padding: 60px 32px;
  background: linear-gradient(135deg, #FFF8F5 0%, #FDF0EC 100%);
  border-radius: 24px;
  border: 1px dashed var(--blush);
  margin-bottom: 32px;
}
.generate-illustration { margin-bottom: 24px; }
.gen-flower { font-size: 5rem; margin-bottom: 8px; }
.gen-sparkles { font-size: 1.5rem; letter-spacing: 8px; }
.generate-zone h2 { font-size: 1.8rem; margin-bottom: 12px; color: var(--charcoal); }
.generate-zone p { font-size: 0.95rem; color: var(--warm-gray); max-width: 480px; margin: 0 auto 32px; line-height: 1.6; }
.btn-xl { padding: 18px 44px; font-size: 1.05rem; }

.ai-loading { text-align: center; padding: 80px 32px; }
.ai-thinking { display: flex; flex-direction: column; align-items: center; gap: 16px; margin-bottom: 24px; }
.thinking-flower { font-size: 4rem; animation: float 2s ease-in-out infinite; }
@keyframes float { 0%, 100% { transform: translateY(0); } 50% { transform: translateY(-12px); } }
.thinking-dots { display: flex; gap: 8px; }
.thinking-dots span { width: 10px; height: 10px; background: var(--rose); border-radius: 50%; animation: bounce 1.2s ease infinite; }
.thinking-dots span:nth-child(2) { animation-delay: 0.2s; }
.thinking-dots span:nth-child(3) { animation-delay: 0.4s; }
@keyframes bounce { 0%, 80%, 100% { transform: translateY(0); } 40% { transform: translateY(-10px); } }
.ai-loading h3 { font-size: 1.4rem; color: var(--charcoal); margin-bottom: 8px; }
.ai-loading p { color: var(--warm-gray); font-size: 0.9rem; }

.error-state { text-align: center; padding: 48px; background: #FFF5F5; border-radius: var(--radius); border: 1px solid #FFCDD2; }
.error-state span { font-size: 2rem; display: block; margin-bottom: 12px; }
.error-state p { color: var(--warm-gray); margin-bottom: 20px; }

.designs-label { text-align: center; margin-bottom: 32px; }
.designs-label h2 { font-size: 1.8rem; margin-bottom: 8px; }
.designs-label p { color: var(--warm-gray); font-size: 0.9rem; }

.designs-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: 24px; margin-bottom: 32px; }

.regen-wrap { text-align: center; margin-bottom: 32px; }
.quota-small { font-size: 0.75rem; color: var(--warm-gray); margin-left: 4px; }
.quota-exhausted { color: #F57C00; }

.agent-bubble { padding-top: 28px; }
.agent-label { font-size: 0.75rem; font-weight: 600; letter-spacing: 0.08em; text-transform: uppercase; color: var(--rose); margin-bottom: 8px; }

.step-actions { display: flex; justify-content: space-between; align-items: center; margin-top: 32px; padding-top: 24px; border-top: 1px solid var(--light-gray); }
.btn-sm { padding: 8px 16px; font-size: 0.82rem; }
.btn-lg { padding: 16px 36px; font-size: 1rem; }
</style>
