<template>
  <div class="shell">
    <aside class="sidebar">
      <div class="logo">go<span>ERP</span></div>
      <RouterLink to="/">Dashboard</RouterLink>
      <RouterLink to="/customers">Clientes</RouterLink>
      <RouterLink to="/products">Produtos</RouterLink>
      <RouterLink to="/suppliers">Fornecedores</RouterLink>
      <RouterLink to="/sales">Vendas</RouterLink>
      <template v-if="isAdminUser">
        <RouterLink to="/users">Usuários</RouterLink>
        <RouterLink to="/audit">Auditoria</RouterLink>
      </template>
    </aside>
    <div class="main">
      <header class="topbar">
        <span class="user">{{ session.user?.name ?? 'Usuário' }}</span>
        <button
          type="button"
          class="theme-toggle"
          :aria-label="isDark ? 'Mudar para tema claro' : 'Mudar para tema escuro'"
          :title="isDark ? 'Tema claro' : 'Tema escuro'"
          @click="toggle"
        >
          <Moon v-if="!isDark" :size="17" :stroke-width="2" aria-hidden="true" />
          <Sun v-else :size="17" :stroke-width="2" aria-hidden="true" />
        </button>
        <LogoutButton />
      </header>
      <main class="content">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Moon, Sun } from 'lucide-vue-next'
import { useSessionStore } from '@/entities/session/model/session.store'
import { isAdmin } from '@/shared/lib/roles'
import { useTheme } from '@/shared/lib/use-theme'
import LogoutButton from '@/features/auth/logout/ui/LogoutButton.vue'

const session = useSessionStore()
const { isDark, toggle } = useTheme()
const isAdminUser = computed(() => isAdmin(session.user?.role))
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

@media (max-width: 768px) {
  .shell {
    flex-direction: column;
  }

  .sidebar {
    border-right: none;
    border-bottom: 1px solid var(--color-border);
    flex-direction: row;
    align-items: center;
    gap: 2px;
    overflow-x: auto;
    padding: 10px 16px;
    width: 100%;
  }

  .logo {
    margin: 0 14px 0 0;
    white-space: nowrap;
  }

  .sidebar a {
    padding: 7px 10px;
    white-space: nowrap;
  }

  .topbar {
    padding: 10px 16px;
  }

  .content {
    padding: 16px;
  }
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
  background: var(--color-brand-soft);
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

.theme-toggle {
  align-items: center;
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  cursor: pointer;
  display: inline-flex;
  height: 34px;
  justify-content: center;
  transition: border-color 0.15s ease, color 0.15s ease;
  width: 34px;
}

.theme-toggle:hover {
  border-color: var(--color-brand);
  color: var(--color-brand);
}
</style>
