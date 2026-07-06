<template>
  <AppShell>
    <AppPageHeader title="Dashboard" description="Resumo da operação" />
    <p v-if="isError" class="error">Erro ao carregar o resumo. Tente novamente.</p>
    <section class="grid">
      <article
        v-for="metric in metrics"
        :key="metric.key"
        :class="{ clickable: metric.clickable }"
        :role="metric.clickable ? 'button' : undefined"
        :tabindex="metric.clickable ? 0 : undefined"
        @click="metric.clickable ? onMetricClick(metric) : undefined"
        @keydown.enter="metric.clickable ? onMetricClick(metric) : undefined"
      >
        {{ metric.label }}
        <strong v-if="isPending" class="skeleton" aria-hidden="true" />
        <strong v-else>{{ formatMetric(metric) }}</strong>
      </article>
    </section>
  </AppShell>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import AppShell from '@/widgets/app-shell/ui/AppShell.vue'
import AppPageHeader from '@/shared/ui/AppPageHeader.vue'
import { useDashboardSummary } from '@/features/dashboard/summary/model/use-dashboard-summary'
import type { DashboardSummary } from '@/entities/dashboard/model/types'

interface DashboardMetric {
  label: string
  key: keyof DashboardSummary
  format: 'count' | 'currency'
  clickable?: boolean
  route?: string
}

const metrics: DashboardMetric[] = [
  { label: 'Clientes ativos', key: 'activeCustomers', format: 'count' },
  { label: 'Novos clientes', key: 'newCustomers30d', format: 'count' },
  { label: 'Pedidos em aberto', key: 'draftSales', format: 'count' },
  { label: 'Alertas', key: 'lowStockAlerts', format: 'count', clickable: true, route: '/products?lowStock=1' },
  { label: 'Vendas confirmadas', key: 'confirmedSalesCount', format: 'count' },
  { label: 'Faturamento', key: 'confirmedSalesTotal', format: 'currency' },
]

const router = useRouter()
const { data, isPending, isError } = useDashboardSummary()

function formatCount(value: number | undefined) {
  return new Intl.NumberFormat('pt-BR').format(value ?? 0)
}

function formatCurrency(value: string | undefined) {
  const n = Number(value ?? 0)
  if (Number.isNaN(n)) return value ?? 'R$ 0,00'
  return new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(n)
}

function formatMetric(metric: DashboardMetric) {
  const value = data.value?.[metric.key]
  if (metric.format === 'currency') {
    return formatCurrency(typeof value === 'string' ? value : String(value ?? 0))
  }
  return formatCount(typeof value === 'number' ? value : Number(value ?? 0))
}

function onMetricClick(metric: DashboardMetric) {
  if (metric.route) {
    router.push(metric.route)
  }
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

article.clickable {
  cursor: pointer;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

article.clickable:hover,
article.clickable:focus-visible {
  border-color: var(--color-brand);
  box-shadow: 0 0 0 1px var(--color-brand);
  outline: none;
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
