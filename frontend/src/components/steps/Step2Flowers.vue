<template>
  <div class="step-container fade-in">
    <div class="step-header">
      <h1>Pilih Bungamu 🌺</h1>
      <p>Pilih satu atau lebih jenis bunga untuk bouquetmu. Bunga bertanda ★ direkomendasikan AI.</p>
    </div>

    <!-- Recommendations banner -->
    <div v-if="store.recommendedFlowerIds.length > 0" class="rec-banner">
      <span>✨ AI merekomendasikan {{ store.recommendedFlowerIds.length }} bunga untuk {{ store.selectedBouquetType?.name }}</span>
      <button class="btn btn-ghost btn-sm" @click="toggleShowAll">
        {{ showAll ? 'Tampilkan Rekomendasi' : 'Tampilkan Semua' }}
      </button>
    </div>

    <!-- Search & filter -->
    <div class="filter-bar">
      <div class="search-wrap">
        <span class="search-icon">🔍</span>
        <input v-model="search" type="text" placeholder="Cari bunga..." class="form-input search-input" />
      </div>
      <div class="filter-tabs">
        <button
          v-for="filter in filters"
          :key="filter.value"
          class="filter-tab"
          :class="{ active: activeFilter === filter.value }"
          @click="activeFilter = filter.value"
        >{{ filter.label }}</button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="loading-grid">
      <div v-for="i in 8" :key="i" class="skeleton" style="height: 300px; border-radius: 16px;"></div>
    </div>

    <!-- Error state -->
    <div v-else-if="loadError" class="empty-state">
      <span>⚠️</span>
      <p style="color: #c62828;">Gagal memuat data bunga. Pastikan server backend berjalan, lalu refresh halaman.</p>
    </div>

    <!-- Flowers grid -->
    <!-- BUG FIX #1: margin-bottom besar supaya konten tidak tertutup sticky bar -->
    <div v-else class="flowers-grid">
      <FlowerCard
        v-for="flower in displayFlowers"
        :key="flower.id"
        :flower="flower"
        :is-selected="store.isFlowerSelected(flower.id)"
        :is-recommended="store.recommendedFlowerIds.includes(flower.id)"
        :quantity="store.getFlowerQuantity(flower.id)"
        @toggle="store.toggleFlower"
        @update-qty="store.updateFlowerQuantity"
      />
    </div>

    <!-- Empty state -->
    <div v-if="!loading && displayFlowers.length === 0" class="empty-state">
      <span>🌿</span>
      <p>Tidak ada bunga yang cocok dengan pencarianmu</p>
    </div>

    <!-- BUG FIX #1: Sticky summary bar — z-index rendah, TIDAK menutup qty control di card -->
    <!-- Pakai position sticky bukan fixed supaya tidak overlap konten -->
    <transition name="slide-up">
      <div v-if="store.selectedFlowers.length > 0" class="selection-summary">
        <div class="summary-info">
          <span class="summary-count">{{ store.flowerCount }} tangkai dipilih</span>
          <div class="summary-flowers">
            <span v-for="f in store.selectedFlowers" :key="f.flower_id" class="summary-tag">
              {{ f.name }} × {{ f.quantity }}
            </span>
          </div>
        </div>

        <!-- Optional style & description inputs -->
        <div class="ai-prompt-inputs">
          <div class="prompt-input-row">
            <label class="prompt-label">🎨 Gaya Desain <span class="optional-tag">opsional</span></label>
            <input
              v-model="store.aiStyleHint"
              class="prompt-input"
              placeholder="Contoh: Romantic, Minimalist, Rustic, Elegant..."
              maxlength="100"
            />
          </div>
          <div class="prompt-input-row">
            <label class="prompt-label">📝 Deskripsi Tambahan <span class="optional-tag">opsional</span></label>
            <textarea
              v-model="store.aiDescriptionHint"
              class="prompt-textarea"
              placeholder="Contoh: Untuk ulang tahun ibu, suasana hangat dan penuh kasih sayang..."
              maxlength="300"
              rows="2"
            ></textarea>
          </div>
        </div>

        <div class="summary-actions">
          <button class="btn btn-ghost btn-sm" @click="store.setStep(1)">← Kembali</button>
          <button class="btn btn-primary" @click="store.setStep(3)">Generate Desain →</button>
        </div>
      </div>
    </transition>

    <!-- Bottom nav if nothing selected -->
    <div v-if="store.selectedFlowers.length === 0" class="step-actions">
      <button class="btn btn-ghost" @click="store.setStep(1)">← Kembali</button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useOrderStore } from '@/stores/order'
import { getFlowers } from '@/services/api'
import FlowerCard from '@/components/FlowerCard.vue'

const store = useOrderStore()
const allFlowers = ref([])
const loading = ref(true)
const search = ref('')
const activeFilter = ref('all')
const showAll = ref(false)
const loadError = ref(false)

const filters = [
  { label: 'Semua', value: 'all' },
  { label: 'Ready', value: 'available' },
  { label: 'Rekomendasi AI', value: 'recommended' },
]

onMounted(async () => {
  try {
    const res = await getFlowers()
    allFlowers.value = res.data.data || []
    if (store.recommendedFlowerIds.length === 0) showAll.value = true
  } catch (e) {
    console.error('[Step2] Gagal load flowers:', e)
    loadError.value = true
  } finally {
    loading.value = false
  }
})

function toggleShowAll() {
  showAll.value = !showAll.value
  activeFilter.value = showAll.value ? 'all' : 'recommended'
}

