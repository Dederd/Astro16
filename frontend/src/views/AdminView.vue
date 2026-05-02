<template>
  <div class="admin-page">
    <!-- Sidebar -->
    <aside class="admin-sidebar">
      <div class="sidebar-brand">
        <span>🌸</span>
        <span>Admin Panel</span>
      </div>
      <nav class="sidebar-nav">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          class="nav-item"
          :class="{ active: activeTab === tab.id }"
          @click="activeTab = tab.id"
        >
          <span class="nav-icon">{{ tab.icon }}</span>
          <span>{{ tab.label }}</span>
        </button>
      </nav>
    </aside>

    <!-- Main -->
    <main class="admin-main">
      <!-- Stats dashboard -->
      <div v-if="activeTab === 'dashboard'">
        <h1 class="admin-title">Dashboard</h1>
        <div v-if="statsLoading" class="loading-row">
          <div class="skeleton stat-skeleton" v-for="i in 4" :key="i"></div>
        </div>
        <div v-else class="stats-grid">
          <div class="stat-card">
            <div class="stat-icon">📋</div>
            <div class="stat-body">
              <span class="stat-value">{{ stats.total_orders || 0 }}</span>
              <span class="stat-label">Total Pesanan</span>
            </div>
          </div>
          <div class="stat-card green">
            <div class="stat-icon">✅</div>
            <div class="stat-body">
              <span class="stat-value">{{ stats.paid_orders || 0 }}</span>
              <span class="stat-label">Pesanan Dibayar</span>
            </div>
          </div>
          <div class="stat-card yellow">
            <div class="stat-icon">⏳</div>
            <div class="stat-body">
              <span class="stat-value">{{ stats.pending_orders || 0 }}</span>
              <span class="stat-label">Pesanan Pending</span>
            </div>
          </div>
          <div class="stat-card rose">
            <div class="stat-icon">💰</div>
            <div class="stat-body">
              <span class="stat-value">Rp{{ formatPrice(stats.total_revenue || 0) }}</span>
              <span class="stat-label">Total Pendapatan</span>
            </div>
          </div>
        </div>

        <!-- Recent orders mini table -->
        <div class="section-card" style="margin-top: 32px;">
          <div class="section-header">
            <h2>Pesanan Terbaru</h2>
            <button class="btn btn-ghost btn-sm" @click="activeTab = 'orders'">Lihat Semua</button>
          </div>
          <table class="admin-table" v-if="orders.length > 0">
            <thead><tr>
              <th>ID Pesanan</th><th>Nama</th><th>Desain</th><th>Total</th><th>Status</th>
            </tr></thead>
            <tbody>
              <tr v-for="order in orders.slice(0, 5)" :key="order.id">
                <td class="mono">{{ order.id }}</td>
                <td>{{ order.customer_name }}</td>
                <td>{{ order.design_name || order.catalog_item_id || '—' }}</td>
                <td>Rp{{ formatPrice(order.total_amount) }}</td>
                <td><span class="status-pill" :class="order.status">{{ order.status }}</span></td>
              </tr>
            </tbody>
          </table>
          <div v-else class="empty-inline">Belum ada pesanan</div>
        </div>
      </div>

      <!-- Orders tab -->
      <div v-if="activeTab === 'orders'">
        <h1 class="admin-title">Manajemen Pesanan</h1>
        <div class="section-card">
          <div v-if="ordersLoading" class="loading-center">Memuat pesanan...</div>
          <div v-else>
            <table class="admin-table" v-if="orders.length > 0">
              <thead><tr>
                <th>ID</th><th>Nama</th><th>Total</th><th>Status</th><th>Kurir</th><th>Resi</th><th>Aksi</th>
              </tr></thead>
              <tbody>
                <tr v-for="order in orders" :key="order.id">
                  <td class="mono">{{ order.id }}</td>
                  <td>
                    <div>{{ order.customer_name }}</div>
                    <div class="sub-text">{{ order.customer_phone }}</div>
                  </td>
                  <td>Rp{{ formatPrice(order.total_amount) }}</td>
                  <td>
                    <select class="status-select" :value="order.status" @change="updateOrderStatus(order, $event.target.value)">
                      <option value="pending">pending</option>
                      <option value="paid">paid</option>
                      <option value="processing">processing</option>
                      <option value="shipped">shipped</option>
                      <option value="delivered">delivered</option>
                      <option value="failed">failed</option>
                      <option value="cancelled">cancelled</option>
                    </select>
                  </td>
                  <td>{{ order.courier_service?.toUpperCase() || '—' }}</td>
                  <td>
                    <input
                      class="resi-input"
                      :value="order.tracking_number || ''"
                      @blur="updateResi(order, $event.target.value)"
                      placeholder="Isi resi..."
                    />
                  </td>
                  <td>
                    <button class="btn-icon" title="Detail" @click="viewOrderDetail(order)">👁️</button>
                  </td>
                </tr>
              </tbody>
            </table>
            <div v-else class="empty-inline">Belum ada pesanan</div>
          </div>
        </div>
      </div>

      <!-- Flowers tab -->
      <div v-if="activeTab === 'flowers'">
        <h1 class="admin-title">Manajemen Bunga</h1>
        <div class="section-card">
          <div v-if="flowersLoading" class="loading-center">Memuat bunga...</div>
          <div v-else>
            <table class="admin-table">
              <thead><tr>
                <th>ID</th><th>Nama</th><th>Harga</th><th>Stok</th><th>Status</th><th>Gambar URL</th><th>Aksi</th>
              </tr></thead>
              <tbody>
                <tr v-for="flower in adminFlowers" :key="flower.id">
                  <td class="mono">{{ flower.id }}</td>
                  <td>
                    <div>{{ flower.name_id }}</div>
                    <div class="sub-text">{{ flower.emoji }}</div>
                  </td>
                  <td>
                    <input
                      class="price-input"
                      type="number"
                      :value="flower.price"
                      @blur="updateFlowerField(flower, 'price', parseInt($event.target.value))"
                    />
                  </td>
                  <td>
                    <input
                      class="stock-input"
                      type="number"
                      :value="flower.stock"
                      @blur="updateFlowerField(flower, 'stock', parseInt($event.target.value))"
                    />
                  </td>
                  <td>
                    <label class="toggle">
                      <input
                        type="checkbox"
                        :checked="flower.is_available"
                        @change="toggleFlowerAvailability(flower)"
                      />
                      <span class="toggle-slider"></span>
                    </label>
                  </td>
                  <td>
                    <input
                      class="url-input"
                      type="text"
                      :value="flower.image_url || ''"
                      placeholder="URL gambar..."
                      @blur="updateFlowerField(flower, 'image_url', $event.target.value)"
                    />
                  </td>
                  <td>
                    <span class="saved-indicator" v-if="savedFlowers[flower.id]">✓</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Catalog tab -->
      <div v-if="activeTab === 'catalog'">
        <div class="section-header-main">
          <h1 class="admin-title">Manajemen Katalog</h1>
          <button class="btn btn-primary btn-sm" @click="showCatalogForm = true">+ Tambah</button>
        </div>

        <!-- Add/Edit form -->
        <div v-if="showCatalogForm" class="catalog-form-card card">
          <h3>{{ editingCatalog ? 'Edit' : 'Tambah' }} Item Katalog</h3>
          <div class="form-grid">
            <div class="form-group">
              <label class="form-label">ID</label>
              <input v-model="catalogForm.id" class="form-input" placeholder="cat_nama_unik" :disabled="!!editingCatalog" />
            </div>
            <div class="form-group">
              <label class="form-label">Nama</label>
              <input v-model="catalogForm.name" class="form-input" placeholder="Nama bouquet" />
            </div>
            <div class="form-group" style="grid-column: span 2">
              <label class="form-label">Deskripsi</label>
              <textarea v-model="catalogForm.description" class="form-input" rows="2" style="resize:vertical;"></textarea>
            </div>
            <div class="form-group" style="grid-column: span 2">
              <label class="form-label">URL Gambar</label>
              <input v-model="catalogForm.image_url" class="form-input" placeholder="https://..." />
            </div>
            <div class="form-group">
              <label class="form-label">Style</label>
              <select v-model="catalogForm.style" class="form-input">
                <option v-for="s in styles" :key="s" :value="s">{{ s }}</option>
              </select>
            </div>
            <div class="form-group">
              <label class="form-label">Occasion</label>
              <select v-model="catalogForm.occasion" class="form-input">
                <option v-for="o in occasions" :key="o.value" :value="o.value">{{ o.label }}</option>
              </select>
            </div>
            <div class="form-group">
              <label class="form-label">Harga</label>
              <input v-model.number="catalogForm.price" type="number" class="form-input" />
            </div>
            <div class="form-group">
              <label class="form-label">Jumlah Tangkai</label>
              <input v-model.number="catalogForm.stem_count" type="number" class="form-input" />
            </div>
            <div class="form-group">
              <label class="form-label">Stok</label>
              <input v-model.number="catalogForm.stock" type="number" class="form-input" />
            </div>
            <div class="form-group">
              <label class="form-label">Sort Order</label>
              <input v-model.number="catalogForm.sort_order" type="number" class="form-input" />
            </div>
          </div>
          <div class="form-actions">
            <button class="btn btn-ghost btn-sm" @click="cancelCatalogForm">Batal</button>
            <button class="btn btn-primary btn-sm" @click="saveCatalog">
              {{ editingCatalog ? 'Update' : 'Simpan' }}
            </button>
          </div>
        </div>

        <div class="section-card">
          <div v-if="catalogLoading" class="loading-center">Memuat katalog...</div>
          <div v-else>
            <table class="admin-table" v-if="adminCatalog.length > 0">
              <thead><tr>
                <th>Nama</th><th>Style</th><th>Occasion</th><th>Harga</th><th>Stok</th><th>Aktif</th><th>Aksi</th>
              </tr></thead>
              <tbody>
                <tr v-for="item in adminCatalog" :key="item.id">
                  <td>{{ item.name }}</td>
                  <td>{{ item.style }}</td>
                  <td>{{ item.occasion }}</td>
                  <td>Rp{{ formatPrice(item.price) }}</td>
                  <td>{{ item.stock }}</td>
                  <td>
                    <span :class="item.is_available ? 'badge-green' : 'badge-gray'">
                      {{ item.is_available ? 'Aktif' : 'Nonaktif' }}
                    </span>
                  </td>
                  <td class="action-cell">
                    <button class="btn-icon" @click="editCatalog(item)">✏️</button>
                    <button class="btn-icon danger" @click="deleteCatalogItem(item.id)">🗑️</button>
                  </td>
                </tr>
              </tbody>
            </table>
            <div v-else class="empty-inline">Belum ada item katalog</div>
          </div>
        </div>
      </div>

      <!-- Order detail modal -->
      <div v-if="selectedOrder" class="modal-overlay" @click.self="selectedOrder = null">
        <div class="modal-card">
          <div class="modal-header">
            <h2>Detail Pesanan</h2>
            <button class="btn-close" @click="selectedOrder = null">✕</button>
          </div>
          <div class="modal-body">
            <div class="detail-grid">
              <div class="detail-item"><span>ID</span><strong class="mono">{{ selectedOrder.id }}</strong></div>
              <div class="detail-item"><span>Nama</span><strong>{{ selectedOrder.customer_name }}</strong></div>
              <div class="detail-item"><span>Email</span><strong>{{ selectedOrder.customer_email }}</strong></div>
              <div class="detail-item"><span>HP</span><strong>{{ selectedOrder.customer_phone }}</strong></div>
              <div class="detail-item"><span>Desain</span><strong>{{ selectedOrder.design_name || '—' }}</strong></div>
              <div class="detail-item"><span>Ukuran</span><strong>{{ selectedOrder.size }}</strong></div>
              <div class="detail-item"><span>Total</span><strong class="price-text">Rp{{ formatPrice(selectedOrder.total_amount) }}</strong></div>
              <div class="detail-item"><span>Status</span><span class="status-pill" :class="selectedOrder.status">{{ selectedOrder.status }}</span></div>
              <div class="detail-item full"><span>Alamat Kirim</span><strong>{{ selectedOrder.shipping_address }}, {{ selectedOrder.shipping_city }} {{ selectedOrder.shipping_postcode }}</strong></div>
              <div class="detail-item"><span>Kurir</span><strong>{{ selectedOrder.courier_service?.toUpperCase() || '—' }}</strong></div>
              <div class="detail-item"><span>Resi</span><strong>{{ selectedOrder.tracking_number || '—' }}</strong></div>
              <div class="detail-item"><span>Status Kirim</span><strong>{{ selectedOrder.shipping_status || '—' }}</strong></div>
              <div class="detail-item full"><span>Catatan</span><strong>{{ selectedOrder.notes || '—' }}</strong></div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {
  adminGetStats, adminGetOrders, adminUpdateOrder,
  adminGetFlowers, adminUpdateFlower,
  adminGetCatalog, adminCreateCatalog, adminUpdateCatalog, adminDeleteCatalog
} from '@/services/api'

