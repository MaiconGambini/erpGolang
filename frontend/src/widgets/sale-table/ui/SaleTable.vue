<template>
  <section class="panel">
    <div class="toolbar">
      <input v-model="searchInput" placeholder="Buscar por cliente" />
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
    <p v-if="isLoading" class="state">Carregando...</p>
    <p v-else-if="isError" class="state error">Erro ao carregar vendas</p>
    <p v-else-if="!sales.length" class="state">Nenhuma venda encontrada</p>
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
  border-radius: 10px;
  overflow: hidden;
}

.toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  padding: 12px;
}

.date-field {
  color: var(--color-text-secondary);
  display: flex;
  flex-direction: column;
  font-size: 12px;
  gap: 4px;
}

.date-field input {
  min-width: 140px;
}

input,
select {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 8px 12px;
}

input {
  flex: 1;
}

table {
  border-collapse: collapse;
  width: 100%;
}

th,
td {
  border-top: 1px solid var(--color-border);
  font-size: 13px;
  padding: 12px 16px;
  text-align: left;
}

.actions-col {
  white-space: nowrap;
  width: 280px;
}

.link {
  background: none;
  border: 0;
  color: var(--color-brand);
  cursor: pointer;
  font: inherit;
  font-size: 13px;
  margin-right: 8px;
  padding: 0;
}

.link.danger {
  color: #dc2626;
}

.link:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.badge {
  border-radius: 999px;
  font-size: 12px;
  padding: 3px 10px;
}

.badge.draft {
  background: #fef3c7;
  color: #92400e;
}

.badge.confirmed {
  background: #dcfce7;
  color: #15803d;
}

.badge.cancelled {
  background: #f3f4f6;
  color: #6b7280;
}

.state {
  color: var(--color-text-muted);
  padding: 24px;
}

.state.error {
  color: #dc2626;
}

.pagination {
  align-items: center;
  border-top: 1px solid var(--color-border);
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  padding: 12px 16px;
}

.pagination button {
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font: inherit;
  padding: 6px 12px;
}

.pagination button:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}
</style>
