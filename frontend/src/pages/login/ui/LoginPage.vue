<template>
  <main class="login-page">
    <section class="login-card">
      <div class="logo">go<span>ERP</span></div>
      <h1>Acesse sua conta</h1>
      <p>Entre para gerenciar sua operação</p>
      <form class="form">
        <label>Empresa<input placeholder="acme" /></label>
        <label>E-mail<input placeholder="ana.souza@empresa.com.br" /></label>
        <label>Senha<input type="password" placeholder="Sua senha" /></label>
        <AppButton type="submit">Entrar</AppButton>
      </form>
    </section>
  </main>
</template>

<script setup lang="ts">
import AppButton from '@/shared/ui/AppButton.vue'
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

.form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

label {
  color: var(--color-text-secondary);
  display: flex;
  flex-direction: column;
  font-size: 13px;
  gap: 6px;
}

input {
  border: 1px solid #cbd5e1;
  border-radius: var(--radius-md);
  font: inherit;
  padding: 10px 12px;
}
</style>