const displayFlowers = computed(() => {
  let flowers = allFlowers.value

  if (!showAll.value && store.recommendedFlowerIds.length > 0) {
    flowers = flowers.filter(f => store.recommendedFlowerIds.includes(f.id))
  }
  if (activeFilter.value === 'available') {
    flowers = flowers.filter(f => f.is_available)
  } else if (activeFilter.value === 'recommended') {
    flowers = flowers.filter(f => store.recommendedFlowerIds.includes(f.id))
  }
  if (search.value) {
    const q = search.value.toLowerCase()
    flowers = flowers.filter(f =>
      f.name_id.toLowerCase().includes(q) ||
      f.name.toLowerCase().includes(q) ||
      f.meaning.toLowerCase().includes(q)
    )
  }
  return [...flowers].sort((a, b) => {
    const aRec = store.recommendedFlowerIds.includes(a.id)
    const bRec = store.recommendedFlowerIds.includes(b.id)
    if (aRec !== bRec) return bRec ? 1 : -1
    if (a.is_available !== b.is_available) return b.is_available ? 1 : -1
    return 0
  })
})
</script>

<style scoped>
.step-container { max-width: 1000px; margin: 0 auto; }

.step-header { text-align: center; margin-bottom: 32px; }
.step-header h1 { font-size: clamp(1.8rem, 4vw, 2.5rem); margin-bottom: 10px; }
.step-header p { color: var(--warm-gray); font-size: 0.95rem; }

.rec-banner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: linear-gradient(135deg, #FFF8F0, #FFF0E8);
  border: 1px solid var(--blush);
  border-radius: var(--radius-sm);
  padding: 12px 20px;
  margin-bottom: 20px;
  font-size: 0.88rem;
  color: var(--warm-gray);
  flex-wrap: wrap;
  gap: 8px;
}

.filter-bar {
  display: flex;
  gap: 16px;
  align-items: center;
  margin-bottom: 28px;
  flex-wrap: wrap;
}

.search-wrap { position: relative; flex: 1; min-width: 200px; }
.search-icon { position: absolute; left: 12px; top: 50%; transform: translateY(-50%); font-size: 0.9rem; }
.search-input { padding-left: 36px; }

.filter-tabs { display: flex; gap: 6px; }

.filter-tab {
  padding: 8px 16px;
  border-radius: var(--radius-pill);
  border: 1.5px solid var(--light-gray);
  background: var(--white);
  font-size: 0.82rem;
  color: var(--warm-gray);
  cursor: pointer;
  transition: var(--transition);
  white-space: nowrap;
}
.filter-tab:hover { border-color: var(--rose); color: var(--rose); }
.filter-tab.active { background: var(--deep-rose); border-color: var(--deep-rose); color: white; }

.loading-grid, .flowers-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 24px;
  /* BUG FIX #1: padding-bottom besar supaya kartu bagian bawah tidak
     tertutup oleh sticky summary bar */
  padding-bottom: 160px;
}

.empty-state { text-align: center; padding: 60px; color: var(--warm-gray); }
.empty-state span { font-size: 3rem; display: block; margin-bottom: 16px; }

/* BUG FIX #1: Sticky bar pakai position sticky + bottom:0 (bukan fixed)
   sehingga hanya "nempel" ketika discroll, tidak mengambang di atas kartu */
.selection-summary {
  position: sticky;
  bottom: 0;
  left: 0; right: 0;
  background: var(--white);
  border-top: 2px solid var(--blush);
  box-shadow: 0 -4px 24px rgba(0,0,0,0.08);
  padding: 14px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  /* z-index rendah — tidak boleh overlap qty-control di dalam kartu */
  z-index: 10;
  border-radius: var(--radius) var(--radius) 0 0;
  margin-top: 16px;
  flex-wrap: wrap;
}

.summary-info { flex: 1; }

.summary-count {
  display: block;
  font-weight: 500;
  color: var(--deep-rose);
  margin-bottom: 6px;
  font-size: 0.9rem;
}

.summary-flowers { display: flex; gap: 6px; flex-wrap: wrap; }

.summary-tag {
  background: var(--cream);
  border: 1px solid var(--blush);
  color: var(--charcoal);
  padding: 2px 10px;
  border-radius: var(--radius-pill);
  font-size: 0.75rem;
}

/* AI Prompt inputs */
.ai-prompt-inputs {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 14px 0 4px;
  border-top: 1px solid var(--blush);
  margin-top: 12px;
}

.prompt-input-row {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.prompt-label {
  font-size: 0.78rem;
  font-weight: 500;
  color: var(--charcoal);
  display: flex;
  align-items: center;
  gap: 6px;
}

.optional-tag {
  font-size: 0.68rem;
  color: var(--warm-gray);
  font-weight: 400;
  background: var(--light-gray);
  padding: 1px 7px;
  border-radius: var(--radius-pill);
}

.prompt-input,
.prompt-textarea {
  width: 100%;
  border: 1.5px solid var(--light-gray);
  border-radius: 8px;
  padding: 8px 12px;
  font-size: 0.83rem;
  color: var(--charcoal);
  background: var(--white);
  font-family: inherit;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.prompt-input:focus,
.prompt-textarea:focus {
  outline: none;
  border-color: var(--rose);
}

.prompt-textarea {
  resize: none;
  line-height: 1.5;
}

.summary-actions { display: flex; gap: 10px; align-items: center; }

.btn-sm { padding: 8px 16px; font-size: 0.82rem; }

.step-actions { display: flex; justify-content: flex-start; margin-top: 32px; }

.slide-up-enter-active, .slide-up-leave-active { transition: all 0.3s ease; }
.slide-up-enter-from { opacity: 0; transform: translateY(20px); }
.slide-up-leave-to { opacity: 0; transform: translateY(20px); }
</style>
