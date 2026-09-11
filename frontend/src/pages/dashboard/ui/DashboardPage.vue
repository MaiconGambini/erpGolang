<template>
  <AppShell>
    <AppPageHeader title="Dashboard" description="Resumo da operação">
      <template #actions>
        <div v-if="canViewFinancialUser" class="header-actions">
          <label class="date-field">
            De
            <input v-model="fromDate" type="date" />
          </label>
          <label class="date-field">
            Até
            <input v-model="toDate" type="date" />
          </label>
          <button type="button" class="outline" :disabled="exportingPdf" @click="onExportPdf">
            Exportar PDF
          </button>
        </div>
      </template>
    </AppPageHeader>
    <p v-if="isError" class="error">Erro ao carregar o resumo. Tente novamente.</p>
    <section class="grid">
      <article
        v-for="metric in visibleMetrics"
        :key="metric.key"
        :class="{ clickable: metric.clickable }"
        :role="metric.clickable ? 'button' : undefined"
        :tabindex="metric.clickable ? 0 : undefined"
        @click="metric.clickable ? onMetricClick(metric) : undefined"
        @keydown.enter="metric.clickable ? onMetricClick(metric) : undefined"
        @keydown.space.prevent="metric.clickable ? onMetricClick(metric) : undefined"
      >
        {{ metric.label }}
        <strong v-if="isPending" class="skeleton" aria-hidden="true" />
        <strong v-else>{{ formatMetric(metric) }}</strong>
      </article>
    </section>

    <section v-if="canViewFinancialUser" class="charts">
      <article class="chart-card">
        <h2>Vendas por dia</h2>
        <p v-if="salesByDayLoading" class="chart-state">Carregando gráfico...</p>
        <p v-else-if="salesByDayError" class="chart-state error">Erro ao carregar vendas por dia</p>
        <div v-else class="chart-wrap">
          <canvas ref="salesChartRef" />
        </div>
      </article>
      <article class="chart-card">
        <h2>Produtos mais vendidos</h2>
        <p v-if="topProductsLoading" class="chart-state">Carregando gráfico...</p>
        <p v-else-if="topProductsError" class="chart-state error">Erro ao carregar produtos</p>
        <div v-else class="chart-wrap">
          <canvas ref="productsChartRef" />
        </div>
      </article>
    </section>
  </AppShell>
</template>

<script setup lang="ts">
import {
  BarController,
  BarElement,
  CategoryScale,
  Chart,
  Filler,
  Legend,
  LineController,
  LineElement,
  LinearScale,
  PointElement,
  Title,
  Tooltip,
} from 'chart.js'
import { useQuery } from '@tanstack/vue-query'
import { format, subDays } from 'date-fns'
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import type { DashboardSummary } from '@/entities/dashboard/model/types'
import { getSalesByDay, getTopProducts, downloadSalesSummaryPdf } from '@/entities/reports/api/reports.api'
import { useDashboardSummary } from '@/features/dashboard/summary/model/use-dashboard-summary'
import { canViewFinancial } from '@/shared/lib/roles'
import { useTheme } from '@/shared/lib/use-theme'
import AppPageHeader from '@/shared/ui/AppPageHeader.vue'
import { downloadBlob } from '@/shared/lib/download'
import { useSessionStore } from '@/entities/session/model/session.store'
import AppShell from '@/widgets/app-shell/ui/AppShell.vue'

Chart.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  LineController,
  BarController,
  Title,
  Tooltip,
  Legend,
  Filler,
)

interface DashboardMetric {
  label: string
  key: keyof DashboardSummary
  format: 'count' | 'currency'
  clickable?: boolean
  route?: string
}

const metrics: DashboardMetric[] = [
  { label: 'Clientes ativos', key: 'activeCustomers', format: 'count', clickable: true, route: '/customers' },
  { label: 'Novos clientes', key: 'newCustomers30d', format: 'count' },
  { label: 'Pedidos em aberto', key: 'draftSales', format: 'count', clickable: true, route: '/sales?status=draft' },
  { label: 'Alertas', key: 'lowStockAlerts', format: 'count', clickable: true, route: '/products?lowStock=1' },
  { label: 'Vendas confirmadas', key: 'confirmedSalesCount', format: 'count', clickable: true, route: '/sales?status=confirmed' },
  { label: 'Faturamento', key: 'confirmedSalesTotal', format: 'currency', clickable: true, route: '/sales?status=confirmed' },
]

const router = useRouter()
const session = useSessionStore()
const canViewFinancialUser = computed(() => canViewFinancial(session.user?.role))
const { data, isPending, isError } = useDashboardSummary()
const { isDark } = useTheme()

const visibleMetrics = computed(() =>
  metrics.filter(
    (metric) =>
      canViewFinancialUser.value
      || (metric.key !== 'confirmedSalesCount' && metric.key !== 'confirmedSalesTotal'),
  ),
)

const fromDate = ref(format(subDays(new Date(), 90), 'yyyy-MM-dd'))
const toDate = ref(format(new Date(), 'yyyy-MM-dd'))
const exportingPdf = ref(false)

const rangeParams = computed(() => ({
  from: fromDate.value,
  to: toDate.value,
}))

