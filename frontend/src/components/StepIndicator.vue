<template>
  <div class="steps-wrapper">
    <div class="step-indicator">
      <template v-for="(step, i) in steps" :key="i">
        <div
          class="step-dot"
          :class="{
            active: currentStep === i + 1,
            done: currentStep > i + 1,
            inactive: currentStep < i + 1
          }"
        >
          <span v-if="currentStep > i + 1">✓</span>
          <span v-else>{{ i + 1 }}</span>
        </div>
        <div v-if="i < steps.length - 1" class="step-line" :class="{ done: currentStep > i + 1 }"></div>
      </template>
    </div>
    <div class="step-labels">
      <span
        v-for="(step, i) in steps"
        :key="i"
        class="step-label"
        :class="{ active: currentStep === i + 1, done: currentStep > i + 1 }"
      >{{ step }}</span>
    </div>
  </div>
</template>

<script setup>
defineProps({
  currentStep: { type: Number, default: 1 },
  steps: { type: Array, default: () => ['Pilih Momen', 'Pilih Bunga', 'Desain AI', 'Checkout'] }
})
</script>

<style scoped>
.steps-wrapper { width: 100%; }

.step-indicator {
  display: flex;
  align-items: center;
  width: 100%;
}

.step-labels {
  display: flex;
  justify-content: space-between;
  margin-top: 8px;
}

.step-label {
  font-size: 0.72rem;
  color: var(--warm-gray);
  text-align: center;
  flex: 1;
  transition: var(--transition);
  letter-spacing: 0.02em;
}

.step-label.active { color: var(--deep-rose); font-weight: 500; }
.step-label.done { color: var(--sage); }

/* Inherited from global: .step-dot, .step-line */
</style>
