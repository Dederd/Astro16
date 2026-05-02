# 🌸 Bloome — Bouquet App v2

Aplikasi pemesanan bouquet bunga dengan AI design generator, katalog pre-made, admin panel, dan tracking pengiriman.

---

## ✅ Bug Fixes

| # | Bug | Solusi |
|---|-----|--------|
| 1 | Tombol +/- tangkai tertutup popup | Pindah qty control **ke dalam card** (bukan floating), pakai `position: sticky` bukan `fixed` untuk summary bar |
| 2 | Stem count AI berbeda dengan pilihan user | Backend mengirim `total_stem_count` ke AI, lalu **override paksa** `stem_count` di response agar selalu sinkron |
| 3 | Nama `claude` masih di kode AI | Seluruh AI agent **direfactor ke Groq** (`llama-3.1-8b-instant`) dengan satu fungsi `callGroq()` terpusat |
| 4 | Error pembayaran tidak terlihat | Backend sekarang log detail Midtrans error + response body; frontend tampilkan `details` field dari response |

## 🆕 Fitur Baru

| Fitur | Detail |
|-------|--------|
| 🛍️ Katalog Bouquet | Halaman `/catalog` — user bisa beli langsung tanpa generate AI |
| 🔒 Batas Generate AI | Gratis 2x per session, setelah itu muncul notifikasi |
| 📦 Alamat Pengiriman | Form checkout lengkap dengan penerima, alamat, kota, kode pos |
| 🚚 Pilih Kurir | JNE / J&T / SiCepat dengan estimasi pengiriman |
| 📍 Halaman Tracking | Status pesanan + resi + history pengiriman dari BinderByte API |
| 👤 Halaman Admin | `/admin` — dashboard stats, kelola pesanan, update resi, kelola bunga & katalog |
| 🗄️ Foto dari PostgreSQL | Bunga & katalog disimpan di DB — admin bisa update URL gambar via admin panel |
| 📖 Swagger UI | `http://localhost:8080/swagger/index.html` |

---

## 🚀 Setup & Menjalankan

### 1. Clone dan konfigurasi

```bash
# Copy env files
cp backend/.env.example backend/.env
# Edit backend/.env — isi GROQ_API_KEY dan MIDTRANS_SERVER_KEY
```

### 2. Jalankan dengan Docker Compose (Recommended)

```bash
docker-compose up -d
```

Akses:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html
- Admin: http://localhost:5173/admin

### 3. Jalankan Manual (Development)

**Backend:**
```bash
cd backend

# Install swag (sekali saja)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger docs
swag init -g cmd/main.go --output docs

# Install dependencies & run
go mod tidy
go run cmd/main.go
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

**PostgreSQL:** Pastikan PostgreSQL berjalan di port 5432. Data bunga & katalog akan di-seed otomatis saat pertama kali dijalankan.

---

## 🔑 Konfigurasi Penting

### backend/.env

```env
# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=bouquet_db

# Groq AI (https://console.groq.com — gratis)
GROQ_API_KEY=gsk_xxx

# Midtrans (https://dashboard.midtrans.com — gunakan Sandbox untuk testing)
MIDTRANS_SERVER_KEY=SB-Mid-server-xxx
MIDTRANS_IS_PRODUCTION=false

# BinderByte tracking (opsional, https://binderbyte.com)
BINDERBYTE_API_KEY=

# URL frontend untuk callback Midtrans
FRONTEND_URL=http://localhost:5173
```

### frontend/index.html

Ganti `data-client-key` di script Midtrans Snap dengan client key kamu:

```html
<script src="https://app.sandbox.midtrans.com/snap/snap.js"
  data-client-key="Mid-client-xxx"></script>
```

---

## 📋 API Endpoints

### Public
| Method | Path | Deskripsi |
|--------|------|-----------|
| GET | `/api/v1/bouquet-types` | Tipe bouquet |
| GET | `/api/v1/flowers` | Daftar bunga dari DB |
| GET | `/api/v1/catalog` | Katalog pre-made |
| POST | `/api/v1/agent/verify-selection` | AI Agent 1 |
| POST | `/api/v1/agent/generate-bouquet` | AI Agent 2 |
| GET | `/api/v1/agent/generate-status` | Cek kuota generate |
| POST | `/api/v1/orders` | Buat order |
| GET | `/api/v1/orders/:id` | Detail order |
| GET | `/api/v1/orders/:id/tracking` | Info pengiriman |
| POST | `/api/v1/payment/token` | Token Midtrans |
| POST | `/api/v1/payment/notification` | Webhook Midtrans |

### Admin (Header: `X-Admin-Key: admin-bouquet-2024`)
| Method | Path | Deskripsi |
|--------|------|-----------|
| GET | `/api/v1/admin/stats` | Dashboard statistik |
| GET | `/api/v1/admin/orders` | Semua pesanan |
| PUT | `/api/v1/admin/orders/:id` | Update status & resi |
| GET | `/api/v1/admin/flowers` | Semua bunga |
| PUT | `/api/v1/admin/flowers/:id` | Update bunga |
| GET | `/api/v1/admin/catalog` | Semua katalog |
| POST | `/api/v1/admin/catalog` | Tambah katalog |
| PUT | `/api/v1/admin/catalog/:id` | Update katalog |
| DELETE | `/api/v1/admin/catalog/:id` | Hapus katalog |

---

## 🏗️ Struktur Proyek

```
bouquet-app-v2/
├── backend/
│   ├── cmd/main.go               ← Entry point + routing
│   ├── internal/
│   │   ├── database/db.go        ← PostgreSQL + seeding
│   │   ├── handlers/handlers.go  ← Semua HTTP handlers
│   │   ├── middleware/cors.go    ← CORS
│   │   ├── models/models.go      ← Struct + DB models
│   │   └── services/
│   │       ├── ai_agent.go       ← Groq AI agents
│   │       ├── data.go           ← Data dari DB
│   │       └── payment.go        ← Midtrans + tracking
│   ├── docs/                     ← Swagger docs (generated)
│   ├── .env                      ← Konfigurasi
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── steps/            ← Step1-4 components
│   │   │   ├── FlowerCard.vue    ← Card bunga (qty fix)
│   │   │   ├── DesignCard.vue    ← Card desain AI
│   │   │   └── CatalogCard.vue   ← Card katalog
│   │   ├── views/
│   │   │   ├── Home.vue
│   │   │   ├── Order.vue
│   │   │   ├── CatalogView.vue   ← Halaman katalog
│   │   │   ├── PaymentFinishView.vue  ← Status + tracking
│   │   │   └── AdminView.vue     ← Admin panel
│   │   ├── stores/order.js       ← Pinia store
│   │   ├── services/api.js       ← Axios API calls
│   │   └── router/index.js
│   └── .env
└── docker-compose.yml
```
