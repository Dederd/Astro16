<template>
  <div class="order-page">
    <!-- Progress bar -->
    <div class="progress-bar-wrap">
      <div class="container">
        <StepIndicator :current-step="store.currentStep" />
      </div>
    </div>

    <!-- Step content -->
    <div class="container order-content">
      <!-- STEP 1: Pilih Momen -->
      <Step1BouquetType v-if="store.currentStep === 1" />

      <!-- STEP 2: Pilih Bunga -->
      <Step2Flowers v-else-if="store.currentStep === 2" />

      <!-- STEP 3: Generate & Pilih Desain -->
      <Step3Design v-else-if="store.currentStep === 3" />

      <!-- STEP 4: Checkout -->
      <Step4Checkout v-else-if="store.currentStep === 4" />
    </div>
  </div>
</template>

<script setup>
import { useOrderStore } from '@/stores/order'
import StepIndicator from '@/components/StepIndicator.vue'
import Step1BouquetType from '@/components/steps/Step1BouquetType.vue'
import Step2Flowers from '@/components/steps/Step2Flowers.vue'
import Step3Design from '@/components/steps/Step3Design.vue'
import Step4Checkout from '@/components/steps/Step4Checkout.vue'

const store = useOrderStore()
</script>

<style scoped>
.order-page { min-height: 80vh; }

.progress-bar-wrap {
  background: var(--white);
  border-bottom: 1px solid var(--light-gray);
  padding: 20px 0;
  position: sticky;
  top: 64px;
  z-index: 50;
}

.order-content {
  padding-top: 48px;
  padding-bottom: 80px;
}
</style>
