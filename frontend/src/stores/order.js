import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useOrderStore = defineStore('order', () => {
  // Step tracking
  const currentStep = ref(1) // 1:type 2:flowers 3:designs 4:checkout

  // Step 1: Bouquet type
  const selectedBouquetType = ref(null)
  const agentMessage = ref('')
  const agentTips = ref('')
  const recommendedFlowerIds = ref([])

  // Step 2: Selected flowers
  const selectedFlowers = ref([]) // [{flower_id, name, quantity: 1}, ...]

  // Step 3: Generated designs
  const generatedDesigns = ref([])
  const designMessage = ref('')
  const selectedDesign = ref(null)
  const selectedSize = ref(null) // 'small' | 'large'
  const generateCount = ref(0)
  const generateLimit = ref(2)
  const extraQuota = ref(0)
  const extraQuotaFee = ref(0)

  // Optional AI prompt hints (from Step 2)
  const aiStyleHint = ref('')
  const aiDescriptionHint = ref('')

  // Catalog mode
  const orderMode = ref('ai') // 'ai' | 'catalog'
  const selectedCatalogItem = ref(null)

  // Step 4
  const createdOrder = ref(null)

  // Session ID (untuk rate limiting)
  const sessionId = ref(
    localStorage.getItem('bouquet_session') || (() => {
      const id = 'sess_' + Math.random().toString(36).slice(2) + Date.now()
      localStorage.setItem('bouquet_session', id)
      return id
    })()
  )

  // Biaya tetap
  const MAKING_FEE = 5000   // biaya pembuatan bouquet
  const AI_FEE = 5000       // biaya AI generate (jika pakai & sudah habis kuota gratis)
  const SHIPPING_FEE_DEFAULT = 15000 // default, akan di-override dari kurir

  // Apakah generate ini berbayar (user sudah habis kuota gratis)
  const aiFeePaid = ref(false)   // di-set dari response AgentGenerateBouquet

  // Harga sekarang flat dari design (mini 35k, premium 75k)
  const flowerCost = computed(() => 0)

  const shippingCost = ref(0) // di-set saat user pilih kurir

  // Total breakdown
  const totalPrice = computed(() => {
    if (orderMode.value === 'catalog' && selectedCatalogItem.value) {
      return selectedCatalogItem.value.price + shippingCost.value
    }
    if (!selectedDesign.value || !selectedSize.value) return 0
    const stemPrice = selectedSize.value === 'small'
      ? selectedDesign.value.small.price
      : selectedDesign.value.large.price
    const ai = aiFeePaid.value ? AI_FEE : 0
    return stemPrice + MAKING_FEE + ai + extraQuotaFee.value + shippingCost.value
  })

  const priceBreakdown = computed(() => {
    if (orderMode.value === 'catalog' && selectedCatalogItem.value) {
      return {
        flower_cost: selectedCatalogItem.value.price,
        making_fee: 0,
        ai_fee: 0,
        shipping_cost: shippingCost.value,
        total: selectedCatalogItem.value.price + shippingCost.value,
      }
    }
    if (!selectedDesign.value || !selectedSize.value) return null
    const stemPrice = selectedSize.value === 'small'
      ? selectedDesign.value.small.price
      : selectedDesign.value.large.price
    const ai = aiFeePaid.value ? AI_FEE : 0
    return {
      flower_cost: stemPrice,
      making_fee: MAKING_FEE,
      ai_fee: ai,
      extra_quota_fee: extraQuotaFee.value,
      shipping_cost: shippingCost.value,
      total: stemPrice + MAKING_FEE + ai + extraQuotaFee.value + shippingCost.value,
    }
  })

  const flowerCount = computed(() => selectedFlowers.value.length)

  const isGenerateLimited = computed(() =>
    generateCount.value >= generateLimit.value
  )

  function setExtraQuota(eq, fee) {
    extraQuota.value = eq
    extraQuotaFee.value = fee
    generateLimit.value = 2 + eq
  }

  // Actions
  function setStep(step) { currentStep.value = step }

  function selectBouquetType(type) { selectedBouquetType.value = type }

  function setAgentResponse(msg, tips, recs) {
    agentMessage.value = msg
    agentTips.value = tips
    recommendedFlowerIds.value = recs
  }

  function toggleFlower(flower) {
    const idx = selectedFlowers.value.findIndex(f => f.flower_id === flower.id)
    if (idx === -1) {
      selectedFlowers.value.push({ flower_id: flower.id, name: flower.name_id, quantity: 1 })
    } else {
      selectedFlowers.value.splice(idx, 1)
    }
  }

  function isFlowerSelected(flowerId) {
    return selectedFlowers.value.some(f => f.flower_id === flowerId)
  }

  function setGeneratedDesigns(designs, message, count, limit) {
    generatedDesigns.value = designs
    designMessage.value = message
    if (count !== undefined) generateCount.value = count
    if (limit !== undefined) generateLimit.value = limit
  }

  function selectDesign(design, size) {
    selectedDesign.value = design
    selectedSize.value = size
  }

  function setOrderMode(mode) {
    orderMode.value = mode
    selectedCatalogItem.value = null
    selectedDesign.value = null
    selectedSize.value = null
  }

  function selectCatalogItem(item) {
    selectedCatalogItem.value = item
    orderMode.value = 'catalog'
  }

  function setCreatedOrder(order) { createdOrder.value = order }

  function setAIFeePaid(paid) { aiFeePaid.value = paid }
  function setShippingCost(cost) { shippingCost.value = cost }

  function reset() {
    currentStep.value = 1
    selectedBouquetType.value = null
    agentMessage.value = ''
    agentTips.value = ''
    recommendedFlowerIds.value = []
    selectedFlowers.value = []
    generatedDesigns.value = []
    designMessage.value = ''
    selectedDesign.value = null
    selectedSize.value = null
    orderMode.value = 'ai'
    selectedCatalogItem.value = null
    createdOrder.value = null
    aiFeePaid.value = false
    shippingCost.value = 0
    extraQuota.value = 0
    extraQuotaFee.value = 0
    aiStyleHint.value = ''
    aiDescriptionHint.value = ''
  }

  return {
    currentStep, selectedBouquetType, agentMessage, agentTips,
    recommendedFlowerIds, selectedFlowers, generatedDesigns, designMessage,
    selectedDesign, selectedSize, generateCount, generateLimit,
    extraQuota, extraQuotaFee,
    aiStyleHint, aiDescriptionHint,
    orderMode, selectedCatalogItem, createdOrder, totalPrice, flowerCount,
    isGenerateLimited, sessionId, aiFeePaid, shippingCost, priceBreakdown,
    MAKING_FEE, AI_FEE,
    setStep, selectBouquetType, setAgentResponse, toggleFlower,
    isFlowerSelected,
    setGeneratedDesigns, selectDesign, setOrderMode, selectCatalogItem,
    setCreatedOrder, setAIFeePaid, setShippingCost, setExtraQuota, reset,
  }
})