const activeTab = ref('dashboard')
const tabs = [
  { id: 'dashboard', icon: '📊', label: 'Dashboard' },
  { id: 'orders', icon: '📋', label: 'Pesanan' },
  { id: 'flowers', icon: '🌸', label: 'Bunga' },
  { id: 'catalog', icon: '🛍️', label: 'Katalog' },
]

// Stats
const stats = ref({})
const statsLoading = ref(true)

// Orders
const orders = ref([])
const ordersLoading = ref(true)
const selectedOrder = ref(null)

// Flowers
const adminFlowers = ref([])
const flowersLoading = ref(true)
const savedFlowers = ref({})

// Catalog
const adminCatalog = ref([])
const catalogLoading = ref(true)
const showCatalogForm = ref(false)
const editingCatalog = ref(null)
const catalogForm = ref(defaultCatalogForm())

const styles = ['Classic', 'Premium', 'Romantic', 'Playful', 'Elegant', 'Warm', 'Modern', 'Natural']
const occasions = [
  { value: 'graduation', label: '🎓 Wisuda' },
  { value: 'wedding', label: '💍 Pernikahan' },
  { value: 'birthday', label: '🎂 Ulang Tahun' },
  { value: 'valentines', label: '❤️ Valentine' },
  { value: 'mothers_day', label: '🌸 Hari Ibu' },
  { value: 'anniversary', label: '💑 Anniversary' },
  { value: 'congratulations', label: '🏆 Selamat' },
  { value: 'sympathy', label: '🕊️ Dukacita' },
]

