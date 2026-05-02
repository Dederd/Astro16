<template>
  <teleport to="body">
    <div class="modal-backdrop" @click.self="$emit('close')">
      <div class="auth-modal">
        <!-- Header -->
        <div class="modal-header">
          <h2>{{ isLogin ? '👋 Masuk ke Bloome' : '🌸 Daftar Akun Baru' }}</h2>
          <button class="close-btn" @click="$emit('close')">✕</button>
        </div>

        <!-- Tab switcher -->
        <div class="auth-tabs">
          <button :class="['auth-tab', { active: isLogin }]" @click="isLogin = true">Masuk</button>
          <button :class="['auth-tab', { active: !isLogin }]" @click="isLogin = false">Daftar</button>
        </div>

        <!-- Form -->
        <form class="auth-form" @submit.prevent="submit">
          <!-- Register: Nama -->
          <div v-if="!isLogin" class="form-group">
            <label>Nama Lengkap</label>
            <input v-model="form.name" type="text" class="form-input" placeholder="Nama kamu" required />
          </div>

          <!-- Email -->
          <div class="form-group">
            <label>Email</label>
            <input v-model="form.email" type="email" class="form-input" placeholder="email@contoh.com" required />
          </div>

          <!-- Register: No HP -->
          <div v-if="!isLogin" class="form-group">
            <label>No. HP (WhatsApp)</label>
            <input v-model="form.phone" type="tel" class="form-input" placeholder="08xxxxxxxxxx" required />
          </div>

          <!-- Password -->
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

          <!-- Error -->
          <div v-if="authStore.error" class="auth-error">
            ⚠️ {{ authStore.error }}
          </div>

          <!-- Submit -->
          <button type="submit" class="btn btn-primary btn-full" :disabled="authStore.loading">
            <span v-if="authStore.loading">⏳ Mohon tunggu...</span>
            <span v-else>{{ isLogin ? 'Masuk' : 'Buat Akun' }}</span>
          </button>
        </form>

        <!-- Switch mode -->
        <p class="switch-mode">
          {{ isLogin ? 'Belum punya akun?' : 'Sudah punya akun?' }}
          <button class="link-btn" @click="isLogin = !isLogin">
            {{ isLogin ? 'Daftar sekarang' : 'Masuk di sini' }}
          </button>
        </p>
      </div>
    </div>
  </teleport>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'

const props = defineProps({
  initialTab: { type: String, default: 'login' } // 'login' | 'register'
})
const emit = defineEmits(['close', 'success'])

const authStore = useAuthStore()
const isLogin = ref(props.initialTab === 'login')
const showPwd = ref(false)
const form = reactive({ name: '', email: '', password: '', phone: '' })

watch(isLogin, () => {
  authStore.error = ''
  form.name = ''
  form.email = ''
  form.password = ''
  form.phone = ''
})

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

.auth-error {
  background: #FFF3F3;
  border: 1px solid #FFCCCC;
  color: #C62828;
  padding: 10px 14px;
  border-radius: var(--radius-sm);
  font-size: 0.82rem;
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
