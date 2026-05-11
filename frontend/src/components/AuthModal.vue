<template>
  <teleport to="body">
    <div class="modal-backdrop" @click.self="$emit('close')">
      <div class="auth-modal">

        <!-- FORGOT PASSWORD VIEW -->
        <template v-if="view === 'forgot'">
          <div class="modal-header">
            <h2>🔑 Lupa Password</h2>
            <button class="close-btn" @click="$emit('close')">✕</button>
          </div>
          <p class="view-desc">Masukkan email terdaftarmu. Kami akan kirimkan link untuk membuat password baru.</p>
          <div class="auth-form">
            <div class="form-group">
              <label>Email</label>
              <input v-model="forgotEmail" type="email" class="form-input" placeholder="email@contoh.com" />
            </div>
            <div v-if="forgotError" class="auth-error">⚠️ {{ forgotError }}</div>
            <div v-if="forgotSuccess" class="auth-success">✅ {{ forgotSuccess }}</div>
            <button class="btn btn-primary btn-full" :disabled="forgotLoading" @click="submitForgot">
              <span v-if="forgotLoading">⏳ Mengirim...</span>
              <span v-else>Kirim Link Reset</span>
            </button>
          </div>
          <p class="switch-mode">
            <button class="link-btn" @click="view = 'login'">← Kembali ke Login</button>
          </p>
        </template>

        <!-- LOGIN / REGISTER VIEW -->
        <template v-else>
          <div class="modal-header">
            <h2>{{ isLogin ? '👋 Masuk ke Bloome' : '🌸 Daftar Akun Baru' }}</h2>
            <button class="close-btn" @click="$emit('close')">✕</button>
          </div>

          <div class="auth-tabs">
            <button :class="['auth-tab', { active: isLogin }]" @click="switchTab(true)">Masuk</button>
            <button :class="['auth-tab', { active: !isLogin }]" @click="switchTab(false)">Daftar</button>
          </div>

          <form class="auth-form" @submit.prevent="submit">
            <div v-if="!isLogin" class="form-group">
              <label>Nama Lengkap</label>
              <input v-model="form.name" type="text" class="form-input" placeholder="Nama kamu" required />
            </div>

            <div class="form-group">
              <label>Email</label>
              <input v-model="form.email" type="email" class="form-input" placeholder="email@contoh.com" required />
            </div>

            <div v-if="!isLogin" class="form-group">
              <label>No. HP (WhatsApp)</label>
              <input v-model="form.phone" type="tel" class="form-input" placeholder="08xxxxxxxxxx" required />
            </div>

            <div class="form-group">
              <label>Password</label>
              <div class="input-wrap">
                <input
                  v-model="form.password"
                  :type="showPwd ? 'text' : 'password'"
                  class="form-input"
                  placeholder="Minimal 6 karakter"
                  required
                  minlength="6"
                />
                <button type="button" class="pwd-toggle" @click="showPwd = !showPwd">
                  {{ showPwd ? '🙈' : '👁️' }}
                </button>
              </div>
            </div>

            <div v-if="isLogin" class="forgot-wrap">
              <button type="button" class="link-btn" @click="view = 'forgot'">Lupa password?</button>
            </div>

            <div v-if="authStore.error" class="auth-error">⚠️ {{ authStore.error }}</div>

            <button type="submit" class="btn btn-primary btn-full" :disabled="authStore.loading">
              <span v-if="authStore.loading">⏳ Mohon tunggu...</span>
              <span v-else>{{ isLogin ? 'Masuk' : 'Buat Akun' }}</span>
            </button>
          </form>

          <p class="switch-mode">
            {{ isLogin ? 'Belum punya akun?' : 'Sudah punya akun?' }}
            <button class="link-btn" @click="switchTab(!isLogin)">
              {{ isLogin ? 'Daftar sekarang' : 'Masuk di sini' }}
            </button>
          </p>
        </template>

      </div>
    </div>
  </teleport>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { forgotPassword } from '@/services/api'

