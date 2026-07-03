<template>
  <div class="shell">
    <aside class="sidebar">
      <div class="logo">go<span>ERP</span></div>
      <RouterLink to="/">Dashboard</RouterLink>
      <RouterLink to="/customers">Clientes</RouterLink>
      <RouterLink to="/products">Produtos</RouterLink>
    </aside>
    <div class="main">
      <header class="topbar">
        <span class="user">{{ session.user?.name ?? 'Usuário' }}</span>
        <LogoutButton />
      </header>
      <main class="content">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useSessionStore } from '@/entities/session/model/session.store'
import LogoutButton from '@/features/auth/logout/ui/LogoutButton.vue'

const session = useSessionStore()
</script>

<style scoped>
.shell {
  display: flex;
  min-height: 100vh;
}

.sidebar {
  background: var(--color-surface);
  border-right: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 20px 12px;
  width: 240px;
}

.logo {
  font-size: 20px;
  font-weight: 700;
  margin: 0 8px 20px;
}

.logo span {
  color: var(--color-brand);
}

a {
  border-radius: var(--radius-sm);
  color: var(--color-text-secondary);
  padding: 8px 12px;
  text-decoration: none;
}

a.router-link-active {
  background: #eff6ff;
  color: var(--color-brand-hover);
}

.main {
  display: flex;
  flex: 1;
  flex-direction: column;
}

.topbar {
  align-items: center;
  border-bottom: 1px solid var(--color-border);
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  padding: 12px 32px;
}

.user {
  color: var(--color-text-secondary);
  font-size: 14px;
}

.content {
  flex: 1;
  padding: 28px 32px;
}
</style>
