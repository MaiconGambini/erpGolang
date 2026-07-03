<template>
  <div v-if="visible && customer" class="overlay" @click.self="emit('close')">
    <div class="dialog" role="dialog" aria-labelledby="edit-title">
      <h2 id="edit-title">Editar cliente</h2>
      <form class="form" @submit.prevent="onSubmit">
        <label>
          Nome *
          <input v-model="form.name" required />
          <span v-if="errors.name" class="field-error">{{ errors.name }}</span>
        </label>
        <label>
          Documento
          <input v-model="form.document" />
        </label>
        <label>
          E-mail
          <input v-model="form.email" type="email" />
          <span v-if="errors.email" class="field-error">{{ errors.email }}</span>
        </label>
        <label>
          Telefone
          <input v-model="form.phone" />
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
import type { Customer } from '@/entities/customer/model/types'
import { editCustomerSchema } from '../model/schema'
import { useUpdateCustomer } from '../model/use-edit-customer'
import AppButton from '@/shared/ui/AppButton.vue'

const props = defineProps<{ visible: boolean; customer: Customer | null }>()
const emit = defineEmits<{ close: [] }>()

const form = reactive({
  name: '',
  document: '',
  email: '',
  phone: '',
  active: true,
})
const errors = ref<Record<string, string>>({})
const submitError = ref('')

const { mutate, isPending } = useUpdateCustomer()

watch(
  () => [props.visible, props.customer] as const,
  ([open, customer]) => {
    if (open && customer) {
      form.name = customer.name
      form.document = customer.document ?? ''
      form.email = customer.email ?? ''
      form.phone = customer.phone ?? ''
      form.active = customer.active
      errors.value = {}
      submitError.value = ''
    }
  },
  { immediate: true },
)

function onSubmit() {
  if (!props.customer) return
  errors.value = {}
  submitError.value = ''
  const parsed = editCustomerSchema.safeParse({
    name: form.name,
    document: form.document || undefined,
    email: form.email || undefined,
    phone: form.phone || undefined,
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
      id: props.customer.id,
      data: {
        name: parsed.data.name,
        document: parsed.data.document,
        email: parsed.data.email,
        phone: parsed.data.phone,
        active: parsed.data.active,
      },
    },
    {
      onSuccess: () => emit('close'),
      onError: () => { submitError.value = 'Não foi possível atualizar o cliente' },
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
input[type='email'],
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
