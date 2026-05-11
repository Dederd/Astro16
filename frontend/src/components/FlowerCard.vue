<template>
  <div
    class="flower-card"
    :class="{
      selected: isSelected,
      recommended: isRecommended,
      unavailable: !flower.is_available,
      disabled: disabled
    }"
    @click="flower.is_available && !disabled && $emit('toggle', flower)"
  >
    <!-- Recommended badge -->
    <div v-if="isRecommended" class="rec-badge">★ Recommended</div>

    <!-- Image area -->
    <div class="flower-img-wrap">
      <img
        v-if="flower.image_url && !imgError"
        :src="resolveImageUrl(flower.image_url)"
        :alt="flower.name_id"
        class="flower-img"
        @error="imgError = true"
      />
      <div v-if="!flower.image_url || imgError" class="flower-emoji-placeholder">
        {{ flower.emoji || flowerEmoji }}
      </div>
      <div v-if="!flower.is_available" class="unavailable-overlay">
        <span>Stok Habis</span>
      </div>
      <div v-if="isSelected" class="selected-check">✓</div>
      <!-- Bug fix 4: small zoom button instead of whole image click -->
      <button
        v-if="flower.image_url && !imgError"
        class="zoom-btn"
        @click.stop="openImageZoom"
        title="Lihat foto"
      >🔍</button>
    </div>

    <!-- Info -->
    <div class="flower-info">
      <div class="flower-header">
        <span class="flower-name">{{ flower.name_id }}</span>
        <span class="badge" :class="flower.is_available ? 'badge-available' : 'badge-unavailable'">
          {{ flower.is_available ? 'Ready' : 'Habis' }}
        </span>
      </div>
      <p class="flower-meaning">{{ flower.meaning }}</p>
      <div class="flower-colors">
        <span
          v-for="color in flower.colors"
          :key="color"
          class="color-dot"
          :style="`background: ${color}; border: 1.5px solid ${color === '#FAFAFA' || color === '#FFFFFF' ? '#ddd' : color}`"
        ></span>
      </div>
    </div>
  </div>

  <ImageZoomModal
    :is-open="showImageZoom"
    :image-url="resolveImageUrl(flower.image_url)"
    :image-alt="flower.name_id"
    @close="showImageZoom = false"
  />
</template>

<script setup>
import { ref, computed } from 'vue'
import ImageZoomModal from './ImageZoomModal.vue'

const props = defineProps({
  flower: { type: Object, required: true },
  isSelected: { type: Boolean, default: false },
  isRecommended: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
})

const emit = defineEmits(['toggle'])
const imgError = ref(false)
const showImageZoom = ref(false)

// Convert Google Drive share/view URL to direct thumbnail URL
function resolveImageUrl(url) {
  if (!url) return url
  const ucMatch = url.match(/drive\.google\.com\/uc\?.*[&?]id=([^&]+)/)
  if (ucMatch) return `https://drive.google.com/thumbnail?id=${ucMatch[1]}&sz=w400`
  const fileMatch = url.match(/drive\.google\.com\/file\/d\/([^/]+)/)
  if (fileMatch) return `https://drive.google.com/thumbnail?id=${fileMatch[1]}&sz=w400`
  return url
}

const flowerEmoji = computed(() => {
  const map = {
    rose_red: '🌹', rose_pink: '🌸', rose_white: '🤍',
    tulip_yellow: '🌷', tulip_purple: '💜',
    sunflower: '🌻', lily_white: '🌼',
    chrysanthemum: '💮', baby_breath: '🌬️',
    carnation_red: '🌺', carnation_pink: '🌸',
    orchid_purple: '🌸', lavender: '💜',
    daisy: '🌼', peony: '🌹'
  }
  return map[props.flower.id] || '🌸'
})

function formatPrice(p) {
  return p.toLocaleString('id-ID')
}

function openImageZoom() {
  if (props.flower.image_url && !imgError.value) {
    showImageZoom.value = true
  }
}
</script>

