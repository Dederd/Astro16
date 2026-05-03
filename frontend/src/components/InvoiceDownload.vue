<template>
  <div>
    <button class="btn btn-outline btn-invoice" @click="downloadInvoice">
      📄 Download Invoice
    </button>

    <!-- Invoice preview (hidden, used for print/pdf) -->
    <div ref="invoiceEl" class="invoice-sheet" style="display:none;">
      <div class="inv-header">
        <div class="inv-brand">
          <span class="inv-logo">🌸</span>
          <span class="inv-brand-name">Bloome</span>
        </div>
        <div class="inv-meta">
          <h2>INVOICE</h2>
          <p><strong>No. Pesanan:</strong> {{ order.id }}</p>
          <p><strong>Tanggal:</strong> {{ formatDate(order.created_at) }}</p>
          <p><strong>Status:</strong> {{ statusLabel(order.status) }}</p>
        </div>
      </div>

      <div class="inv-parties">
        <div class="inv-party">
          <h4>Dari</h4>
          <p>Bloome Bouquet Studio</p>
          <p>support@bloome.id</p>
        </div>
        <div class="inv-party">
          <h4>Kepada</h4>
          <p><strong>{{ order.customer_name }}</strong></p>
          <p>{{ order.customer_email }}</p>
          <p>{{ order.customer_phone }}</p>
        </div>
      </div>

      <!-- Items table -->
      <table class="inv-table">
        <thead>
          <tr>
            <th>Item</th>
            <th>Qty</th>
            <th class="text-right">Harga</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="order.flower_cost > 0">
            <td>{{ order.design_name || 'Custom Bouquet' }} — Bunga</td>
            <td>1</td>
            <td class="text-right">Rp{{ formatPrice(order.flower_cost) }}</td>
          </tr>
          <tr v-if="order.making_fee > 0">
            <td>Biaya Pembuatan</td>
            <td>1</td>
            <td class="text-right">Rp{{ formatPrice(order.making_fee) }}</td>
          </tr>
          <tr v-if="order.ai_fee > 0">
            <td>Biaya AI Generate</td>
            <td>1</td>
            <td class="text-right">Rp{{ formatPrice(order.ai_fee) }}</td>
          </tr>
          <tr v-if="order.shipping_cost > 0">
            <td>Ongkos Kirim ({{ order.courier_service?.toUpperCase() }})</td>
            <td>1</td>
            <td class="text-right">Rp{{ formatPrice(order.shipping_cost) }}</td>
          </tr>
          <!-- Fallback jika breakdown tidak ada -->
          <tr v-if="!order.flower_cost && !order.making_fee">
            <td>{{ order.design_name || order.catalog_item_id || 'Bouquet' }}</td>
            <td>1</td>
            <td class="text-right">Rp{{ formatPrice(order.total_amount) }}</td>
          </tr>
        </tbody>
        <tfoot>
          <tr class="inv-total-row">
            <td colspan="2"><strong>Total</strong></td>
            <td class="text-right"><strong>Rp{{ formatPrice(order.total_amount) }}</strong></td>
          </tr>
        </tfoot>
      </table>

      <!-- Shipping info -->
      <div class="inv-shipping">
        <h4>📦 Informasi Pengiriman</h4>
        <div class="inv-shipping-grid">
          <div><span>Penerima</span><p>{{ order.customer_name }}</p></div>
          <div><span>No. HP</span><p>{{ order.customer_phone }}</p></div>
          <div class="full-width"><span>Alamat</span><p>{{ order.shipping_address }}, {{ order.shipping_city }} {{ order.shipping_postcode }}</p></div>
          <div><span>Kurir</span><p>{{ order.courier_service?.toUpperCase() || '—' }}</p></div>
          <div v-if="order.tracking_number"><span>No. Resi</span><p>{{ order.tracking_number }}</p></div>
        </div>
      </div>

      <!-- Notes -->
      <div v-if="order.notes" class="inv-notes">
        <strong>Catatan:</strong> {{ order.notes }}
      </div>

      <div class="inv-footer">
        <p>Terima kasih telah berbelanja di Bloome 🌸</p>
        <p>Pertanyaan? Hubungi kami di support@bloome.id</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  order: { type: Object, required: true }
})