const props = defineProps({
  initialTab: { type: String, default: 'login' }
})
const emit = defineEmits(['close', 'success'])

const authStore = useAuthStore()
const view = ref('login')
const isLogin = ref(props.initialTab === 'login')
const showPwd = ref(false)
const form = reactive({ name: '', email: '', password: '', phone: '' })

const forgotEmail = ref('')
const forgotError = ref('')
const forgotSuccess = ref('')
const forgotLoading = ref(false)

function switchTab(toLogin) {
  isLogin.value = toLogin
  authStore.error = ''
  form.name = ''
  form.email = ''
  form.password = ''
  form.phone = ''
}

async function submit() {
  let result
  if (isLogin.value) {
    result = await authStore.login(form.email, form.password)
  } else {
    result = await authStore.register(form.name, form.email, form.password, form.phone)
  }
  if (result.success) {
    emit('success')
    emit('close')
  }
}

async function submitForgot() {
  forgotError.value = ''
  forgotSuccess.value = ''
  if (!forgotEmail.value) {
    forgotError.value = 'Masukkan email kamu'
    return
  }
  forgotLoading.value = true
  try {
    await forgotPassword(forgotEmail.value)
    forgotSuccess.value = 'Link reset password sudah dikirim ke email kamu. Cek inbox (atau folder spam).'
    forgotEmail.value = ''
  } catch (e) {
    forgotError.value = e.response?.data?.error || 'Gagal mengirim email. Coba lagi.'
  } finally {
    forgotLoading.value = false
  }
}
</script>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(44, 44, 44, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 500;
  padding: 20px;
}

.auth-modal {
  background: var(--white);
  border-radius: var(--radius);
  padding: 32px;
  width: 100%;
  max-width: 420px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.15);
  animation: slideUp 0.25s ease;
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(24px); }
  to   { opacity: 1; transform: translateY(0); }
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.modal-header h2 {
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--charcoal);
}

.close-btn {
  background: none;
  border: none;
  font-size: 1rem;
  cursor: pointer;
  color: var(--warm-gray);
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  transition: var(--transition);
}
.close-btn:hover { background: var(--cream); }

.auth-tabs {
  display: flex;
  background: var(--cream);
  border-radius: var(--radius-pill);
  padding: 4px;
  margin-bottom: 24px;
}

.auth-tab {
  flex: 1;
  padding: 8px;
  border: none;
  background: transparent;
  border-radius: var(--radius-pill);
  font-size: 0.88rem;
  font-weight: 500;
  cursor: pointer;
  color: var(--warm-gray);
  transition: var(--transition);
}

.auth-tab.active {
  background: var(--white);
  color: var(--deep-rose);
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}

.auth-form { display: flex; flex-direction: column; gap: 16px; }

.form-group { display: flex; flex-direction: column; gap: 6px; }
.form-group label {
  font-size: 0.82rem;
  font-weight: 500;
  color: var(--charcoal);
}

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

.forgot-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: -8px;
}

.auth-error {
  background: #FFF3F3;
  border: 1px solid #FFCCCC;
  color: #C62828;
  padding: 10px 14px;
  border-radius: var(--radius-sm);
  font-size: 0.82rem;
}

.auth-success {
  background: #F0FFF4;
  border: 1px solid #B2DFDB;
  color: #1B5E20;
  padding: 10px 14px;
  border-radius: var(--radius-sm);
  font-size: 0.82rem;
  line-height: 1.5;
}

.view-desc {
  font-size: 0.85rem;
  color: var(--warm-gray);
  margin-bottom: 20px;
  line-height: 1.5;
}

.btn-full { width: 100%; justify-content: center; }

.switch-mode {
  text-align: center;
  font-size: 0.82rem;
  color: var(--warm-gray);
  margin-top: 20px;
}

.link-btn {
  background: none;
  border: none;
  color: var(--deep-rose);
  font-size: 0.82rem;
  font-weight: 500;
  cursor: pointer;
  text-decoration: underline;
  padding: 0;
}
</style>
