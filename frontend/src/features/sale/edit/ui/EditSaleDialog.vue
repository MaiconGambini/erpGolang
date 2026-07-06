<template>
  <div v-if="visible" class="overlay" @click.self="emit('close')">
    <div class="dialog" role="dialog" aria-labelledby="edit-title">
      <h2 id="edit-title">Editar venda</h2>
      <p v-if="isLoading" class="state">Carregando...</p>
      <p v-else-if="isError" class="state error">Erro ao carregar venda</p>
      <p v-else-if="sale && sale.status !== 'draft'" class="state error">
        Apenas vendas em rascunho podem ser editadas
      </p>
      <form v-else-if="sale" class="form" @submit.prevent="onSubmit">
        <label>
          Cliente *
          <select v-model="form.customerId" required>
            <option value="">Selecione...</option>
            <option v-for="c in customers" :key="c.id" :value="c.id">{{ c.name }}</option>
          </select>
          <span v-if="errors.customerId" class="field-error">{{ errors.customerId }}</span>
        </label>

        <fieldset class="items">
          <legend>Itens *</legend>
          <div v-for="(row, index) in form.items" :key="index" class="item-row">
            <select v-model="row.productId" required>
              <option value="">Produto...</option>
              <option v-for="p in products" :key="p.id" :value="p.id">
                {{ p.name }} ({{ p.sku }}) — estoque {{ p.stock }}
              </option>
            </select>
            <input v-model.number="row.quantity" type="number" min="1" step="1" placeholder="Qtd" />
            <button v-if="form.items.length > 1" type="button" class="link danger" @click="removeRow(index)">
              Remover
            </button>
          </div>
          <button type="button" class="link" @click="addRow">+ Adicionar item</button>
          <span v-if="errors.items" class="field-error">{{ errors.items }}</span>
        </fieldset>

        <label>
          Observações
          <textarea v-model="form.notes" rows="2" />
        </label>

        <p v-if="submitError" class="submit-error">{{ submitError }}</p>
        <div class="actions">
          <button type="button" class="secondary" @click="emit('close')">Cancelar</button>
          <AppButton type="submit" :disabled="isPending">Salvar rascunho</AppButton>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { listCustomers } from '@/entities/customer/api/customer.api'
import { listProducts } from '@/entities/product/api/product.api'
import { getSale } from '@/entities/sale/api/sale.api'
import { createSaleSchema } from '@/entities/sale/model/schemas'
import { useEditSale } from '../model/use-edit-sale'
import { getApiErrorMessage } from '@/shared/api/errors'
import AppButton from '@/shared/ui/AppButton.vue'

const props = defineProps<{ visible: boolean; saleId: string | null }>()
const emit = defineEmits<{ close: [] }>()

const form = reactive({
  customerId: '',
  notes: '',
  items: [{ productId: '', quantity: 1 }],
})
const errors = ref<Record<string, string>>({})
const submitError = ref('')

const saleIdRef = computed(() => props.saleId)

const { data: sale, isLoading, isError } = useQuery({
  queryKey: computed(() => ['sale', props.saleId]),
  queryFn: () => getSale(props.saleId!),
  enabled: computed(() => props.visible && !!props.saleId),
})

const { data: customersData } = useQuery({
  queryKey: ['customers', 'select'],
  queryFn: () => listCustomers({ limit: 100, active: true }),
  enabled: computed(() => props.visible),
})
const { data: productsData } = useQuery({
  queryKey: ['products', 'select'],
  queryFn: () => listProducts({ limit: 100, active: true }),
  enabled: computed(() => props.visible),
})

const customers = computed(() => customersData.value?.data ?? [])
const products = computed(() => productsData.value?.data ?? [])

const { mutate, isPending } = useEditSale()

watch(
  () => [props.visible, sale.value] as const,
  ([open, currentSale]) => {
    if (open && currentSale?.status === 'draft') {
      form.customerId = currentSale.customerId
      form.notes = currentSale.notes ?? ''
      form.items = (currentSale.items ?? []).map((item) => ({
        productId: item.productId,
        quantity: item.quantity,
      }))
      if (!form.items.length) {
        form.items = [{ productId: '', quantity: 1 }]
      }
      errors.value = {}
      submitError.value = ''
    }
  },
)

function addRow() {
  form.items.push({ productId: '', quantity: 1 })
}

function removeRow(index: number) {
  form.items.splice(index, 1)
}

function onSubmit() {
  if (!saleIdRef.value) return
  errors.value = {}
  submitError.value = ''
  const parsed = createSaleSchema.safeParse({
    customerId: form.customerId,
    notes: form.notes || undefined,
    items: form.items,
  })
  if (!parsed.success) {
    for (const issue of parsed.error.issues) {
      const key = issue.path.length ? String(issue.path[0]) : '_'
      errors.value[key] = issue.message
    }
    return
  }
  mutate(
    {
      id: saleIdRef.value,
      data: {
        customerId: parsed.data.customerId,
        notes: parsed.data.notes,
        items: parsed.data.items,
      },
    },
    {
      onSuccess: () => emit('close'),
      onError: (error) => {
        submitError.value = getApiErrorMessage(error, 'Não foi possível atualizar a venda')
      },
    },
  )
}
</script>

<style scoped>
.overlay {
  align-items: center;
  background: rgba(15, 23, 42, 0.45);
  display: flex;
  inset: 0;
  justify-content: center;
  position: fixed;
  z-index: 50;
}

.dialog {
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-soft);
  max-width: 560px;
  max-height: 90vh;
  overflow-y: auto;
  padding: 24px;
  width: 100%;
}

h2 {
  font-size: 18px;
  margin: 0 0 20px;
}

.state {
  color: var(--color-text-muted);
  font-size: 14px;
}

.state.error {
  color: #dc2626;
}

.form {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

label,
fieldset {
  color: var(--color-text-secondary);
  display: flex;
  flex-direction: column;
  font-size: 13px;
  gap: 6px;
}

fieldset.items {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 12px;
}

.item-row {
  align-items: center;
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

.item-row select {
  flex: 2;
}

.item-row input {
  width: 80px;
}

select,
textarea,
input {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font: inherit;
  padding: 8px 12px;
}

.link {
  background: none;
  border: 0;
  color: var(--color-brand);
  cursor: pointer;
  font: inherit;
  font-size: 13px;
  padding: 0;
}

.link.danger {
  color: #dc2626;
}

.field-error,
.submit-error {
  color: #dc2626;
  font-size: 12px;
}

.actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 8px;
}

.secondary {
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  cursor: pointer;
  font: inherit;
  padding: 8px 14px;
}
</style>
