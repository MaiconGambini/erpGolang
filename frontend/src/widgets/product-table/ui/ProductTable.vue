<template>
  <section class="panel">
    <div class="toolbar">
      <input v-model="searchInput" placeholder="Buscar por nome ou SKU" />
    </div>
    <p v-if="isLoading" class="state">Carregando...</p>
    <p v-else-if="isError" class="state error">Erro ao carregar produtos</p>
    <p v-else-if="!products.length" class="state">Nenhum produto encontrado</p>
    <table v-else>
      <thead>
        <tr>
          <th>Nome</th>
          <th>SKU</th>
          <th>Preço</th>
          <th>Estoque</th>
          <th>Status</th>
          <th class="actions-col">Ações</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="product in products" :key="product.id">
          <td>{{ product.name }}</td>
          <td>{{ product.sku }}</td>
          <td>R$ {{ formatPrice(product.price) }}</td>
          <td>
            <span :class="{ 'low-stock': product.stock <= 5 }">{{ product.stock }}</span>
          </td>
          <td>
            <span class="badge" :class="{ inactive: !product.active }">
              {{ product.active ? 'Ativo' : 'Inativo' }}
            </span>
          </td>
          <td class="actions-col">
            <button type="button" class="link" @click="emit('edit', product)">Editar</button>
            <button type="button" class="link danger" @click="emit('delete', product)">Excluir</button>
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
import type { Product } from '@/entities/product/model/types'
import { useListProducts } from '@/features/product/list/model/use-list-products'

const emit = defineEmits<{
  edit: [product: Product]
  delete: [product: Product]
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

const { data, isLoading, isError } = useListProducts({
  search: searchModel,
  active: ref(undefined),
  limit,
  offset,
})

const products = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.pagination.total ?? 0)

function formatPrice(value: string) {
  const n = Number(value)
  if (Number.isNaN(n)) return value
  return n.toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}
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

.low-stock {
  color: #dc2626;
  font-weight: 600;
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