<style scoped>
.flower-card {
  background: var(--white);
  border-radius: var(--radius);
  border: 1.5px solid var(--light-gray);
  overflow: visible; /* penting: jangan hidden supaya qty control terlihat */
  cursor: pointer;
  transition: var(--transition);
  position: relative;
  user-select: none;
}

.flower-card:hover:not(.unavailable) {
  border-color: var(--rose);
  box-shadow: 0 4px 20px rgba(198, 40, 40, 0.1);
  transform: translateY(-2px);
}

.flower-card.selected {
  border-color: var(--deep-rose);
  box-shadow: 0 4px 20px rgba(139, 0, 0, 0.15);
}

.flower-card.recommended {
  border-color: var(--gold);
}

.rec-badge {
  position: absolute;
  top: -10px;
  left: 50%;
  transform: translateX(-50%);
  background: var(--gold);
  color: white;
  font-size: 0.7rem;
  font-weight: 600;
  padding: 3px 12px;
  border-radius: var(--radius-pill);
  white-space: nowrap;
  z-index: 2;
}

.flower-img-wrap {
  position: relative;
  aspect-ratio: 1;
  overflow: hidden;
  background: linear-gradient(135deg, var(--cream) 0%, var(--blush) 100%);
  border-radius: var(--radius) var(--radius) 0 0;
}

.flower-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.flower-emoji-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 3.5rem;
}

.unavailable-overlay {
  position: absolute;
  inset: 0;
  background: rgba(44,44,44,0.55);
  display: flex;
  align-items: center;
  justify-content: center;
}

.unavailable-overlay span {
  background: rgba(255,255,255,0.9);
  color: var(--charcoal);
  padding: 6px 14px;
  border-radius: var(--radius-pill);
  font-size: 0.8rem;
  font-weight: 500;
}

.selected-check {
  position: absolute;
  top: 10px;
  left: 10px;
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

.zoom-btn {
  position: absolute;
  bottom: 8px;
  right: 8px;
  width: 28px; height: 28px;
  background: rgba(255,255,255,0.85);
  border: 1px solid rgba(0,0,0,0.1);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  cursor: pointer;
  transition: background 0.15s;
  padding: 0;
  line-height: 1;
}
.zoom-btn:hover { background: white; }

.unavailable { opacity: 0.7; cursor: not-allowed; }
.disabled { opacity: 0.45; cursor: not-allowed; pointer-events: none; }
.disabled:hover { transform: none; box-shadow: none; }

.flower-info { padding: 14px; }

.flower-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 4px;
  flex-wrap: wrap;
}

.flower-name { font-weight: 500; font-size: 0.9rem; color: var(--charcoal); }

.flower-meaning {
  font-size: 0.78rem;
  color: var(--warm-gray);
  margin-bottom: 6px;
  line-height: 1.4;
}

.flower-price {
  font-size: 0.85rem;
  font-weight: 500;
  color: var(--deep-rose);
  margin-bottom: 8px;
}

.flower-colors {
  display: flex;
  gap: 4px;
  margin-bottom: 8px;
}

.color-dot {
  width: 14px; height: 14px;
  border-radius: 50%;
  display: inline-block;
}

/* BUG FIX #1: Qty control styling — tidak dipengaruhi z-index popup */
.qty-control {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-top: 10px;
  padding: 6px 12px;
  background: var(--cream);
  border: 1px solid var(--blush);
  border-radius: var(--radius-pill);
}

.qty-btn {
  width: 28px; height: 28px;
  border-radius: 50%;
  border: 1.5px solid var(--rose);
  background: var(--white);
  color: var(--deep-rose);
  font-size: 1.1rem;
  font-weight: 700;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
  transition: var(--transition);
}

.qty-btn:hover {
  background: var(--deep-rose);
  color: white;
}

.qty-value {
  font-weight: 600;
  font-size: 1rem;
  color: var(--charcoal);
  min-width: 24px;
  text-align: center;
}
</style>
