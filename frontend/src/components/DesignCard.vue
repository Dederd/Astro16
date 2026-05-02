<template>
  <div class="design-card" :class="{ selected: isSelected }">
    <!-- Design header -->
    <div class="design-visual" :style="`background: ${gradient}`">
      <span class="design-emoji">{{ emoji }}</span>
      <div class="design-style-badge">{{ design.style }}</div>
      <div v-if="isSelected" class="selected-overlay">
        <span class="selected-icon">✓</span>
      </div>
    </div>

    <div class="design-body">
      <h3 class="design-name">{{ design.name }}</h3>
      <p class="design-desc">{{ design.description }}</p>

      <!-- Size selector tabs -->
      <div class="size-tabs">
        <button
          class="size-tab"
          :class="{ active: isSelected && selectedSize === 'small' }"
          @click="$emit('select', design, 'small')"
        >
          <span class="size-icon">🌿</span>
          <span class="size-label">{{ design.small.label }}</span>
        </button>
        <button
          class="size-tab"
          :class="{ active: isSelected && selectedSize === 'large' }"
          @click="$emit('select', design, 'large')"
        >
          <span class="size-icon">💐</span>
          <span class="size-label">{{ design.large.label }}</span>
        </button>
      </div>

      <!-- Active size details -->
      <div v-if="isSelected && activeVariant" class="size-detail">
        <div class="size-detail-row">
          <span>{{ activeVariant.stem_count }} tangkai</span>
          <span class="size-price">Rp{{ formatPrice(activeVariant.price) }}</span>
        </div>
        <p class="size-desc-text">{{ activeVariant.description }}</p>
      </div>

      <!-- Both sizes preview if not selected -->
      <div v-if="!isSelected" class="sizes-preview">
        <div class="size-preview-item" @click="$emit('select', design, 'small')">
          <span class="size-preview-label">{{ design.small.label }}</span>
          <span class="size-preview-detail">{{ design.small.stem_count }} tangkai</span>
          <span class="size-preview-price">Rp{{ formatPrice(design.small.price) }}</span>
        </div>
        <div class="size-preview-divider"></div>
        <div class="size-preview-item" @click="$emit('select', design, 'large')">
          <span class="size-preview-label">{{ design.large.label }}</span>
          <span class="size-preview-detail">{{ design.large.stem_count }} tangkai</span>
          <span class="size-preview-price">Rp{{ formatPrice(design.large.price) }}</span>
        </div>
      </div>

      <button
        v-if="!isSelected"
        class="btn btn-outline select-btn"
        @click="$emit('select', design, 'small')"
      >
        Pilih Desain Ini
      </button>
      <button
        v-else
        class="btn btn-primary select-btn"
        disabled
      >
        ✓ Dipilih
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  design: { type: Object, required: true },
  isSelected: { type: Boolean, default: false },
  selectedSize: { type: String, default: 'small' },
})

defineEmits(['select'])

const gradients = {
  Romantic: 'linear-gradient(135deg, #FFECD2 0%, #FCB69F 100%)',
  Elegant: 'linear-gradient(135deg, #E8D5C4 0%, #C4A882 100%)',
  Modern: 'linear-gradient(135deg, #D4E8E0 0%, #8ABFAD 100%)',
  Natural: 'linear-gradient(135deg, #E8F5E8 0%, #8BC888 100%)',
  Classic: 'linear-gradient(135deg, #FFF0E6 0%, #F0C8A0 100%)',
}

const emojis = {
  Romantic: '💐', Elegant: '🌹', Modern: '🌿', Natural: '🌻', Classic: '🌸',
}

const gradient = computed(() => gradients[props.design.style] || gradients.Romantic)
const emoji = computed(() => emojis[props.design.style] || '💐')

const activeVariant = computed(() => {
  if (!props.isSelected) return null
  return props.selectedSize === 'small' ? props.design.small : props.design.large
})

function formatPrice(p) {
  return (p || 0).toLocaleString('id-ID')
}
</script>

<style scoped>
.design-card {
  background: var(--white);
  border-radius: var(--radius);
  border: 1.5px solid var(--light-gray);
  overflow: hidden;
  transition: var(--transition);
  display: flex;
  flex-direction: column;
}

.design-card:hover {
  border-color: var(--rose);
  box-shadow: 0 6px 24px rgba(198, 40, 40, 0.1);
  transform: translateY(-3px);
}

.design-card.selected {
  border-color: var(--deep-rose);
  box-shadow: 0 6px 28px rgba(139, 0, 0, 0.18);
}

.design-visual {
  position: relative;
  height: 160px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.design-emoji { font-size: 4rem; }

.design-style-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  background: rgba(255,255,255,0.85);
  color: var(--charcoal);
  font-size: 0.72rem;
  font-weight: 600;
  padding: 3px 10px;
  border-radius: var(--radius-pill);
  letter-spacing: 0.05em;
}

.selected-overlay {
  position: absolute;
  inset: 0;
  background: rgba(139,0,0,0.12);
  display: flex;
  align-items: center;
  justify-content: center;
}

.selected-icon {
  width: 44px; height: 44px;
  background: var(--deep-rose);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
  font-weight: 700;
}

.design-body {
  padding: 20px;
  display: flex;
  flex-direction: column;
  flex: 1;
}

.design-name { font-size: 1.05rem; font-weight: 600; color: var(--charcoal); margin-bottom: 6px; }
.design-desc { font-size: 0.82rem; color: var(--warm-gray); line-height: 1.5; margin-bottom: 16px; }

/* Size tabs */
.size-tabs { display: flex; gap: 8px; margin-bottom: 14px; }
.size-tab {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 3px;
  padding: 10px 6px;
  border: 1.5px solid var(--light-gray);
  border-radius: var(--radius-sm);
  background: var(--cream);
  cursor: pointer;
  transition: var(--transition);
  font-size: 0.78rem;
}
.size-tab:hover { border-color: var(--rose); }
.size-tab.active { border-color: var(--deep-rose); background: #FFF0F0; }
.size-icon { font-size: 1.1rem; }
.size-label { font-weight: 500; color: var(--charcoal); }

/* Size detail (when selected) */
.size-detail {
  background: var(--cream);
  border: 1px solid var(--blush);
  border-radius: var(--radius-sm);
  padding: 12px 14px;
  margin-bottom: 14px;
}
.size-detail-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}
.size-price { font-weight: 600; color: var(--deep-rose); font-size: 1.05rem; }
.size-desc-text { font-size: 0.78rem; color: var(--warm-gray); }

/* Sizes preview (when not selected) */
.sizes-preview {
  display: flex;
  gap: 0;
  border: 1px solid var(--light-gray);
  border-radius: var(--radius-sm);
  overflow: hidden;
  margin-bottom: 14px;
}
.size-preview-item {
  flex: 1;
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  cursor: pointer;
  transition: var(--transition);
}
.size-preview-item:hover { background: var(--cream); }
.size-preview-divider { width: 1px; background: var(--light-gray); }
.size-preview-label { font-size: 0.75rem; font-weight: 600; color: var(--charcoal); }
.size-preview-detail { font-size: 0.7rem; color: var(--warm-gray); }
.size-preview-price { font-size: 0.82rem; font-weight: 500; color: var(--deep-rose); }

.select-btn { width: 100%; justify-content: center; }
</style>
