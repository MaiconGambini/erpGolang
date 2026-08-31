<template>
  <main class="login-page">
    <section class="login-card" aria-labelledby="login-title">
      <div class="brand-lockup" aria-label="goERP">
        <span class="brand-mark" aria-hidden="true">g</span>
        <span class="logo">go<span>ERP</span></span>
      </div>
      <div class="login-intro">
        <h1 id="login-title">Acesse sua conta</h1>
        <p>Entre para gerenciar sua operação</p>
      </div>
      <LoginForm />
    </section>
  </main>
</template>

<script setup lang="ts">
import LoginForm from '@/features/auth/login/ui/LoginForm.vue'
</script>

<style scoped>
.login-page {
  align-items: center;
  background:
    radial-gradient(circle at 50% 0%, var(--color-brand-soft), transparent 42%),
    var(--color-canvas);
  display: flex;
  justify-content: center;
  min-height: 100vh;
  padding: 24px 16px;
}

.login-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 16px;
  box-shadow: var(--shadow-card);
  max-width: 420px;
  overflow: hidden;
  padding: 40px;
  width: 100%;
}


.brand-lockup {
  align-items: center;
  display: flex;
  gap: 10px;
  margin-bottom: 34px;
}

.brand-mark {
  align-items: center;
  background: var(--color-brand);
  border-radius: 9px;
  box-shadow: var(--shadow-soft);
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

.login-intro h1 {
  font-size: 26px;
  letter-spacing: -0.03em;
  line-height: 1.15;
  margin: 0 0 8px;
}

.login-intro p {
  color: var(--color-text-muted);
  margin: 0 0 30px;
}
</style>
