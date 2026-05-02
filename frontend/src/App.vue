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
          <router-link to="/" class="nav-link">Beranda</router-link>
          <router-link to="/catalog" class="nav-link">Katalog</router-link>
          <router-link to="/order" class="btn btn-primary btn-nav">Buat Bouquet</router-link>
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
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const isAdminRoute = computed(() => route.path.startsWith('/admin'))
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

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
}

.logo-icon { font-size: 1.4rem; }

.logo-text {
  font-family: var(--font-display);
  font-size: 1.6rem;
  font-weight: 600;
  color: var(--deep-rose);
  letter-spacing: 0.02em;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 20px;
}

.nav-link {
  font-size: 0.88rem;
  color: var(--warm-gray);
  text-decoration: none;
  font-weight: 500;
  transition: var(--transition);
  padding: 6px 0;
  border-bottom: 2px solid transparent;
}
.nav-link:hover,
.nav-link.router-link-active {
  color: var(--deep-rose);
  border-bottom-color: var(--deep-rose);
}

.btn-nav {
  padding: 10px 22px;
  font-size: 0.85rem;
}

main { min-height: calc(100vh - 64px - 120px); padding: 0 24px; }
main.main-admin { padding: 0; }

.footer {
  margin-top: 80px;
  background: var(--charcoal);
  color: var(--cream);
  padding: 40px 0;
}

.footer-inner {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.footer .logo-text { color: var(--blush); font-size: 1.3rem; }

.footer-copy {
  font-size: 0.82rem;
  color: rgba(250,246,240,0.5);
  text-align: center;
}

/* Page transition */
.page-enter-active, .page-leave-active { transition: opacity 0.3s ease, transform 0.3s ease; }
.page-enter-from { opacity: 0; transform: translateY(12px); }
.page-leave-to { opacity: 0; transform: translateY(-8px); }
</style>
