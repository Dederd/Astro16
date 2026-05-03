<template>
  <div id="app-root">
    <!-- Navbar — hidden on admin page for cleaner layout -->
    <nav v-if="!isAdminRoute" class="navbar">
      <div class="container navbar-inner">
        <router-link to="/" class="logo">
          <span class="logo-icon">🌸</span>
          <span class="logo-text">Bloome</span>
        </router-link>
        <div class="nav-links">
          <router-link to="/" class="nav-link hide-mobile">Beranda</router-link>
          <router-link to="/catalog" class="nav-link catalog-link">Katalog</router-link>
          <router-link to="/order" class="btn btn-primary btn-nav">Buat Bouquet</router-link>

          <!-- Auth area -->
          <div v-if="authStore.isLoggedIn" class="user-menu" ref="userMenuRef">
            <button class="user-btn" @click="showUserMenu = !showUserMenu">
              <span class="user-avatar">{{ initials }}</span>
              <span class="user-name">{{ authStore.user?.name?.split(' ')[0] }}</span>
              <span class="chevron">▾</span>
            </button>
            <transition name="dropdown">
              <div v-if="showUserMenu" class="user-dropdown">
                <router-link to="/my-orders" class="dropdown-item" @click="showUserMenu = false">
                  📦 Pesanan Saya
                </router-link>
                <hr class="dropdown-hr" />
                <button class="dropdown-item logout" @click="logout">
                  🚪 Keluar
                </button>
              </div>
            </transition>
          </div>

          <button v-else class="btn btn-outline btn-nav" @click="openAuth('login')">
            Masuk
          </button>
        </div>
      </div>
    </nav>

    <!-- Main content -->
    <main :class="{ 'main-admin': isAdminRoute }">
      <router-view v-slot="{ Component }">
        <transition name="page" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <!-- Footer — hidden on admin -->
    <footer v-if="!isAdminRoute" class="footer">
      <div class="container footer-inner">
        <div class="footer-brand">
          <span class="logo-icon">🌸</span>
          <span class="logo-text">Bloome</span>
        </div>
        <p class="footer-copy">© 2025 Bloome · Dibuat dengan cinta untuk setiap momen spesialmu</p>
      </div>
    </footer>

    <!-- Auth Modal -->
    <AuthModal
      v-if="showAuthModal"
      :initial-tab="authTab"
      @close="showAuthModal = false"
      @success="onAuthSuccess"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, provide } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import AuthModal from '@/components/AuthModal.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const isAdminRoute = computed(() => route.path.startsWith('/admin'))

// Auth modal
const showAuthModal = ref(false)
const authTab = ref('login')

function openAuth(tab = 'login') {
  authTab.value = tab
  showAuthModal.value = true
}

function onAuthSuccess() {
  if (route.query.redirect) {
    router.push(route.query.redirect)
  }
}

// User dropdown
const showUserMenu = ref(false)
const userMenuRef = ref(null)

const initials = computed(() => {
  const name = authStore.user?.name || ''
  return name.split(' ').map(n => n[0]).slice(0, 2).join('').toUpperCase()
})

function logout() {
  authStore.logout()
  showUserMenu.value = false
  router.push('/')
}

function handleOutsideClick(e) {
  if (userMenuRef.value && !userMenuRef.value.contains(e.target)) {
    showUserMenu.value = false
  }
}
onMounted(() => document.addEventListener('click', handleOutsideClick))
onUnmounted(() => document.removeEventListener('click', handleOutsideClick))

provide('openAuth', openAuth)
</script>

