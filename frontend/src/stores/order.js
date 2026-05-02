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
  const selectedFlowers = ref([]) // [{flower_id, name, quantity}]

  // Step 3: Generated designs
  const generatedDesigns = ref([])
  const designMessage = ref('')
  const selectedDesign = ref(null)
  const selectedSize = ref(null) // 'small' | 'large'
  const generateCount = ref(0)
  const generateLimit = ref(2)

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

  // Computed
  const totalPrice = computed(() => {
    if (orderMode.value === 'catalog' && selectedCatalogItem.value) {
      return selectedCatalogItem.value.price
    }
    if (!selectedDesign.value || !selectedSize.value) return 0
    return selectedSize.value === 'small'
      ? selectedDesign.value.small.price
      : selectedDesign.value.large.price
  })

  const flowerCount = computed(() =>
    selectedFlowers.value.reduce((sum, f) => sum + f.quantity, 0)
  )

  const isGenerateLimited = computed(() =>
    generateCount.value >= generateLimit.value
  )

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

  function updateFlowerQuantity(flowerId, quantity) {
    const f = selectedFlowers.value.find(f => f.flower_id === flowerId)
    if (f) {
      if (quantity <= 0) {
        selectedFlowers.value = selectedFlowers.value.filter(f => f.flower_id !== flowerId)
      } else {
        f.quantity = quantity
      }
    }
  }

  function isFlowerSelected(flowerId) {
    return selectedFlowers.value.some(f => f.flower_id === flowerId)
  }

  function getFlowerQuantity(flowerId) {
    const f = selectedFlowers.value.find(f => f.flower_id === flowerId)
    return f ? f.quantity : 0
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
  }

  return {
    currentStep, selectedBouquetType, agentMessage, agentTips,
    recommendedFlowerIds, selectedFlowers, generatedDesigns, designMessage,
    selectedDesign, selectedSize, generateCount, generateLimit,
    orderMode, selectedCatalogItem, createdOrder, totalPrice, flowerCount,
    isGenerateLimited, sessionId,
    setStep, selectBouquetType, setAgentResponse, toggleFlower,
    updateFlowerQuantity, isFlowerSelected, getFlowerQuantity,
    setGeneratedDesigns, selectDesign, setOrderMode, selectCatalogItem,
    setCreatedOrder, reset,
  }
})
