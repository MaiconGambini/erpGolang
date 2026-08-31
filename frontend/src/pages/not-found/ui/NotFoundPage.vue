<template>
  <main class="not-found" aria-labelledby="not-found-title">
    <div class="brand-lockup" aria-label="goERP">
      <span class="brand-mark" aria-hidden="true">g</span>
      <span class="logo">go<span>ERP</span></span>
    </div>
    <span class="not-found-code">404</span>
    <h1 id="not-found-title">Página não encontrada</h1>
    <p>O endereço que você acessou não existe ou foi movido.</p>
    <RouterLink to="/" class="home-link">Voltar ao dashboard</RouterLink>
  </main>
</template>

<style scoped>
.not-found {
  align-items: center;
  background: var(--color-canvas);
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 100vh;
  padding: 24px 16px;
  text-align: center;
}

.brand-lockup {
  align-items: center;
  display: flex;
  gap: 10px;
  margin-bottom: 44px;
}

.brand-mark {
  align-items: center;
  background: var(--color-brand);
  border-radius: 9px;
  color: var(--color-text-on-brand);
  display: inline-flex;
  font-size: 18px;
  font-weight: 800;
  height: 32px;
  justify-content: center;
  letter-spacing: -0.06em;
  width: 32px;
}

.logo {
  font-size: 22px;
  font-weight: 750;
  letter-spacing: -0.04em;
}

.logo span {
  color: var(--color-brand);
}

.not-found-code {
  color: var(--color-brand);
  font-size: clamp(64px, 12vw, 112px);
  font-weight: 750;
  letter-spacing: -0.08em;
  line-height: 0.9;
}

h1 {
  font-size: 26px;
  letter-spacing: -0.03em;
  line-height: 1.15;
  margin: 24px 0 8px;
}

p {
  color: var(--color-text-muted);
  margin: 0 0 24px;
  max-width: 42ch;
}

.home-link {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-brand);
  font-weight: 650;
  min-height: 40px;
  padding: 9px 15px;
  text-decoration: none;
  transition: background-color 0.15s ease, border-color 0.15s ease;
}

.home-link:hover {
  background: var(--color-brand-soft);
  border-color: var(--color-brand);
}
</style>
