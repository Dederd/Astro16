<template>
  <div class="step-container fade-in">
    <div class="step-header">
      <h1>Apa Momen Spesialmu? 🎉</h1>
      <p>Pilih jenis acara untuk mendapatkan rekomendasi bouquet yang tepat dari AI kami</p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <span>Memuat pilihan momen...</span>
    </div>

    <!-- Error state -->
    <div v-else-if="loadError" class="loading-state">
      <span>⚠️</span>
      <p style="color: #c62828; margin-top: 8px;">Gagal memuat pilihan momen. Pastikan server backend berjalan, lalu refresh halaman.</p>
    </div>

    <!-- Bouquet type grid -->
    <div v-else class="bouquet-type-grid">
      <button
        v-for="type in bouquetTypes"
        :key="type.id"
        class="type-card"
        :class="{ selected: store.selectedBouquetType?.id === type.id }"
        :style="`--theme: ${type.theme}`"
        @click="selectType(type)"
      >
        <div class="type-icon">{{ type.icon }}</div>
        <div class="type-info">
          <span class="type-name">{{ type.name }}</span>
          <span class="type-desc">{{ type.description }}</span>
        </div>
        <div v-if="store.selectedBouquetType?.id === type.id" class="type-check">✓</div>
      </button>
    </div>

    <!-- Agent response -->
    <transition name="slide-up">
      <div v-if="agentLoading || store.agentMessage" class="agent-section">
        <div v-if="agentLoading" class="agent-loading">
          <div class="dots-loader">
            <span></span><span></span><span></span>
          </div>
          <p>AI sedang menganalisis pilihanmu...</p>
        </div>

        <div v-else-if="store.agentMessage" class="agent-bubble fade-in">
          <div class="agent-label">Agent AI 🤖</div>
          <p class="agent-msg">{{ store.agentMessage }}</p>
          <div v-if="store.agentTips" class="agent-tips">
            <span class="tips-label">💡 Tips</span>
            <p>{{ store.agentTips }}</p>
          </div>
        </div>
      </div>
    </transition>

    <!-- Next button -->
    <div class="step-actions">
      <button
        class="btn btn-primary btn-lg"
        :disabled="!store.selectedBouquetType || agentLoading"
        @click="goNext"
      >
        Lanjut Pilih Bunga →
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useOrderStore } from '@/stores/order'
import { getBouquetTypes, agentVerifySelection } from '@/services/api'

const store = useOrderStore()
const bouquetTypes = ref([])
const loading = ref(true)
const agentLoading = ref(false)
const loadError = ref(false)

onMounted(async () => {
  try {
    const res = await getBouquetTypes()
    bouquetTypes.value = res.data.data || []
  } catch (e) {
    console.error('[Step1] Gagal load bouquet types:', e)
    loadError.value = true
  } finally {
    loading.value = false
  }
})

async function selectType(type) {
  store.selectBouquetType(type)
  store.setAgentResponse('', '', [])
  agentLoading.value = true

  try {
    const res = await agentVerifySelection({
      bouquet_type_id: type.id,
      bouquet_type: type.name
    })
    const d = res.data.data
    store.setAgentResponse(d.message, d.tips, d.recommendations)
  } catch (e) {
    console.error(e)
    store.setAgentResponse(
      `Pilihan ${type.name} yang indah! Kami memiliki banyak bunga cantik untuk acara ini.`,
      'Pilih bunga sesuai selera dan maknanya.',
      []
    )
  } finally {
    agentLoading.value = false
  }
}

function goNext() {
  if (store.selectedBouquetType) store.setStep(2)
}
</script>

<style scoped>
.step-container { max-width: 800px; margin: 0 auto; }

.step-header {
  text-align: center;
  margin-bottom: 48px;
}

.step-header h1 {
  font-size: clamp(1.8rem, 4vw, 2.5rem);
  color: var(--charcoal);
  margin-bottom: 12px;
}

.step-header p {
  font-size: 1rem;
  color: var(--warm-gray);
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 60px;
  color: var(--warm-gray);
}

.bouquet-type-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
  margin-bottom: 32px;
}

.type-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px 24px;
  background: var(--white);
  border: 2px solid var(--light-gray);
  border-radius: var(--radius);
  cursor: pointer;
  transition: var(--transition);
  text-align: left;
  position: relative;
  overflow: hidden;
}

.type-card::before {
  content: '';
  position: absolute;
  left: 0; top: 0; bottom: 0;
  width: 4px;
  background: var(--theme, var(--rose));
  opacity: 0;
  transition: var(--transition);
}

.type-card:hover {
  border-color: var(--theme, var(--rose));
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.type-card:hover::before { opacity: 1; }

.type-card.selected {
  border-color: var(--theme, var(--rose));
  background: linear-gradient(135deg, white 80%, rgba(244,141,176,0.1) 100%);
}

.type-card.selected::before { opacity: 1; }

.type-icon { font-size: 2rem; flex-shrink: 0; }

.type-info { flex: 1; }

.type-name {
  display: block;
  font-weight: 500;
  font-size: 1rem;
  color: var(--charcoal);
  margin-bottom: 2px;
}

.type-desc {
  display: block;
  font-size: 0.8rem;
  color: var(--warm-gray);
}

.type-check {
  width: 28px; height: 28px;
  background: var(--deep-rose);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8rem;
  font-weight: 700;
}

/* Agent section */
.agent-section { margin: 24px 0; }

.agent-loading {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 20px 24px;
  background: var(--white);
  border-radius: var(--radius);
  border: 1px solid var(--blush);
  color: var(--warm-gray);
}

.dots-loader {
  display: flex;
  gap: 5px;
}

.dots-loader span {
  width: 8px; height: 8px;
  background: var(--rose);
  border-radius: 50%;
  animation: bounce 1.2s ease infinite;
}
.dots-loader span:nth-child(2) { animation-delay: 0.2s; }
.dots-loader span:nth-child(3) { animation-delay: 0.4s; }

@keyframes bounce {
  0%, 80%, 100% { transform: translateY(0); }
  40% { transform: translateY(-8px); }
}

.agent-bubble { padding-top: 28px; }

.agent-label {
  font-size: 0.75rem;
  font-weight: 600;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--rose);
  margin-bottom: 8px;
}

.agent-msg {
  font-size: 0.95rem;
  color: var(--charcoal);
  line-height: 1.6;
  margin-bottom: 12px;
}

.agent-tips {
  background: rgba(138,158,140,0.12);
  border-radius: var(--radius-sm);
  padding: 12px 16px;
}

.tips-label {
  font-size: 0.78rem;
  font-weight: 600;
  color: var(--sage);
  display: block;
  margin-bottom: 4px;
}

.agent-tips p { font-size: 0.88rem; color: var(--warm-gray); }

.step-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 32px;
}

.btn-lg { padding: 16px 36px; font-size: 1rem; }

.slide-up-enter-active, .slide-up-leave-active {
  transition: all 0.4s ease;
}
.slide-up-enter-from { opacity: 0; transform: translateY(20px); }
.slide-up-leave-to { opacity: 0; transform: translateY(-10px); }
</style>