const invoiceEl = ref(null)

function formatDate(dt) {
  if (!dt) return '-'
  return new Date(dt).toLocaleDateString('id-ID', {
    day: 'numeric', month: 'long', year: 'numeric',
    hour: '2-digit', minute: '2-digit'
  })
}

function formatPrice(p) {
  return (p || 0).toLocaleString('id-ID')
}

function statusLabel(s) {
  const map = { pending: 'Menunggu Pembayaran', paid: 'Lunas', processing: 'Diproses', shipped: 'Dikirim', delivered: 'Terkirim' }
  return map[s] || s
}

function downloadInvoice() {
  if (!invoiceEl.value) return

  // Show the invoice element temporarily for printing
  invoiceEl.value.style.display = 'block'

  const printWindow = window.open('', '_blank', 'width=800,height=900')
  printWindow.document.write(`
    <!DOCTYPE html>
    <html>
    <head>
      <meta charset="utf-8">
      <title>Invoice ${props.order.id}</title>
      <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { font-family: 'Segoe UI', system-ui, sans-serif; color: #2C2C2C; padding: 40px; background: white; }
        .inv-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 32px; padding-bottom: 20px; border-bottom: 2px solid #F2D4C8; }
        .inv-brand { display: flex; align-items: center; gap: 10px; }
        .inv-logo { font-size: 2rem; }
        .inv-brand-name { font-size: 1.8rem; font-weight: 700; color: #8B3A3A; letter-spacing: 0.02em; }
        .inv-meta { text-align: right; font-size: 0.88rem; color: #6B6560; line-height: 1.8; }
        .inv-meta h2 { font-size: 1.5rem; color: #2C2C2C; margin-bottom: 4px; letter-spacing: 0.1em; }
        .inv-parties { display: grid; grid-template-columns: 1fr 1fr; gap: 24px; margin-bottom: 28px; padding: 20px; background: #FAF6F0; border-radius: 12px; }
        .inv-party h4 { font-size: 0.72rem; text-transform: uppercase; letter-spacing: 0.1em; color: #6B6560; margin-bottom: 8px; }
        .inv-party p { font-size: 0.88rem; line-height: 1.6; }
        .inv-table { width: 100%; border-collapse: collapse; margin-bottom: 24px; }
        .inv-table thead th { background: #8B3A3A; color: white; padding: 10px 14px; text-align: left; font-size: 0.82rem; font-weight: 500; }
        .inv-table thead th.text-right { text-align: right; }
        .inv-table tbody td { padding: 10px 14px; border-bottom: 1px solid #E8E0D8; font-size: 0.88rem; }
        .inv-table .text-right { text-align: right; }
        .inv-total-row td { padding: 14px 14px; background: #FAF6F0; font-size: 0.95rem; }
        .inv-shipping { margin-bottom: 20px; padding: 20px; border: 1.5px solid #E8E0D8; border-radius: 12px; }
        .inv-shipping h4 { font-size: 0.9rem; margin-bottom: 14px; }
        .inv-shipping-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
        .full-width { grid-column: 1 / -1; }
        .inv-shipping-grid span { font-size: 0.72rem; text-transform: uppercase; letter-spacing: 0.08em; color: #6B6560; display: block; margin-bottom: 2px; }
        .inv-shipping-grid p { font-size: 0.88rem; }
        .inv-notes { font-size: 0.82rem; color: #6B6560; padding: 12px 16px; background: #FFF8F0; border-radius: 8px; margin-bottom: 20px; }
        .inv-footer { text-align: center; font-size: 0.82rem; color: #6B6560; padding-top: 20px; border-top: 1px solid #E8E0D8; line-height: 1.8; }
        @media print { body { padding: 20px; } }
      </style>
    </head>
    <body>
      ${invoiceEl.value.innerHTML}
    </body>
    </html>
  `)
  printWindow.document.close()
  printWindow.focus()
  setTimeout(() => {
    printWindow.print()
    printWindow.close()
  }, 500)

  invoiceEl.value.style.display = 'none'
}
</script>

<style scoped>
.btn-invoice {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  font-size: 0.85rem;
}
</style>
