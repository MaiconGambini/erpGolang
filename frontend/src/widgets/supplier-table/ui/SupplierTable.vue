<template>
  <section class="panel">
    <div class="toolbar">
      <input v-model="searchInput" aria-label="Buscar fornecedores" placeholder="Buscar por nome ou documento" />
    </div>
    <TableSkeleton v-if="isLoading" :cols="5" />
    <p v-else-if="isError" class="state error">Erro ao carregar fornecedores</p>
    <EmptyState
      v-else-if="!suppliers.length"
      :icon="Truck"
      title="Nenhum fornecedor encontrado"
      hint="Ajuste a busca ou cadastre o primeiro fornecedor usando “Novo fornecedor”."
    />
    <table v-else>
      <thead>
        <tr>
          <th>Nome</th>
          <th>Documento</th>
          <th>E-mail</th>
          <th>Status</th>
          <th v-if="showActions" class="actions-col">Ações</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="supplier in suppliers" :key="supplier.id">
          <td>{{ supplier.name }}</td>
          <td>{{ supplier.document ?? '—' }}</td>
          <td>{{ supplier.email ?? '—' }}</td>
          <td>
            <span class="badge" :class="{ inactive: !supplier.active }">
              {{ supplier.active ? 'Ativo' : 'Inativo' }}
            </span>
          </td>
          <td v-if="showActions" class="actions-col">
            <button v-if="canWriteUser" type="button" class="link" @click="emit('edit', supplier)">Editar</button>
            <button
              v-if="canDeleteUser"
              type="button"
              class="link danger"
              @click="emit('delete', supplier)"
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
import type { Supplier } from '@/entities/supplier/model/types'
import { useSessionStore } from '@/entities/session/model/session.store'
import { useListSuppliers } from '@/features/supplier/list/model/use-list-suppliers'
import { canDelete, canWrite } from '@/shared/lib/roles'
import EmptyState from '@/shared/ui/EmptyState.vue'
import TableSkeleton from '@/shared/ui/TableSkeleton.vue'
import { Truck } from 'lucide-vue-next'

const session = useSessionStore()
const canWriteUser = computed(() => canWrite(session.user?.role))
const canDeleteUser = computed(() => canDelete(session.user?.role))
const showActions = computed(() => canWriteUser.value || canDeleteUser.value)

const emit = defineEmits<{
  edit: [supplier: Supplier]
  delete: [supplier: Supplier]
}>()

const searchInput = ref('')
const searchModel = ref('')
const limit = ref(20)
const offset = ref(0)

let searchDebounce: ReturnType<typeof setTimeout> | undefined

watch(searchInput, (value) => {
  if (searchDebounce) clearTimeout(searchDebounce)
  searchDebounce = setTimeout(() => {
    searchModel.value = value
    offset.value = 0
  }, 300)
})

const { data, isLoading, isError } = useListSuppliers({
  search: searchModel,
  active: ref(undefined),
  limit,
  offset,
})

const suppliers = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.pagination.total ?? 0)
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
  gap: 10px;
  padding: 16px;
}

input {
  background: var(--color-surface-raised);
  border: 1px solid var(--color-border-strong);
  border-radius: var(--radius-md);
  color: var(--color-text-primary);
  flex: 1;
  min-height: 40px;
  outline: none;
  padding: 0 12px;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

input:focus {
  border-color: var(--color-brand);
  box-shadow: var(--focus-ring);
}

table {
  border-collapse: collapse;
  min-width: 620px;
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
  width: 140px;
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

.badge {
  background: var(--color-success-soft);
  border-radius: 999px;
  color: var(--color-success-strong);
  font-size: 12px;
  font-weight: 650;
  line-height: 1;
  padding: 6px 9px;
  white-space: nowrap;
}

.badge.inactive {
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