function defaultCatalogForm() {
  return { id: '', name: '', description: '', image_url: '', style: 'Classic', occasion: 'graduation', price: 0, stem_count: 0, stock: 10, sort_order: 0, is_available: true }
}

onMounted(async () => {
  await Promise.all([loadStats(), loadOrders(), loadFlowers(), loadCatalog()])
})

async function loadStats() {
  statsLoading.value = true
  try { const r = await adminGetStats(); stats.value = r.data } catch { /* ignore */ } finally { statsLoading.value = false }
}

async function loadOrders() {
  ordersLoading.value = true
  try { const r = await adminGetOrders(); orders.value = r.data.data || [] } catch { /* ignore */ } finally { ordersLoading.value = false }
}

async function loadFlowers() {
  flowersLoading.value = true
  try { const r = await adminGetFlowers(); adminFlowers.value = r.data.data || [] } catch { /* ignore */ } finally { flowersLoading.value = false }
}

async function loadCatalog() {
  catalogLoading.value = true
  try { const r = await adminGetCatalog(); adminCatalog.value = r.data.data || [] } catch { /* ignore */ } finally { catalogLoading.value = false }
}

function viewOrderDetail(order) { selectedOrder.value = order }

async function updateOrderStatus(order, newStatus) {
  try {
    await adminUpdateOrder(order.id, {
      status: newStatus,
      tracking_number: order.tracking_number || '',
      shipping_status: order.shipping_status || '',
      courier_service: order.courier_service || '',
    })
    order.status = newStatus
  } catch (e) { alert('Gagal update status: ' + (e?.response?.data?.error || e.message)) }
}

