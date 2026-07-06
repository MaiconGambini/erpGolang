<template>
  <AppShell>
    <AppPageHeader title="Dashboard" description="Resumo da operação" />
    <p v-if="isError" class="error">Erro ao carregar o resumo. Tente novamente.</p>
    <section class="grid">
      <article v-for="metric in metrics" :key="metric.key">
        {{ metric.label }}
        <strong v-if="isPending" class="skeleton" aria-hidden="true" />
        <strong v-else>{{ formatCount(data?.[metric.key]) }}</strong>
      </article>
    </section>
  </AppShell>
</template>

<script setup lang="ts">
import AppShell from '@/widgets/app-shell/ui/AppShell.vue'
import AppPageHeader from '@/shared/ui/AppPageHeader.vue'
import { useDashboardSummary } from '@/features/dashboard/summary/model/use-dashboard-summary'
import type { DashboardSummary } from '@/entities/dashboard/model/types'

const metrics: { label: string; key: keyof DashboardSummary }[] = [
  { label: 'Clientes ativos', key: 'active_customers' },
  { label: 'Novos clientes', key: 'new_customers_30d' },
  { label: 'Pedidos em aberto', key: 'draft_sales' },
  { label: 'Alertas', key: 'low_stock_alerts' },
]

const { data, isPending, isError } = useDashboardSummary()

function formatCount(value: number | undefined) {
  return new Intl.NumberFormat('pt-BR').format(value ?? 0)
}
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

.skeleton {
  background: var(--color-border);
  border-radius: 6px;
  display: block;
  height: 32px;
  width: 72px;
}

.error {
  color: var(--color-danger, #b91c1c);
  margin: 0 0 16px;
}
</style>
