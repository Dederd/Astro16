<template>
  <div class="home">
    <!-- Hero -->
    <section class="hero">
      <div class="hero-bg">
        <div class="petal petal-1">🌸</div>
        <div class="petal petal-2">🌺</div>
        <div class="petal petal-3">🌼</div>
        <div class="petal petal-4">🌷</div>
        <div class="petal petal-5">🌹</div>
      </div>
      <div class="container hero-content">
        <div class="hero-badge badge badge-gold">✦ AI-Powered Bouquet Designer</div>
        <h1 class="hero-title">
          Ciptakan Bouquet<br>
          <em>Impianmu Sendiri</em>
        </h1>
        <p class="hero-desc">
          Pilih bungamu, biarkan AI kami merancang bouquet yang sempurna untukmu.
          Setiap momen layak mendapat bunga yang diciptakan dengan penuh cinta.
        </p>
        <router-link to="/order" class="btn btn-primary hero-cta">
          Mulai Merancang
          <span>→</span>
        </router-link>
      </div>
    </section>

    <!-- How it works -->
    <section class="how-it-works">
      <div class="container">
        <div class="section-header">
          <span class="section-label">Cara Kerja</span>
          <h2 class="section-title">Semudah 4 Langkah</h2>
        </div>
        <div class="steps-grid">
          <div class="step-item" v-for="(step, i) in steps" :key="i" :style="`--delay: ${i * 0.1}s`">
            <div class="step-icon">{{ step.icon }}</div>
            <div class="step-num">0{{ i + 1 }}</div>
            <h3>{{ step.title }}</h3>
            <p>{{ step.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Occasions -->
    <section class="occasions">
      <div class="container">
        <div class="section-header">
          <span class="section-label">Untuk Semua Momen</span>
          <h2 class="section-title">Pilih Acara Spesialmu</h2>
        </div>
        <div class="occasions-grid">
          <router-link
            to="/order"
            class="occasion-card"
            v-for="occ in occasions" :key="occ.name"
            :style="`--theme: ${occ.color}`"
          >
            <span class="occ-icon">{{ occ.icon }}</span>
            <span class="occ-name">{{ occ.name }}</span>
          </router-link>
        </div>
      </div>
    </section>

    <!-- Features -->
    <section class="features">
      <div class="container">
        <div class="features-grid">
          <div class="feature-card" v-for="f in features" :key="f.title">
            <div class="feature-icon">{{ f.icon }}</div>
            <h3>{{ f.title }}</h3>
            <p>{{ f.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section class="cta-section">
      <div class="container">
        <div class="cta-card">
          <div class="cta-petals">🌸 🌺 🌼 🌷 🌹</div>
          <h2>Siap Merangkai Bouquet Impianmu?</h2>
          <p>Dengan teknologi AI kami, setiap bouquet adalah karya seni yang unik.</p>
          <router-link to="/order" class="btn btn-primary">
            Buat Bouquet Sekarang ✦
          </router-link>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
const steps = [
  { icon: '🎯', title: 'Pilih Momen', desc: 'Tentukan acara apa yang ingin kamu rayakan — wisuda, pernikahan, ulang tahun, dan lainnya.' },
  { icon: '🤖', title: 'AI Merekomendasikan', desc: 'Agent AI kami akan memverifikasi pilihanmu dan merekomendasikan bunga-bunga terbaik.' },
  { icon: '🌺', title: 'Pilih Bungamu', desc: 'Pilih dari koleksi bunga premium kami. Kamu bisa memilih lebih dari satu jenis bunga.' },
  { icon: '✨', title: 'Generate Desain', desc: 'AI kami menciptakan beberapa pilihan desain bouquet berdasarkan bunga-bunga pilihanmu.' },
]

const occasions = [
  { icon: '🎓', name: 'Wisuda', color: '#7C4DFF' },
  { icon: '💍', name: 'Pernikahan', color: '#E91E8C' },
  { icon: '🎂', name: 'Ulang Tahun', color: '#FF6B35' },
  { icon: '💑', name: 'Anniversary', color: '#C62828' },
  { icon: '🏆', name: 'Selamat', color: '#F9A825' },
  { icon: '❤️', name: 'Valentine', color: '#D32F2F' },
  { icon: '🌸', name: 'Hari Ibu', color: '#AD1457' },
  { icon: '🕊️', name: 'Dukacita', color: '#546E7A' },
]

const features = [
  { icon: '🤖', title: 'Dirancang oleh AI', desc: 'Teknologi AI menghasilkan desain bouquet yang unik dan personal sesuai pilihanmu.' },
  { icon: '🌿', title: 'Bunga Segar Premium', desc: 'Kami hanya menggunakan bunga-bunga segar berkualitas tinggi dari kebun terpilih.' },
  { icon: '💳', title: 'Pembayaran Aman', desc: 'Bayar dengan tenang menggunakan berbagai metode pembayaran yang aman via Midtrans.' },
  { icon: '🚚', title: 'Pengiriman Cepat', desc: 'Bouquet siap diantar ke alamatmu tepat waktu untuk momen yang tak terlupakan.' },
]
</script>

<style scoped>
/* Hero */
.hero {
  position: relative;
  min-height: 80vh;
  display: flex;
  align-items: center;
  overflow: hidden;
  background: linear-gradient(160deg, var(--cream) 0%, #FDF0EC 60%, var(--blush) 100%);
}

.hero-bg {
  position: absolute;
  inset: 0;
  pointer-events: none;
  overflow: hidden;
}

.petal {
  position: absolute;
  font-size: 2.5rem;
  opacity: 0.15;
  animation: float 8s ease-in-out infinite;
}
.petal-1 { top: 10%; left: 5%; animation-delay: 0s; font-size: 4rem; }
.petal-2 { top: 60%; left: 2%; animation-delay: 1.5s; }
.petal-3 { top: 20%; right: 8%; animation-delay: 0.8s; font-size: 3rem; }
.petal-4 { bottom: 15%; right: 5%; animation-delay: 2s; font-size: 3.5rem; }
.petal-5 { top: 75%; left: 20%; animation-delay: 0.4s; }

@keyframes float {
  0%, 100% { transform: translateY(0) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(10deg); }
}

.hero-content {
  position: relative;
  z-index: 2;
  max-width: 680px;
  padding-top: 80px;
  padding-bottom: 80px;
}

.hero-badge {
  margin-bottom: 24px;
  font-size: 0.75rem;
}

.hero-title {
  font-size: clamp(2.8rem, 6vw, 4.5rem);
  line-height: 1.1;
  color: var(--charcoal);
  margin-bottom: 24px;
}

.hero-title em {
  color: var(--deep-rose);
  font-style: italic;
}

.hero-desc {
  font-size: 1.05rem;
  color: var(--warm-gray);
  max-width: 520px;
  margin-bottom: 40px;
  line-height: 1.7;
}

.hero-cta {
  font-size: 1rem;
  padding: 16px 36px;
  box-shadow: 0 8px 24px rgba(139,58,58,0.25);
}

/* Section commons */
section { padding: 80px 0; }

.section-header {
  text-align: center;
  margin-bottom: 56px;
}

.section-label {
  display: inline-block;
  font-size: 0.75rem;
  letter-spacing: 0.15em;
  text-transform: uppercase;
  color: var(--rose);
  margin-bottom: 12px;
  font-weight: 500;
}

.section-title {
  font-size: clamp(2rem, 4vw, 2.8rem);
  color: var(--charcoal);
}

/* Steps */
.how-it-works { background: var(--white); }

.steps-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 32px;
}

.step-item {
  text-align: center;
  animation: fadeInUp 0.6s ease both;
  animation-delay: var(--delay, 0s);
}

.step-icon { font-size: 2.5rem; margin-bottom: 12px; }

.step-num {
  font-family: var(--font-display);
  font-size: 3rem;
  color: var(--blush);
  line-height: 1;
  margin-bottom: 8px;
}

.step-item h3 {
  font-size: 1.2rem;
  margin-bottom: 8px;
  color: var(--charcoal);
}

.step-item p {
  font-size: 0.9rem;
  color: var(--warm-gray);
  line-height: 1.6;
}

/* Occasions */
.occasions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 16px;
}

.occasion-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 28px 16px;
  background: var(--white);
  border-radius: var(--radius);
  border: 1px solid var(--light-gray);
  transition: var(--transition);
  text-decoration: none;
}

.occasion-card:hover {
  border-color: var(--theme, var(--rose));
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.1);
}

.occ-icon { font-size: 2rem; }

.occ-name {
  font-size: 0.85rem;
  font-weight: 500;
  color: var(--charcoal);
  text-align: center;
}

/* Features */
.features { background: linear-gradient(135deg, #FDF0EC 0%, var(--cream) 100%); }

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 24px;
}

.feature-card {
  background: var(--white);
  border-radius: var(--radius);
  padding: 32px 24px;
  border: 1px solid var(--light-gray);
  transition: var(--transition);
}

.feature-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
}

.feature-icon { font-size: 2rem; margin-bottom: 16px; }

.feature-card h3 {
  font-size: 1.1rem;
  margin-bottom: 8px;
}

.feature-card p {
  font-size: 0.88rem;
  color: var(--warm-gray);
  line-height: 1.6;
}

/* CTA */
.cta-section { background: var(--white); }

.cta-card {
  background: linear-gradient(135deg, var(--deep-rose) 0%, #A0434A 100%);
  border-radius: 24px;
  padding: 64px 48px;
  text-align: center;
  color: var(--white);
}

.cta-petals { font-size: 1.8rem; margin-bottom: 24px; letter-spacing: 8px; }

.cta-card h2 {
  font-size: clamp(1.8rem, 4vw, 2.8rem);
  margin-bottom: 16px;
  color: var(--white);
}

.cta-card p {
  font-size: 1rem;
  color: rgba(255,255,255,0.8);
  margin-bottom: 36px;
}

.cta-card .btn-primary {
  background: var(--white);
  color: var(--deep-rose);
  border-color: var(--white);
  font-size: 1rem;
  padding: 16px 36px;
}

.cta-card .btn-primary:hover {
  background: var(--cream);
  border-color: var(--cream);
}

/* ── Mobile responsive ── */
@media (max-width: 640px) {
  .section { padding: 48px 0; }
  .section-title { font-size: clamp(1.6rem, 6vw, 2rem); }
  .hero { padding: 48px 0 32px; }
  .hero-badge { font-size: 0.7rem; padding: 6px 14px; }

  .steps-grid { grid-template-columns: 1fr !important; gap: 12px !important; }
  .occasions-grid { grid-template-columns: 1fr 1fr !important; gap: 10px !important; }
  .features-grid { grid-template-columns: 1fr !important; }

  .cta-card { padding: 40px 24px; border-radius: 16px; }
  .cta-card h2 { font-size: 1.6rem; }
  .cta-card .btn-primary { padding: 14px 28px; }
}
</style>
