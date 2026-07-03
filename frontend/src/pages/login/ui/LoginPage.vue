<template>
  <main class="login-page">
    <section class="login-card">
      <div class="logo">go<span>ERP</span></div>
      <h1>Acesse sua conta</h1>
      <p>Entre para gerenciar sua operação</p>
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
  display: flex;
  justify-content: center;
  min-height: 100vh;
}

.login-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-soft);
  padding: 40px;
  width: 420px;
}

.logo {
  font-size: 22px;
  font-weight: 700;
  margin-bottom: 32px;
}

.logo span {
  color: var(--color-brand);
}

h1 {
  font-size: 24px;
  margin: 0 0 6px;
}

p {
  color: var(--color-text-muted);
  margin: 0 0 28px;
}
</style>