<style scoped>
.navbar {
  position: sticky;
  top: 0;
  z-index: 100;
  background: rgba(250,246,240,0.92);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--light-gray);
}
.navbar-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
}
.logo { display: flex; align-items: center; gap: 10px; text-decoration: none; }
.logo-icon { font-size: 1.4rem; }
.logo-text {
  font-family: var(--font-display);
  font-size: 1.6rem;
  font-weight: 600;
  color: var(--deep-rose);
  letter-spacing: 0.02em;
}
.nav-links { display: flex; align-items: center; gap: 16px; }
.nav-link {
  font-size: 0.88rem;
  color: var(--warm-gray);
  text-decoration: none;
  font-weight: 500;
  transition: var(--transition);
  padding: 6px 0;
  border-bottom: 2px solid transparent;
}
.nav-link:hover, .nav-link.router-link-active { color: var(--deep-rose); border-bottom-color: var(--deep-rose); }
.btn-nav { padding: 10px 22px; font-size: 0.85rem; }

/* User menu */
.user-menu { position: relative; }
.user-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  background: var(--cream);
  border: 1.5px solid var(--blush);
  border-radius: var(--radius-pill);
  padding: 6px 14px 6px 6px;
  cursor: pointer;
  transition: var(--transition);
}
.user-btn:hover { border-color: var(--rose); }
.user-avatar {
  width: 28px; height: 28px;
  border-radius: 50%;
  background: var(--deep-rose);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.7rem;
  font-weight: 700;
}
.user-name { font-size: 0.82rem; font-weight: 500; color: var(--charcoal); }
.chevron { font-size: 0.6rem; color: var(--warm-gray); }
.user-dropdown {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  background: var(--white);
  border: 1.5px solid var(--light-gray);
  border-radius: var(--radius-sm);
  box-shadow: 0 8px 32px rgba(0,0,0,0.12);
  min-width: 180px;
  z-index: 200;
  overflow: hidden;
}
.dropdown-item {
  display: block;
  width: 100%;
  padding: 12px 16px;
  font-size: 0.85rem;
  color: var(--charcoal);
  text-decoration: none;
  background: none;
  border: none;
  text-align: left;
  cursor: pointer;
  transition: var(--transition);
}
.dropdown-item:hover { background: var(--cream); }
.dropdown-item.logout { color: #C62828; }
.dropdown-hr { border: none; border-top: 1px solid var(--light-gray); margin: 0; }

.dropdown-enter-active, .dropdown-leave-active { transition: all 0.18s ease; }
.dropdown-enter-from, .dropdown-leave-to { opacity: 0; transform: translateY(-8px); }

main { min-height: calc(100vh - 64px - 120px); padding: 0 24px; }
main.main-admin { padding: 0; }
.footer { margin-top: 80px; background: var(--charcoal); color: var(--cream); padding: 40px 0; }
.footer-inner { display: flex; flex-direction: column; align-items: center; gap: 16px; }
.footer .logo-text { color: var(--blush); font-size: 1.3rem; }
.footer-copy { font-size: 0.82rem; color: rgba(250,246,240,0.5); text-align: center; }

.page-enter-active, .page-leave-active { transition: opacity 0.3s ease, transform 0.3s ease; }
.page-enter-from { opacity: 0; transform: translateY(12px); }
.page-leave-to { opacity: 0; transform: translateY(-8px); }

/* ── Mobile navbar responsive ── */
@media (max-width: 640px) {
  .navbar-inner { height: 52px; }
  .logo-text { font-size: 1.2rem; }

  .nav-links { gap: 6px; }
  .hide-mobile { display: none; }  /* sembunyikan Beranda */

  /* Katalog tetap tampil di mobile */
  .catalog-link {
    display: inline-flex !important;
    font-size: 0.82rem;
    padding: 6px 10px;
    border: 1px solid var(--blush);
    border-radius: var(--radius-pill);
    color: var(--deep-rose);
  }
  .catalog-link.router-link-active {
    background: var(--cream);
    border-color: var(--rose);
  }

  .btn-nav {
    padding: 8px 12px;
    font-size: 0.78rem;
  }

  .user-name { display: none; }
  .user-btn { padding: 5px 10px 5px 5px; }

  main { padding: 0; min-height: calc(100vh - 52px - 80px); }
  .footer { margin-top: 40px; padding: 28px 0; }
}

</style>
