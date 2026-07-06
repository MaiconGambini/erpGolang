<template>
  <div v-if="visible && product" class="overlay" @click.self="emit('close')">
    <div class="dialog" role="dialog" aria-labelledby="edit-title">
      <h2 id="edit-title">Editar produto</h2>
      <form class="form" @submit.prevent="onSubmit">
        <label>
          Nome *
          <input v-model="form.name" required />
          <span v-if="errors.name" class="field-error">{{ errors.name }}</span>
        </label>
        <label>
          SKU *
          <input v-model="form.sku" required />
          <span v-if="errors.sku" class="field-error">{{ errors.sku }}</span>
        </label>
        <label>
          Preço *
          <input v-model="form.price" inputmode="decimal" required />
          <span v-if="errors.price" class="field-error">{{ errors.price }}</span>
        </label>
        <label>
          Estoque
          <input v-model.number="form.stock" type="number" min="0" step="1" />
          <span v-if="errors.stock" class="field-error">{{ errors.stock }}</span>
        </label>
        <label>
          Unidade
          <input v-model="form.unit" placeholder="UN" />
          <span v-if="errors.unit" class="field-error">{{ errors.unit }}</span>
        </label>
        <label>
          Código de barras
          <input v-model="form.barcode" />
        </label>
        <label class="checkbox">
          <input v-model="form.active" type="checkbox" />
          Ativo
        </label>
        <p v-if="submitError" class="submit-error">{{ submitError }}</p>
        <div class="actions">
          <button type="button" class="secondary" @click="emit('close')">Cancelar</button>
          <AppButton type="submit" :disabled="isPending">Salvar</AppButton>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import type { Product } from '@/entities/product/model/types'
import { editProductSchema } from '../model/schema'
import { useUpdateProduct } from '../model/use-edit-product'
import { getApiErrorMessage } from '@/shared/api/errors'
import AppButton from '@/shared/ui/AppButton.vue'

const props = defineProps<{ visible: boolean; product: Product | null }>()
const emit = defineEmits<{ close: [] }>()

const form = reactive({
  name: '',
  sku: '',
  price: '',
  stock: 0,
  unit: 'UN',
  barcode: '',
  active: true,
})
const errors = ref<Record<string, string>>({})
const submitError = ref('')

const { mutate, isPending } = useUpdateProduct()

watch(
  () => [props.visible, props.product] as const,
  ([open, product]) => {
    if (open && product) {
      form.name = product.name
      form.sku = product.sku
      form.price = product.price
      form.stock = product.stock
      form.unit = product.unit || 'UN'
      form.barcode = product.barcode ?? ''
      form.active = product.active
      errors.value = {}
      submitError.value = ''
    }
  },
  { immediate: true },
)

function normalizePrice(value: string) {
  return value.replace(',', '.')
}

function onSubmit() {
  if (!props.product) return
  errors.value = {}
  submitError.value = ''
  const parsed = editProductSchema.safeParse({
    name: form.name,
    sku: form.sku,
    price: normalizePrice(form.price),
    stock: form.stock,
    unit: form.unit || 'UN',
    barcode: form.barcode || undefined,
    active: form.active,
  })
  if (!parsed.success) {
    for (const issue of parsed.error.issues) {
      const key = String(issue.path[0] ?? '_')
      errors.value[key] = issue.message
    }
    return
  }
  mutate(
    {
      id: props.product.id,
      data: {
        name: parsed.data.name,
        sku: parsed.data.sku,
        price: normalizePrice(parsed.data.price),
        stock: parsed.data.stock,
        unit: parsed.data.unit,
        barcode: parsed.data.barcode,
        active: parsed.data.active,
      },
    },
    {
      onSuccess: () => emit('close'),
      onError: (error) => {
        submitError.value = getApiErrorMessage(error, 'Não foi possível atualizar o produto')
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
  max-width: 480px;
  padding: 24px;
  width: 100%;
}

h2 {
  font-size: 18px;
  margin: 0 0 20px;
}

.form {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

label {
  color: var(--color-text-secondary);
  display: flex;
  flex-direction: column;
  font-size: 13px;
  gap: 6px;
}

label.checkbox {
  flex-direction: row;
  align-items: center;
}

input[type='text'],
input[type='number'],
input:not([type]) {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font: inherit;
  padding: 8px 12px;
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
