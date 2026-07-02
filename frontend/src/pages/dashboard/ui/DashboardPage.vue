<template>
  <AppShell>
    <AppPageHeader title="Dashboard" description="Resumo da operação" />
    <section class="grid">
      <article>Clientes ativos<strong>1.284</strong></article>
      <article>Novos clientes<strong>37</strong></article>
      <article>Pedidos em aberto<strong>18</strong></article>
      <article>Alertas<strong>3</strong></article>
    </section>
  </AppShell>
</template>

<script setup lang="ts">
import AppShell from '@/widgets/app-shell/ui/AppShell.vue'
import AppPageHeader from '@/shared/ui/AppPageHeader.vue'
</script>

<style scoped>
.grid {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
}

article {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 10px;
  color: var(--color-text-muted);
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 18px 20px;
}

strong {
  color: var(--color-text-primary);
  font-size: 26px;
}
</style>
