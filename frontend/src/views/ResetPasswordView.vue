<template>
  <div class="reset-page">
    <div class="reset-card">
      <div class="reset-logo">🌸 Bloome</div>

      <!-- Loading token check -->
      <div v-if="checking" class="state-box">
        <p>⏳ Memeriksa link...</p>
      </div>

      <!-- Invalid / expired token -->
      <div v-else-if="tokenInvalid" class="state-box error">
        <p class="state-icon">❌</p>
        <h2>Link Tidak Valid</h2>
        <p>Link reset password sudah expired atau sudah digunakan. Silakan minta link baru.</p>
        <router-link to="/" class="btn btn-primary" style="margin-top:16px;">Kembali ke Beranda</router-link>
      </div>

      <!-- Success -->
      <div v-else-if="success" class="state-box success">
        <p class="state-icon">✅</p>
        <h2>Password Berhasil Diubah</h2>
        <p>Silakan login menggunakan password baru kamu.</p>
        <router-link to="/" class="btn btn-primary" style="margin-top:16px;">Login Sekarang</router-link>
      </div>

      <!-- Reset form -->
      <div v-else>
        <h2>Buat Password Baru</h2>
        <p class="desc">Masukkan password baru untuk akunmu.</p>

        <div class="form">
          <div class="form-group">
            <label>Password Baru</label>
            <div class="input-wrap">
              <input
                v-model="password"
                :type="showPwd ? 'text' : 'password'"
                class="form-input"
                placeholder="Minimal 6 karakter"
                minlength="6"
              />
              <button type="button" class="pwd-toggle" @click="showPwd = !showPwd">
                {{ showPwd ? '🙈' : '👁️' }}
              </button>
            </div>
          </div>

          <div class="form-group">
            <label>Konfirmasi Password</label>
            <div class="input-wrap">
              <input
                v-model="confirm"
                :type="showConfirm ? 'text' : 'password'"
                class="form-input"
                placeholder="Ulangi password baru"
              />
              <button type="button" class="pwd-toggle" @click="showConfirm = !showConfirm">
                {{ showConfirm ? '🙈' : '👁️' }}
              </button>
            </div>
          </div>

          <div v-if="error" class="auth-error">⚠️ {{ error }}</div>

          <button class="btn btn-primary btn-full" :disabled="loading" @click="submit">
            <span v-if="loading">⏳ Menyimpan...</span>
            <span v-else>Simpan Password Baru</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { resetPassword } from '@/services/api'

const route = useRoute()
const token = route.query.token

const checking = ref(false)
const tokenInvalid = ref(!token)
const success = ref(false)
const error = ref('')
const loading = ref(false)

const password = ref('')
const confirm = ref('')
const showPwd = ref(false)
const showConfirm = ref(false)

async function submit() {
  error.value = ''
  if (!password.value || password.value.length < 6) {
    error.value = 'Password minimal 6 karakter'
    return
  }
  if (password.value !== confirm.value) {
    error.value = 'Password dan konfirmasi tidak cocok'
    return
  }
  loading.value = true
  try {
    await resetPassword(token, password.value)
    success.value = true
  } catch (e) {
    const msg = e.response?.data?.error || ''
    if (msg.includes('expired') || msg.includes('valid') || msg.includes('digunakan')) {
      tokenInvalid.value = true
    } else {
      error.value = msg || 'Gagal mengubah password. Coba lagi.'
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.reset-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--cream);
  padding: 20px;
}

.reset-card {
  background: var(--white);
  border-radius: var(--radius);
  padding: 40px 32px;
  width: 100%;
  max-width: 420px;
  box-shadow: 0 8px 40px rgba(0,0,0,0.10);
  text-align: center;
}

.reset-logo {
  font-size: 1.4rem;
  font-weight: 700;
  color: var(--deep-rose);
  margin-bottom: 28px;
}

h2 {
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--charcoal);
  margin-bottom: 8px;
}

.desc {
  font-size: 0.85rem;
  color: var(--warm-gray);
  margin-bottom: 24px;
}

.form { display: flex; flex-direction: column; gap: 16px; text-align: left; }

.form-group { display: flex; flex-direction: column; gap: 6px; }
.form-group label { font-size: 0.82rem; font-weight: 500; color: var(--charcoal); }

.input-wrap { position: relative; }
.input-wrap .form-input { width: 100%; padding-right: 44px; }

.pwd-toggle {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1rem;
  padding: 4px;
}

.auth-error {
  background: #FFF3F3;
  border: 1px solid #FFCCCC;
  color: #C62828;
  padding: 10px 14px;
  border-radius: var(--radius-sm);
  font-size: 0.82rem;
}

.btn-full { width: 100%; justify-content: center; }

.state-box { padding: 16px 0; }
.state-box.error { color: #C62828; }
.state-box.success { color: #1B5E20; }
.state-icon { font-size: 2.5rem; margin-bottom: 12px; }
.state-box h2 { margin-bottom: 8px; }
.state-box p { font-size: 0.88rem; color: var(--warm-gray); }
</style>
