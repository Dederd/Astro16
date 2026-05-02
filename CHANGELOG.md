# Changelog — Update Auth & My Orders

## Bug Fix
### ❌ Badge "★ Recommended" Muncul Double
**Root cause:** Ada dua sumber badge yang aktif bersamaan:
1. `frontend/src/assets/main.css` — `.flower-card.recommended::after { content: '★ Recommended' }`
2. `frontend/src/components/FlowerCard.vue` — `<div v-if="isRecommended" class="rec-badge">`

**Fix:** Hapus blok `::after` di `main.css` (line 240-252). Badge sekarang hanya dari komponen Vue.

---

## Fitur Baru: Auth (Login / Register)

### Frontend
- **`src/stores/auth.js`** — Pinia store baru: `login()`, `register()`, `logout()`, persist token ke localStorage
- **`src/components/AuthModal.vue`** — Modal login/register dengan tab switcher, validasi form, show/hide password
- **`src/App.vue`** — Navbar diupdate:
  - Jika belum login: tampil tombol **"Masuk"** → buka AuthModal
  - Jika sudah login: tampil avatar + nama → dropdown **"Pesanan Saya"** dan **"Keluar"**
- **`src/services/api.js`** — Ditambah `export default api`, endpoint auth (`/auth/register`, `/auth/login`, `/auth/me`), dan `getUserOrders`

### Backend
- **`internal/models/models.go`** — Tambah `UserDB`, `RegisterRequest`, `LoginRequest`, `AuthResponse`, dan field `UserID *uint` di `OrderDB`
- **`internal/handlers/auth.go`** — Handler baru: `Register`, `Login`, `GetMe`, `GetUserOrders`, `AuthMiddleware`, `OptionalAuthMiddleware`
- **`internal/database/db.go`** — Auto-migrate `UserDB`
- **`cmd/main.go`** — Route baru:
  - `POST /api/v1/auth/register`
  - `POST /api/v1/auth/login`
  - `GET  /api/v1/auth/me`
  - `GET  /api/v1/user/orders`
  - `POST /api/v1/orders` — sekarang pakai `OptionalAuthMiddleware` (link ke user jika login)
- **`go.mod`** — Tambah `github.com/golang-jwt/jwt/v5 v5.2.1` dan promote `golang.org/x/crypto` ke direct dependency

---

## Fitur Baru: My Orders (Daftar & Tracking Pesanan)

### Frontend
- **`src/views/MyOrdersView.vue`** — Halaman baru `/my-orders`:
  - List semua pesanan user (urutkan terbaru dulu)
  - Klik kartu → expand detail + **tracking visual** (step-by-step: Pesanan Diterima → Pembayaran → Dirangkai → Pengiriman → Terkirim)
  - Badge status berwarna (pending/paid/processing/shipped/delivered/cancelled)
  - Info lengkap: nama, no HP, alamat, kurir, ID pembayaran
- **`src/router/index.js`** — Route `/my-orders` dengan guard `requiresAuth: true` (redirect jika belum login)

---

## Guard Checkout (Wajib Login)

- **`src/components/steps/Step4Checkout.vue`**:
  - Banner kuning muncul jika belum login: *"Kamu perlu masuk atau daftar untuk menyelesaikan pembelian"*
  - Klik **"Bayar Sekarang"** saat belum login → buka AuthModal otomatis
  - Setelah login berhasil → lanjut proses pembayaran normal

---

## Setup Backend (Tambahan)

```bash
# Install dependencies baru
cd backend
go mod tidy

# Variabel .env baru (opsional, ada default)
JWT_SECRET=your-secret-key-here
```

JWT token berlaku 30 hari. Password di-hash dengan bcrypt (cost 10).