async function updateResi(order, resi) {
  if (resi === (order.tracking_number || '')) return
  try {
    await adminUpdateOrder(order.id, {
      status: order.status,
      tracking_number: resi,
      shipping_status: resi ? 'shipped' : order.shipping_status,
      courier_service: order.courier_service || '',
    })
    order.tracking_number = resi
    if (resi) order.shipping_status = 'shipped'
  } catch (e) { alert('Gagal update resi: ' + (e?.response?.data?.error || e.message)) }
}

async function updateFlowerField(flower, field, value) {
  const payload = { [field]: value }
  try {
    await adminUpdateFlower(flower.id, payload)
    flower[field] = value
    savedFlowers.value[flower.id] = true
    setTimeout(() => { savedFlowers.value[flower.id] = false }, 2000)
  } catch (e) { alert('Gagal update bunga: ' + (e?.response?.data?.error || e.message)) }
}

async function toggleFlowerAvailability(flower) {
  await updateFlowerField(flower, 'is_available', !flower.is_available)
}

function editCatalog(item) {
  editingCatalog.value = item.id
  catalogForm.value = { ...item }
  showCatalogForm.value = true
}

function cancelCatalogForm() {
  showCatalogForm.value = false
  editingCatalog.value = null
  catalogForm.value = defaultCatalogForm()
}

