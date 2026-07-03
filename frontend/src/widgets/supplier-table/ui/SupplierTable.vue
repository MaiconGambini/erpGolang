<template>
  <section class="panel">
    <div class="toolbar">
      <input v-model="searchInput" placeholder="Buscar por nome ou documento" />
    </div>
    <p v-if="isLoading" class="state">Carregando...</p>
    <p v-else-if="isError" class="state error">Erro ao carregar fornecedores</p>
    <p v-else-if="!suppliers.length" class="state">Nenhum fornecedor encontrado</p>
    <table v-else>
      <thead>
        <tr>
          <th>Nome</th>
          <th>Documento</th>
          <th>E-mail</th>
          <th>Status</th>
          <th class="actions-col">Ações</th>
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
          <td class="actions-col">
            <button type="button" class="link" @click="emit('edit', supplier)">Editar</button>
            <button type="button" class="link danger" @click="emit('delete', supplier)">Excluir</button>
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
import { useListSuppliers } from '@/features/supplier/list/model/use-list-suppliers'

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
  border-radius: 10px;
  overflow: hidden;
}

.toolbar {
  display: flex;
  gap: 10px;
  padding: 12px;
}

input {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  flex: 1;
  padding: 8px 12px;
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
  width: 140px;
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

.badge {
  background: #dcfce7;
  border-radius: 999px;
  color: #15803d;
  padding: 3px 10px;
}

.badge.inactive {
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
