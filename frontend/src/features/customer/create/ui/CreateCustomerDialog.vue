<template>
  <div v-if="visible" class="overlay" @click.self="emit('close')">
    <div class="dialog" role="dialog" aria-labelledby="create-title">
      <h2 id="create-title">Novo cliente</h2>
      <form class="form" @submit.prevent="onSubmit">
        <label>
          Nome *
          <input v-model="form.name" required />
          <span v-if="errors.name" class="field-error">{{ errors.name }}</span>
        </label>
        <PartyFormFields v-model="partyForm" :errors="errors" />
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
import { computed, reactive, ref, watch } from 'vue'
import { createCustomerSchema } from '../model/schema'
import { useCreateCustomer } from '../model/use-create-customer'
import AppButton from '@/shared/ui/AppButton.vue'
import PartyFormFields, { type PartyFormState } from '@/shared/ui/PartyFormFields.vue'
import { toPartyInput } from '@/shared/lib/party-payload'

const props = defineProps<{ visible: boolean }>()
const emit = defineEmits<{ close: [] }>()

const emptyPartyFields = (): PartyFormState => ({
  documentType: '',
  document: '',
  postalCode: '',
  street: '',
  streetNumber: '',
  city: '',
  state: '',
})

const form = reactive({
  name: '',
  email: '',
  phone: '',
  active: true,
  ...emptyPartyFields(),
})
const errors = ref<Record<string, string>>({})
const submitError = ref('')

const partyForm = computed({
  get: (): PartyFormState => ({
    documentType: form.documentType,
    document: form.document,
    postalCode: form.postalCode,
    street: form.street,
    streetNumber: form.streetNumber,
    city: form.city,
    state: form.state,
  }),
  set: (value: PartyFormState) => {
    Object.assign(form, value)
  },
})

const { mutate, isPending } = useCreateCustomer()

watch(() => props.visible, (open) => {
  if (open) {
    form.name = ''
    form.email = ''
    form.phone = ''
    form.active = true
    Object.assign(form, emptyPartyFields())
    errors.value = {}
    submitError.value = ''
  }
})

function onSubmit() {
  errors.value = {}
  submitError.value = ''
  const parsed = createCustomerSchema.safeParse(form)
  if (!parsed.success) {
    for (const issue of parsed.error.issues) {
      const key = String(issue.path[0] ?? '_')
      errors.value[key] = issue.message
    }
    return
  }
  mutate(toPartyInput(parsed.data), {
    onSuccess: () => emit('close'),
    onError: () => { submitError.value = 'Não foi possível criar o cliente' },
  })
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
  max-width: 520px;
  max-height: 90vh;
  overflow-y: auto;
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
