<template>
  <div class="shell">
    <aside class="sidebar">
      <div class="logo">go<span>ERP</span></div>
      <RouterLink to="/">Dashboard</RouterLink>
      <RouterLink to="/customers">Clientes</RouterLink>
    </aside>
    <main class="content">
      <slot />
    </main>
  </div>
</template>

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

.content {
  flex: 1;
  padding: 28px 32px;
}
</style>