async function saveCatalog() {
  try {
    if (editingCatalog.value) {
      await adminUpdateCatalog(editingCatalog.value, catalogForm.value)
    } else {
      await adminCreateCatalog(catalogForm.value)
    }
    await loadCatalog()
    cancelCatalogForm()
  } catch (e) { alert('Gagal simpan: ' + (e?.response?.data?.error || e.message)) }
}

async function deleteCatalogItem(id) {
  if (!confirm('Yakin hapus item ini?')) return
  try { await adminDeleteCatalog(id); await loadCatalog() }
  catch (e) { alert('Gagal hapus: ' + (e?.response?.data?.error || e.message)) }
}

function formatPrice(p) { return (p || 0).toLocaleString('id-ID') }
</script>

<style scoped>
.admin-page { display: flex; min-height: 100vh; background: #F8F4EF; }

/* Sidebar */
.admin-sidebar {
  width: 220px;
  background: var(--charcoal);
  color: white;
  padding: 24px 0;
  flex-shrink: 0;
  position: sticky;
  top: 0;
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 20px 24px;
  font-size: 1.1rem;
  font-weight: 600;
  border-bottom: 1px solid rgba(255,255,255,0.1);
  margin-bottom: 16px;
}

.sidebar-nav { display: flex; flex-direction: column; gap: 2px; padding: 0 10px; }
.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  border-radius: 10px;
  border: none;
  background: none;
  color: rgba(255,255,255,0.7);
  cursor: pointer;
  font-size: 0.9rem;
  transition: var(--transition);
  text-align: left;
}
.nav-item:hover { background: rgba(255,255,255,0.08); color: white; }
.nav-item.active { background: var(--deep-rose); color: white; }
.nav-icon { font-size: 1.1rem; }

/* Main content */
.admin-main { flex: 1; padding: 32px; overflow-y: auto; max-height: 100vh; }
.admin-title { font-size: 1.6rem; color: var(--charcoal); margin-bottom: 24px; }
.section-header-main { display: flex; align-items: center; justify-content: space-between; margin-bottom: 0; }

