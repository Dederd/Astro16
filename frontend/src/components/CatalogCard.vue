<template>
  <div
    class="catalog-card"
    :class="{ selected: isSelected, 'out-of-stock': !item.is_available }"
    @click="item.is_available && $emit('select', item)"
  >
    <!-- Image -->
    <div class="catalog-img-wrap" :style="`background: ${gradientFor(item.style)}`">
      <img
        v-if="item.image_url"
        :src="item.image_url"
        :alt="item.name"
        class="catalog-img"
        @error="imgError = true"
      />
      <div v-if="!item.image_url || imgError" class="catalog-emoji">
        {{ emojiFor(item.style) }}
      </div>
      <div v-if="!item.is_available" class="out-badge">Stok Habis</div>
      <div v-if="isSelected" class="selected-badge">✓ Dipilih</div>
      <div class="occasion-tag">{{ occasionLabel }}</div>
    </div>

    <div class="catalog-info">
      <h3 class="catalog-name">{{ item.name }}</h3>
      <p class="catalog-desc">{{ item.description }}</p>

      <div class="catalog-meta">
        <div class="meta-item">
          <span class="meta-icon">🌸</span>
          <span>{{ item.stem_count }} tangkai</span>
        </div>
        <div class="meta-item">
          <span class="meta-icon">🎨</span>
          <span>{{ item.style }}</span>
        </div>
      </div>

      <div class="catalog-footer">
        <span class="catalog-price">Rp{{ formatPrice(item.price) }}</span>
        <button
          v-if="!isSelected"
          class="btn btn-outline btn-sm"
          :disabled="!item.is_available"
          @click.stop="$emit('select', item)"
        >
          Pilih
        </button>
        <button v-else class="btn btn-primary btn-sm" disabled>✓ Dipilih</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  item: { type: Object, required: true },
  isSelected: { type: Boolean, default: false },
})

defineEmits(['select'])

const imgError = ref(false)

const occasionMap = {
  graduation: '🎓 Wisuda',
  wedding: '💍 Pernikahan',
  birthday: '🎂 Ulang Tahun',
  valentines: '❤️ Valentine',
  mothers_day: '🌸 Hari Ibu',
  anniversary: '💑 Anniversary',
  congratulations: '🏆 Selamat',
  sympathy: '🕊️ Dukacita',
}

const occasionLabel = computed(() => occasionMap[props.item.occasion] || props.item.occasion)

function gradientFor(style) {
  const g = {
    Classic:  'linear-gradient(135deg, #FFF0E6, #F0C8A0)',
    Premium:  'linear-gradient(135deg, #F5E6FF, #C485FF)',
    Romantic: 'linear-gradient(135deg, #FFECD2, #FCB69F)',
    Playful:  'linear-gradient(135deg, #FFE0F0, #FF8CC8)',
    Elegant:  'linear-gradient(135deg, #E8D5C4, #C4A882)',
    Warm:     'linear-gradient(135deg, #FFE0E0, #FF9F9F)',
  }
  return g[style] || g.Classic
}

function emojiFor(style) {
  const e = { Classic: '🌸', Premium: '👑', Romantic: '💐', Playful: '🌻', Elegant: '🌹', Warm: '❤️' }
  return e[style] || '💐'
}

function formatPrice(p) { return (p || 0).toLocaleString('id-ID') }
</script>

<style scoped>
.catalog-card {
  background: var(--white);
  border-radius: var(--radius);
  border: 1.5px solid var(--light-gray);
  overflow: hidden;
  cursor: pointer;
  transition: var(--transition);
  display: flex;
  flex-direction: column;
}
.catalog-card:hover:not(.out-of-stock) {
  border-color: var(--rose);
  box-shadow: 0 6px 24px rgba(198, 40, 40, 0.1);
  transform: translateY(-3px);
}
.catalog-card.selected {
  border-color: var(--deep-rose);
  box-shadow: 0 6px 28px rgba(139, 0, 0, 0.18);
}
.catalog-card.out-of-stock { opacity: 0.65; cursor: not-allowed; }

.catalog-img-wrap {
  position: relative;
  height: 180px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}
.catalog-img { width: 100%; height: 100%; object-fit: cover; }
.catalog-emoji { font-size: 4.5rem; }

.out-badge, .selected-badge {
  position: absolute;
  top: 10px;
  right: 10px;
  padding: 4px 12px;
  border-radius: var(--radius-pill);
  font-size: 0.75rem;
  font-weight: 600;
}
.out-badge { background: rgba(44,44,44,0.75); color: white; }
.selected-badge { background: var(--deep-rose); color: white; }

.occasion-tag {
  position: absolute;
  bottom: 10px;
  left: 10px;
  background: rgba(255,255,255,0.88);
  color: var(--charcoal);
  font-size: 0.72rem;
  font-weight: 500;
  padding: 3px 10px;
  border-radius: var(--radius-pill);
}

.catalog-info { padding: 18px; flex: 1; display: flex; flex-direction: column; }
.catalog-name { font-size: 1rem; font-weight: 600; color: var(--charcoal); margin-bottom: 6px; }
.catalog-desc { font-size: 0.8rem; color: var(--warm-gray); line-height: 1.5; margin-bottom: 12px; flex: 1; }

.catalog-meta { display: flex; gap: 12px; margin-bottom: 14px; }
.meta-item { display: flex; align-items: center; gap: 4px; font-size: 0.78rem; color: var(--warm-gray); }

.catalog-footer { display: flex; align-items: center; justify-content: space-between; }
.catalog-price { font-family: var(--font-display); font-size: 1.15rem; color: var(--deep-rose); font-weight: 500; }
.btn-sm { padding: 8px 16px; font-size: 0.82rem; }
</style>
