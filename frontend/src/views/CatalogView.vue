<template>
  <div class="catalog-page fade-in">
    <div class="catalog-hero">
      <h1>🛍️ Katalog Bouquet</h1>
      <p>Pilih bouquet siap pakai dari koleksi kami, atau buat desainmu sendiri bersama AI</p>
      <div class="hero-cta">
        <button class="btn btn-primary" @click="goAI">
          🤖 Buat dengan AI
        </button>
      </div>
    </div>

    <!-- Filter bar -->
    <div class="filter-section">
      <div class="filter-label">Filter Acara:</div>
      <div class="filter-tabs">
        <button
          v-for="occ in occasions"
          :key="occ.value"
          class="filter-tab"
          :class="{ active: activeOccasion === occ.value }"
          @click="activeOccasion = occ.value"
        >{{ occ.label }}</button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="catalog-grid">
      <div v-for="i in 6" :key="i" class="skeleton" style="height: 380px; border-radius: 16px;"></div>
    </div>

    <!-- Grid -->
    <div v-else class="catalog-grid">
      <CatalogCard
        v-for="item in displayItems"
        :key="item.id"
        :item="item"
        :is-selected="store.selectedCatalogItem?.id === item.id && store.orderMode === 'catalog'"
        @select="selectItem"
      />
    </div>

    <!-- Error state -->
    <div v-if="loadError" class="empty-state">
      <span>⚠️</span>
      <p>Gagal memuat katalog. Pastikan server backend berjalan dan coba refresh halaman.</p>
    </div>

    <!-- Empty -->
    <div v-if="!loading && displayItems.length === 0" class="empty-state">
      <span>🌿</span>
      <p>Belum ada katalog untuk acara ini</p>
    </div>

    <!-- Bottom sheet saat ada pilihan -->
    <transition name="slide-up">
      <div v-if="store.selectedCatalogItem && store.orderMode === 'catalog'" class="selection-bar">
        <div class="sel-info">
          <span class="sel-name">{{ store.selectedCatalogItem.name }}</span>
          <span class="sel-price">Rp{{ formatPrice(store.selectedCatalogItem.price) }}</span>
        </div>
        <div class="sel-actions">
          <button class="btn btn-ghost btn-sm" @click="store.setOrderMode('ai')">Batal</button>
          <button class="btn btn-primary" @click="checkoutCatalog">Checkout →</button>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useOrderStore } from '@/stores/order'
import { getCatalog } from '@/services/api'
import CatalogCard from '@/components/CatalogCard.vue'

const store = useOrderStore()
const router = useRouter()
const allItems = ref([])
const loading = ref(true)
const activeOccasion = ref('all')
const loadError = ref(false)

const occasions = [
  { label: 'Semua', value: 'all' },
  { label: '🎓 Wisuda', value: 'graduation' },
  { label: '💍 Pernikahan', value: 'wedding' },
  { label: '🎂 Ulang Tahun', value: 'birthday' },
  { label: '❤️ Valentine', value: 'valentines' },
  { label: '🌸 Hari Ibu', value: 'mothers_day' },
]

onMounted(async () => {
  try {
    const res = await getCatalog()
    allItems.value = res.data.data || []
  } catch (e) {
    console.error(e)
    loadError.value = true
  } finally {
    loading.value = false
  }
})

const displayItems = computed(() => {
  if (activeOccasion.value === 'all') return allItems.value
  return allItems.value.filter(i => i.occasion === activeOccasion.value)
})

function selectItem(item) {
  store.selectCatalogItem(item)
}

function checkoutCatalog() {
  // Untuk catalog, kita skip step 2 & 3 — langsung step 4 (checkout)
  store.setStep(4)
  router.push('/order')
}

function goAI() {
  store.setOrderMode('ai')
  store.reset()
  router.push('/order')
}

function formatPrice(p) { return (p || 0).toLocaleString('id-ID') }
</script>

<style scoped>
.catalog-page { max-width: 1080px; margin: 0 auto; padding-bottom: 120px; }

.catalog-hero {
  text-align: center;
  padding: 48px 24px 40px;
}
.catalog-hero h1 { font-size: clamp(2rem, 5vw, 3rem); margin-bottom: 14px; }
.catalog-hero p { color: var(--warm-gray); font-size: 1rem; max-width: 540px; margin: 0 auto 28px; line-height: 1.6; }

.hero-cta { display: flex; justify-content: center; gap: 12px; flex-wrap: wrap; }

.filter-section {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 32px;
  flex-wrap: wrap;
}
.filter-label { font-size: 0.85rem; font-weight: 500; color: var(--warm-gray); white-space: nowrap; }
.filter-tabs { display: flex; gap: 8px; flex-wrap: wrap; }
.filter-tab {
  padding: 8px 16px;
  border-radius: var(--radius-pill);
  border: 1.5px solid var(--light-gray);
  background: var(--white);
  font-size: 0.82rem;
  color: var(--warm-gray);
  cursor: pointer;
  transition: var(--transition);
}
.filter-tab:hover { border-color: var(--rose); color: var(--rose); }
.filter-tab.active { background: var(--deep-rose); border-color: var(--deep-rose); color: white; }

.catalog-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 24px;
}

.empty-state { text-align: center; padding: 80px; color: var(--warm-gray); }
.empty-state span { font-size: 3rem; display: block; margin-bottom: 16px; }

/* Bottom selection bar */
.selection-bar {
  position: fixed;
  bottom: 0; left: 0; right: 0;
  background: var(--white);
  border-top: 2px solid var(--blush);
  box-shadow: 0 -6px 32px rgba(0,0,0,0.1);
  padding: 16px 32px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  z-index: 50;
  flex-wrap: wrap;
}
.sel-info { display: flex; flex-direction: column; gap: 2px; }
.sel-name { font-weight: 500; font-size: 0.95rem; color: var(--charcoal); }
.sel-price { font-family: var(--font-display); font-size: 1.2rem; color: var(--deep-rose); }
.sel-actions { display: flex; gap: 10px; align-items: center; }
.btn-sm { padding: 8px 16px; font-size: 0.82rem; }

.slide-up-enter-active, .slide-up-leave-active { transition: all 0.3s ease; }
.slide-up-enter-from, .slide-up-leave-to { transform: translateY(100%); opacity: 0; }
</style>