/* Stats */
.loading-row { display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px; }
.stat-skeleton { height: 100px; }
.stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 16px; }
.stat-card { background: white; border-radius: var(--radius); padding: 20px; display: flex; align-items: center; gap: 16px; border: 1px solid var(--light-gray); }
.stat-card.green { border-left: 4px solid #43A047; }
.stat-card.yellow { border-left: 4px solid #F9A825; }
.stat-card.rose { border-left: 4px solid var(--deep-rose); }
.stat-icon { font-size: 2rem; }
.stat-body { display: flex; flex-direction: column; gap: 2px; }
.stat-value { font-size: 1.5rem; font-weight: 700; color: var(--charcoal); }
.stat-label { font-size: 0.78rem; color: var(--warm-gray); }

/* Section card */
.section-card { background: white; border-radius: var(--radius); padding: 24px; margin-top: 20px; }
.section-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.section-header h2 { font-size: 1.05rem; color: var(--charcoal); }

/* Table */
.admin-table { width: 100%; border-collapse: collapse; font-size: 0.85rem; }
.admin-table th { text-align: left; padding: 10px 12px; border-bottom: 2px solid var(--light-gray); color: var(--warm-gray); font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.05em; }
.admin-table td { padding: 12px; border-bottom: 1px solid var(--light-gray); vertical-align: middle; }
.admin-table tr:hover td { background: #FAFAFA; }
.mono { font-family: monospace; font-size: 0.8rem; color: var(--warm-gray); }
.sub-text { font-size: 0.75rem; color: var(--warm-gray); margin-top: 2px; }

/* Status pill */
.status-pill { padding: 3px 10px; border-radius: var(--radius-pill); font-size: 0.73rem; font-weight: 600; text-transform: uppercase; letter-spacing: 0.05em; }
.status-pill.paid, .badge-green { background: #E8F5E9; color: #2E7D32; }
.status-pill.pending { background: #FFF8E1; color: #F57F17; }
.status-pill.processing { background: #E3F2FD; color: #1565C0; }
.status-pill.shipped { background: #F3E5F5; color: #6A1B9A; }
.status-pill.delivered { background: #E8F5E9; color: #1B5E20; }
.status-pill.failed, .status-pill.cancelled { background: #FFEBEE; color: #C62828; }
.badge-gray { background: #ECEFF1; color: #546E7A; padding: 3px 10px; border-radius: var(--radius-pill); font-size: 0.73rem; font-weight: 600; }

/* Inline inputs */
.status-select, .resi-input, .price-input, .stock-input, .url-input {
  border: 1px solid var(--light-gray);
  border-radius: 6px;
  padding: 5px 8px;
  font-size: 0.82rem;
  background: white;
  outline: none;
  transition: var(--transition);
}
.status-select:focus, .resi-input:focus, .price-input:focus, .stock-input:focus, .url-input:focus { border-color: var(--rose); }
.resi-input { width: 130px; }
.price-input, .stock-input { width: 80px; }
.url-input { width: 180px; }

/* Toggle switch */
.toggle { position: relative; display: inline-block; width: 36px; height: 20px; }
.toggle input { display: none; }
.toggle-slider { position: absolute; inset: 0; background: #ccc; border-radius: 20px; cursor: pointer; transition: .3s; }
.toggle-slider:before { content: ''; position: absolute; width: 14px; height: 14px; background: white; border-radius: 50%; left: 3px; bottom: 3px; transition: .3s; }
.toggle input:checked + .toggle-slider { background: var(--deep-rose); }
.toggle input:checked + .toggle-slider:before { transform: translateX(16px); }

.saved-indicator { color: #43A047; font-weight: 600; font-size: 0.85rem; }

/* Catalog form */
.catalog-form-card { padding: 24px; margin: 20px 0; }
.catalog-form-card h3 { font-size: 1rem; margin-bottom: 20px; }
.form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; margin-bottom: 20px; }
.form-actions { display: flex; gap: 10px; justify-content: flex-end; }

/* Buttons */
.btn-icon { background: none; border: none; cursor: pointer; font-size: 1rem; padding: 4px 6px; border-radius: 6px; transition: var(--transition); }
.btn-icon:hover { background: var(--cream); }
.btn-icon.danger:hover { background: #FFEBEE; }
.action-cell { display: flex; gap: 4px; }
.btn-sm { padding: 8px 14px; font-size: 0.82rem; }

/* Modal */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.45); z-index: 1000; display: flex; align-items: center; justify-content: center; padding: 20px; }
.modal-card { background: white; border-radius: var(--radius); width: 100%; max-width: 560px; max-height: 90vh; overflow-y: auto; }
.modal-header { display: flex; align-items: center; justify-content: space-between; padding: 20px 24px; border-bottom: 1px solid var(--light-gray); }
.modal-header h2 { font-size: 1.1rem; }
.btn-close { background: none; border: none; font-size: 1.1rem; cursor: pointer; color: var(--warm-gray); }
.modal-body { padding: 24px; }
.detail-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.detail-item { display: flex; flex-direction: column; gap: 2px; }
.detail-item.full { grid-column: span 2; }
.detail-item span { font-size: 0.72rem; color: var(--warm-gray); text-transform: uppercase; letter-spacing: 0.05em; }
.detail-item strong { font-size: 0.9rem; color: var(--charcoal); }
.price-text { color: var(--deep-rose) !important; }

.empty-inline { color: var(--warm-gray); font-size: 0.88rem; text-align: center; padding: 32px; }
.loading-center { text-align: center; padding: 32px; color: var(--warm-gray); font-size: 0.9rem; }
</style>
