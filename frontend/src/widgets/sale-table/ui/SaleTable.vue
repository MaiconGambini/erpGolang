<template>
  <section class="panel">
    <div class="toolbar">
      <input v-model="searchInput" aria-label="Buscar vendas por cliente" placeholder="Buscar por cliente" />
      <select v-model="statusFilter">
        <option value="">Todos os status</option>
        <option value="draft">Rascunho</option>
        <option value="confirmed">Confirmada</option>
        <option value="cancelled">Cancelada</option>
      </select>
      <label class="date-field">
        De
        <input v-model="fromDate" type="date" />
      </label>
      <label class="date-field">
        Até
        <input v-model="toDate" type="date" />
      </label>
    </div>
    <p v-if="actionError" class="state error">{{ actionError }}</p>
    <TableSkeleton v-if="isLoading" :cols="5" />
    <p v-else-if="isError" class="state error">Erro ao carregar vendas</p>
    <EmptyState
      v-else-if="!sales.length"
      :icon="ReceiptText"
      title="Nenhuma venda encontrada"
      hint="Ajuste os filtros de período, status ou busca por cliente."
    />
    <table v-else>
      <thead>
        <tr>
          <th>Cliente</th>
          <th>Total</th>
          <th>Status</th>
          <th>Data</th>
          <th class="actions-col">Ações</th>        </tr>
      </thead>
      <tbody>
        <tr v-for="sale in sales" :key="sale.id">
          <td>{{ sale.customerName }}</td>
          <td>R$ {{ formatPrice(sale.total) }}</td>
          <td>
            <span class="badge" :class="sale.status">{{ statusLabel(sale.status) }}</span>
          </td>
          <td>{{ formatDate(sale.createdAt) }}</td>
          <td class="actions-col">
            <button type="button" class="link" @click="emit('view', sale.id)">Ver</button>
            <button
              v-if="sale.status === 'draft' && canWriteUser"
              type="button"
              class="link"
              @click="emit('edit', sale.id)"
            >
              Editar
            </button>
            <button
              v-if="sale.status === 'draft' && canWriteUser"
              type="button"
              class="link"
              :disabled="pendingId === sale.id"
              @click="onConfirm(sale.id)"
            >
              Confirmar
            </button>
            <button
              v-if="sale.status === 'confirmed' && canManageSalesUser"
              type="button"
              class="link danger"
              :disabled="pendingId === sale.id"
              @click="onCancel(sale.id)"
            >
              Cancelar
            </button>
            <button
              v-if="sale.status === 'draft' && canManageSalesUser"
              type="button"
              class="link danger"
              :disabled="pendingId === sale.id"
              @click="onDelete(sale.id)"
            >
              Excluir
            </button>
          </td>
        </tr>
      </tbody>
    </table>
    <footer v-if="total > limit" class="pagination">
      <button type="button" :disabled="offset === 0" @click="offset = Math.max(0, offset - limit)">
        Anterior
      </button>
      <span>{{ offset + 1 }}–{{ Math.min(offset + limit, total) }} de {{ total }}</span>
      <button type="button" :disabled="offset + limit >= total" @click="offset += limit">
        Próximo
      </button>
    </footer>
  </section>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import type { SaleStatus } from '@/entities/sale/model/types'
import { useSessionStore } from '@/entities/session/model/session.store'
import { useCancelSale, useConfirmSale, useDeleteSale } from '@/features/sale/actions/model/use-sale-actions'
import { useListSales } from '@/features/sale/list/model/use-list-sales'
import { getApiErrorMessage } from '@/shared/api/errors'
import { canManageSales, canWrite } from '@/shared/lib/roles'
import EmptyState from '@/shared/ui/EmptyState.vue'
import TableSkeleton from '@/shared/ui/TableSkeleton.vue'
import { ReceiptText } from 'lucide-vue-next'

const route = useRoute()
const session = useSessionStore()
const canWriteUser = computed(() => canWrite(session.user?.role))
const canManageSalesUser = computed(() => canManageSales(session.user?.role))

const emit = defineEmits<{
  view: [id: string]
  edit: [id: string]
}>()

const searchInput = ref('')
const searchModel = ref('')
const statusFilter = ref<SaleStatus | ''>((route.query.status as SaleStatus) || '')
const fromDate = ref('')
const toDate = ref('')
const limit = ref(20)
const offset = ref(0)
const pendingId = ref('')
const actionError = ref('')

let searchDebounce: ReturnType<typeof setTimeout> | undefined

watch(
  () => route.query.status,
  (value) => {
    statusFilter.value = (typeof value === 'string' ? value : '') as SaleStatus | ''
    offset.value = 0
  },
)

watch(searchInput, (value) => {
  if (searchDebounce) clearTimeout(searchDebounce)
  searchDebounce = setTimeout(() => {
    searchModel.value = value
    offset.value = 0
  }, 300)
})

watch(statusFilter, () => { offset.value = 0 })
watch([fromDate, toDate], () => { offset.value = 0 })

const statusRef = computed(() => (statusFilter.value || undefined) as SaleStatus | undefined)
const fromRef = computed(() => fromDate.value || undefined)
const toRef = computed(() => toDate.value || undefined)

const { data, isLoading, isError } = useListSales({
  search: searchModel,
  status: statusRef,
  from: fromRef,
  to: toRef,
  limit,
  offset,
})

