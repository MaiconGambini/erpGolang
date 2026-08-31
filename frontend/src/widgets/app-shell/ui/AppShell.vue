<template>
  <div class="shell">
    <aside class="sidebar">
      <div class="brand-lockup" aria-label="goERP">
        <span class="brand-mark" aria-hidden="true">g</span>
        <span class="brand-name">go<span>ERP</span></span>
      </div>
      <nav class="nav" aria-label="Navegação principal">
        <RouterLink to="/">
          <LayoutDashboard :size="17" :stroke-width="2" aria-hidden="true" />
          <span>Dashboard</span>
        </RouterLink>
        <RouterLink to="/customers">
          <Users :size="17" :stroke-width="2" aria-hidden="true" />
          <span>Clientes</span>
        </RouterLink>
        <RouterLink to="/products">
          <Package :size="17" :stroke-width="2" aria-hidden="true" />
          <span>Produtos</span>
        </RouterLink>
        <RouterLink to="/suppliers">
          <Truck :size="17" :stroke-width="2" aria-hidden="true" />
          <span>Fornecedores</span>
        </RouterLink>
        <RouterLink to="/sales">
          <ShoppingCart :size="17" :stroke-width="2" aria-hidden="true" />
          <span>Vendas</span>
        </RouterLink>
        <template v-if="isAdminUser">
          <RouterLink to="/users">
            <UserCog :size="17" :stroke-width="2" aria-hidden="true" />
            <span>Usuários</span>
          </RouterLink>
          <RouterLink to="/audit">
            <ScrollText :size="17" :stroke-width="2" aria-hidden="true" />
            <span>Auditoria</span>
          </RouterLink>
        </template>
      </nav>
    </aside>
    <div class="main">
      <header class="topbar">
        <span class="topbar-spacer" aria-hidden="true" />
        <div class="topbar-user">
          <span class="avatar" aria-hidden="true">{{ userInitials }}</span>
          <div class="user-copy">
            <span class="user-label">Sessão ativa</span>
            <span class="user">{{ userName }}</span>
          </div>
        </div>
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
import { LayoutDashboard, Moon, Package, ScrollText, ShoppingCart, Sun, Truck, UserCog, Users } from 'lucide-vue-next'
import { useSessionStore } from '@/entities/session/model/session.store'
import { isAdmin } from '@/shared/lib/roles'
import { useTheme } from '@/shared/lib/use-theme'
import LogoutButton from '@/features/auth/logout/ui/LogoutButton.vue'

const session = useSessionStore()
const { isDark, toggle } = useTheme()
const isAdminUser = computed(() => isAdmin(session.user?.role))
const userName = computed(() => session.user?.name ?? 'Usuário')
const userInitials = computed(() =>
  userName.value
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0])
    .join('')
    .toUpperCase(),
)
</script>



<style scoped>
.shell {
  background: var(--color-canvas);
  display: flex;
  min-height: 100vh;
}

.sidebar {
  background: var(--color-surface);
  border-right: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  gap: 24px;
  padding: 22px 14px;
  width: 248px;
}

.brand-lockup {
  align-items: center;
  display: flex;
  gap: 10px;
  padding: 0 10px;
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

.brand-name {
  font-size: 20px;
  font-weight: 750;
  letter-spacing: -0.04em;
}

.brand-name span {
  color: var(--color-brand);
}

.nav {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 3px;
}

.nav a {
  align-items: center;
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  display: flex;
  gap: 11px;
  min-height: 42px;
  padding: 0 12px;
  text-decoration: none;
  transition: background-color 0.15s ease, color 0.15s ease, transform 0.15s ease;
}

.nav a svg {
  color: var(--color-text-muted);
  flex-shrink: 0;
  transition: color 0.15s ease;
}

.nav a:hover {
  background: var(--color-surface-muted);
  color: var(--color-text-primary);
  transform: translateX(2px);
}

.nav a.router-link-active {
  background: var(--color-brand-soft);
  color: var(--color-brand-hover);
  font-weight: 650;
}

.nav a.router-link-active svg {
  color: var(--color-brand);
}

.main {
  display: flex;
  flex: 1;
  flex-direction: column;
  min-width: 0;
}

.topbar {
  align-items: center;
  background: color-mix(in srgb, var(--color-surface) 94%, transparent);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  gap: 12px;
  min-height: 72px;
  padding: 12px 32px;
}

.topbar-spacer {
  flex: 1;
}

.topbar-user {
  align-items: center;
  display: flex;
  gap: 10px;
  margin-right: 4px;
}

.avatar {
  align-items: center;
  background: var(--color-brand-soft);
  border: 1px solid color-mix(in srgb, var(--color-brand) 22%, var(--color-border));
  border-radius: 50%;
  color: var(--color-brand-hover);
  display: inline-flex;
  font-size: 12px;
  font-weight: 750;
  height: 34px;
  justify-content: center;
  width: 34px;
}

.user-copy {
  display: flex;
  flex-direction: column;
  line-height: 1.15;
}

.user-label {
  color: var(--color-text-muted);
  font-size: 10px;
  font-weight: 650;
  letter-spacing: 0.07em;
  text-transform: uppercase;
}

.user {
  color: var(--color-text-primary);
  font-size: 13px;
  font-weight: 600;
}

.content {
  flex: 1;
  padding: 32px;
}

.theme-toggle {
  align-items: center;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  cursor: pointer;
  display: inline-flex;
  height: 38px;
  justify-content: center;
  transition: background-color 0.15s ease, border-color 0.15s ease, color 0.15s ease;
  width: 38px;
}

.theme-toggle:hover {
  background: var(--color-brand-soft);
  border-color: var(--color-brand);
  color: var(--color-brand);
}

@media (max-width: 768px) {
  .shell {
    flex-direction: column;
  }

  .sidebar {
    border-bottom: 1px solid var(--color-border);
    border-right: none;
    flex-direction: row;
    gap: 12px;
    overflow-x: auto;
    padding: 10px 16px;
    width: 100%;
  }

  .brand-lockup {
    flex-shrink: 0;
    padding: 0;
  }

  .nav {
    flex-direction: row;
    gap: 3px;
    min-width: max-content;
  }

  .nav a {
    min-height: 38px;
    padding: 0 10px;
    white-space: nowrap;
  }

  .topbar {
    min-height: 64px;
    padding: 10px 16px;
  }

  .content {
    padding: 24px 16px;
  }
}

@media (max-width: 480px) {
  .brand-name {
    display: none;
  }

  .topbar-user {
    margin-right: auto;
  }

  .user-label {
    display: none;
  }
}
</style>
