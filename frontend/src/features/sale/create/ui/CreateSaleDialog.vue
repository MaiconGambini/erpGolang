<template>
  <AppDialog :visible="visible" title-id="create-title" @close="emit('close')">
      <h2 id="create-title">Nova venda</h2>
      <form class="form" @submit.prevent="onSubmit">
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
          <div v-for="(row, index) in form.items" :key="row.key" class="item-row">
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
  </AppDialog>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { listCustomers } from '@/entities/customer/api/customer.api'
import { listProducts } from '@/entities/product/api/product.api'
import { createSaleSchema } from '@/entities/sale/model/schemas'
import { useCreateSale } from '../model/use-create-sale'
import { getApiErrorMessage } from '@/shared/api/errors'
import AppDialog from '@/shared/ui/AppDialog.vue'
import AppButton from '@/shared/ui/AppButton.vue'
const props = defineProps<{ visible: boolean }>()
const emit = defineEmits<{ close: [] }>()

let itemKey = 0
const newItem = () => ({ key: `sale-item-${itemKey++}`, productId: '', quantity: 1 })

const form = reactive({
  customerId: '',
  notes: '',
  items: [newItem()],
})
const errors = ref<Record<string, string>>({})
const submitError = ref('')

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

const { mutate, isPending } = useCreateSale()

watch(() => props.visible, (open) => {
  if (open) {
    form.customerId = ''
    form.notes = ''
    form.items = [newItem()]
    errors.value = {}
    submitError.value = ''
  }
})

function addRow() {
  form.items.push(newItem())
}

function removeRow(index: number) {
  form.items.splice(index, 1)
}

function onSubmit() {
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
      customerId: parsed.data.customerId,
      notes: parsed.data.notes,
      items: parsed.data.items,
    },
    {
      onSuccess: () => emit('close'),
      onError: (error) => {
        submitError.value = getApiErrorMessage(error, 'Não foi possível criar a venda')
      },
    },
  )
}
</script>

<style scoped>

h2 {
  font-size: 18px;
  font-weight: 650;
  margin: 0 0 20px;
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
  color: var(--color-danger);
}

.field-error,
.submit-error {
  color: var(--color-danger);
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