const { mutate: confirmMutate } = useConfirmSale()
const { mutate: cancelMutate } = useCancelSale()
const { mutate: deleteMutate } = useDeleteSale()

const sales = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.pagination.total ?? 0)

function formatPrice(value: string) {
  const n = Number(value)
  if (Number.isNaN(n)) return value
  return n.toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

function formatDate(value: string) {
  return new Date(value).toLocaleDateString('pt-BR')
}

function statusLabel(status: SaleStatus) {
  const map: Record<SaleStatus, string> = {
    draft: 'Rascunho',
    confirmed: 'Confirmada',
    cancelled: 'Cancelada',
  }
  return map[status]
}

function onConfirm(id: string) {
  actionError.value = ''
  pendingId.value = id
  confirmMutate(id, {
    onSettled: () => { pendingId.value = '' },
    onError: (error) => {
      actionError.value = getApiErrorMessage(error, 'Não foi possível confirmar a venda')
    },
  })
}

function onCancel(id: string) {
  if (!confirm('Deseja cancelar esta venda confirmada? O estoque será revertido.')) return
  actionError.value = ''
  pendingId.value = id
  cancelMutate(id, {
    onSettled: () => { pendingId.value = '' },
    onError: (error) => {
      actionError.value = getApiErrorMessage(error, 'Não foi possível cancelar a venda')
    },
  })
}

function onDelete(id: string) {
  if (!confirm('Deseja excluir este rascunho de venda?')) return
  actionError.value = ''
  pendingId.value = id
  deleteMutate(id, {
    onSettled: () => { pendingId.value = '' },
    onError: (error) => {
      actionError.value = getApiErrorMessage(error, 'Não foi possível excluir a venda')
    },
  })
}

const exportParams = computed(() => ({
  search: searchModel.value || undefined,
  status: statusFilter.value || undefined,
  from: fromDate.value || undefined,
  to: toDate.value || undefined,
}))

defineExpose({ exportParams })
</script>

<style scoped>
.panel {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-card);
  overflow-x: auto;
}

.toolbar {
  background: var(--color-surface-muted);
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  padding: 16px;
}

.date-field {
  color: var(--color-text-secondary);
  display: flex;
  flex-direction: column;
  font-size: 12px;
  font-weight: 600;
  gap: 4px;
}

.date-field input {
  min-width: 140px;
}

input,
select {
  background: var(--color-surface-raised);
  border: 1px solid var(--color-border-strong);
  border-radius: var(--radius-md);
  color: var(--color-text-primary);
  font: inherit;
  min-height: 40px;
  outline: none;
  padding: 0 12px;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

input {
  flex: 1;
  min-width: 180px;
}

select {
  min-width: 150px;
}

input:focus,
select:focus {
  border-color: var(--color-brand);
  box-shadow: var(--focus-ring);
}

table {
  border-collapse: collapse;
  min-width: 860px;
  width: 100%;
}

th,
td {
  border-top: 1px solid var(--color-border);
  font-size: 13px;
  padding: 13px 16px;
  text-align: left;
}

th {
  background: var(--color-surface-muted);
  color: var(--color-text-muted);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  white-space: nowrap;
}

tbody tr {
  transition: background-color 0.15s ease;
}

tbody tr:hover {
  background: var(--color-surface-muted);
}

.actions-col {
  white-space: nowrap;
  width: 280px;
}

.link {
  background: none;
  border: 0;
  border-radius: var(--radius-sm);
  color: var(--color-brand);
  cursor: pointer;
  font: inherit;
  font-size: 13px;
  margin-right: 4px;
  min-height: 32px;
  padding: 0 5px;
  transition: background-color 0.15s ease, color 0.15s ease;
}

.link:hover {
  background: var(--color-brand-soft);
  color: var(--color-brand-hover);
  text-decoration: underline;
  text-underline-offset: 3px;
}

.link:focus-visible {
  outline: none;
  box-shadow: var(--focus-ring);
}

.link.danger {
  color: var(--color-danger);
}

.link.danger:hover {
  background: var(--color-danger-soft);
  color: var(--color-danger-strong);
}

.link:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.badge {
  border-radius: 999px;
  font-size: 12px;
  font-weight: 650;
  line-height: 1;
  padding: 6px 9px;
  white-space: nowrap;
}

.badge.draft {
  background: var(--color-warning-soft);
  color: var(--color-warning-strong);
}

.badge.confirmed {
  background: var(--color-success-soft);
  color: var(--color-success-strong);
}

.badge.cancelled {
  background: var(--color-surface-muted);
  color: var(--color-text-muted);
}

.state {
  color: var(--color-text-muted);
  padding: 24px;
}

.state.error {
  color: var(--color-danger);
}

.pagination {
  align-items: center;
  background: var(--color-surface-muted);
  border-top: 1px solid var(--color-border);
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  justify-content: flex-end;
  padding: 12px 16px;
}

.pagination button {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-secondary);
  cursor: pointer;
  font: inherit;
  min-height: 34px;
  padding: 0 12px;
  transition: background-color 0.15s ease, border-color 0.15s ease;
}

.pagination button:hover:not(:disabled) {
  background: var(--color-surface-raised);
  border-color: var(--color-border-strong);
}

.pagination button:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}
</style>