const {
  data: salesByDay,
  isLoading: salesByDayLoading,
  isError: salesByDayError,
} = useQuery({
  queryKey: computed(() => ['reports', 'sales-by-day', rangeParams.value]),
  queryFn: () => getSalesByDay(rangeParams.value),
  enabled: canViewFinancialUser,
})

const {
  data: topProducts,
  isLoading: topProductsLoading,
  isError: topProductsError,
} = useQuery({
  queryKey: computed(() => ['reports', 'top-products', rangeParams.value]),
  queryFn: () => getTopProducts({ ...rangeParams.value, limit: 5 }),
  enabled: canViewFinancialUser,
})

const salesChartRef = ref<HTMLCanvasElement | null>(null)
const productsChartRef = ref<HTMLCanvasElement | null>(null)
let salesChart: Chart | null = null
let productsChart: Chart | null = null

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

async function onExportPdf() {
  exportingPdf.value = true
  try {
    const blob = await downloadSalesSummaryPdf(fromDate.value, toDate.value)
    downloadBlob('vendas-resumo.pdf', blob)
  } finally {
    exportingPdf.value = false
  }
}

function chartColors() {
  const css = getComputedStyle(document.documentElement)
  const read = (name: string, fallback: string) => css.getPropertyValue(name).trim() || fallback
  const brand = read('--color-brand', '#2563eb')
  return {
    brand,
    brandFill: `${brand}1a`,
    grid: read('--color-border', '#e2e8f0'),
    tick: read('--color-text-muted', '#64748b'),
  }
}

function renderSalesChart() {
  const palette = chartColors()
  if (!salesChartRef.value || !salesByDay.value) return
  salesChart?.destroy()
  const labels = salesByDay.value.map((point) => {
    const [year, month, day] = point.date.split('-')
    return `${day}/${month}`
  })
  const totals = salesByDay.value.map((point) => Number(point.total) || 0)
  salesChart = new Chart(salesChartRef.value, {
    type: 'line',
    data: {
      labels,
      datasets: [{
        label: 'Faturamento (R$)',
        data: totals,
        borderColor: palette.brand,
        backgroundColor: palette.brandFill,
        fill: true,
        tension: 0.3,
      }],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: { legend: { display: false } },
      scales: {
        x: {
          grid: { color: palette.grid },
          ticks: { color: palette.tick },
        },
        y: {
          grid: { color: palette.grid },
          ticks: {
            color: palette.tick,
            callback: (value) => `R$ ${Number(value).toLocaleString('pt-BR')}`,
          },
        },
      },
    },
  })
}

function renderProductsChart() {
  if (!productsChartRef.value || !topProducts.value) return
  const palette = chartColors()
  productsChart?.destroy()
  const labels = topProducts.value.map((item) => item.productName)
  const quantities = topProducts.value.map((item) => item.quantity)
  productsChart = new Chart(productsChartRef.value, {
    type: 'bar',
    data: {
      labels,
      datasets: [{
        label: 'Quantidade vendida',
        data: quantities,
        backgroundColor: palette.brand,
        borderRadius: 4,
      }],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
    scales: {
      x: {
        grid: { color: palette.grid },
        ticks: { color: palette.tick },
      },
      y: {
        grid: { color: palette.grid },
        beginAtZero: true,
        ticks: { color: palette.tick },
      },
    },
    },
  })
}

watch(salesByDay, () => {
  if (!salesByDayLoading.value && !salesByDayError.value) {
    renderSalesChart()
  }
}, { flush: 'post' })

watch(topProducts, () => {
  if (!topProductsLoading.value && !topProductsError.value) {
    renderProductsChart()
  }
}, { flush: 'post' })

watch(isDark, () => {
  if (!salesByDayLoading.value && !salesByDayError.value) renderSalesChart()
  if (!topProductsLoading.value && !topProductsError.value) renderProductsChart()
}, { flush: 'post' })

onBeforeUnmount(() => {
  salesChart?.destroy()
  productsChart?.destroy()
})
</script>

<style scoped>
.header-actions {
  align-items: flex-end;
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.date-field {
  color: var(--color-text-secondary);
  display: flex;
  flex-direction: column;
  font-size: 12px;
  gap: 4px;
}

.date-field input {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font: inherit;
  min-width: 140px;
  padding: 8px 12px;
}

.outline {
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  cursor: pointer;
  font: inherit;
  font-weight: 600;
  min-height: 38px;
  padding: 0 16px;
}

.outline:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.grid {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  margin-bottom: 24px;
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
  font-variant-numeric: tabular-nums;
}

.grid .skeleton {
  height: 32px;
  width: 72px;
}

.error {
  color: var(--color-danger);
  margin: 0 0 16px;
}

.charts {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
}

.chart-card {
  min-height: 320px;
}

.chart-card h2 {
  color: var(--color-text-primary);
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 12px;
}

.chart-wrap {
  height: 260px;
  position: relative;
}

.chart-state {
  color: var(--color-text-muted);
  font-size: 14px;
  margin: 0;
}

.chart-state.error {
  color: var(--color-danger);
}
</style>
