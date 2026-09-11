<template>
  <fieldset class="party-fields">
    <legend>Documento</legend>
    <label>
      Tipo
      <select v-model="form.documentType">
        <option value="">—</option>
        <option value="cpf">CPF</option>
        <option value="cnpj">CNPJ</option>
      </select>
    </label>
    <label>
      Documento
      <input
        :value="form.document"
        @input="onDocumentInput"
      />
      <span v-if="errors.document" class="field-error">{{ errors.document }}</span>
    </label>
  </fieldset>
  <fieldset class="party-fields">
    <legend>Endereço</legend>
    <label>
      CEP
      <input
        :value="form.postalCode"
        @input="onPostalCodeInput"
      />
    </label>
    <label>
      Rua
      <input v-model="form.street" />
    </label>
    <label>
      Número
      <input v-model="form.streetNumber" />
    </label>
    <label>
      Cidade
      <input v-model="form.city" />
    </label>
    <label>
      UF
      <input v-model="form.state" maxlength="2" />
      <span v-if="errors.state" class="field-error">{{ errors.state }}</span>
    </label>
  </fieldset>
</template>

<script setup lang="ts">
import type { DocumentType } from '@/shared/lib/document'
import { formatDocument, formatPostalCode } from '@/shared/lib/document'

export interface PartyFormState {
  documentType: DocumentType | ''
  document: string
  postalCode: string
  street: string
  streetNumber: string
  city: string
  state: string
}

const form = defineModel<PartyFormState>({ required: true })

defineProps<{
  errors: Record<string, string>
}>()

function onDocumentInput(event: Event) {
  const target = event.target as HTMLInputElement
  form.value = {
    ...form.value,
    document: formatDocument(target.value, form.value.documentType || undefined),
  }
}

function onPostalCodeInput(event: Event) {
  const target = event.target as HTMLInputElement
  form.value = {
    ...form.value,
    postalCode: formatPostalCode(target.value),
  }
}
</script>

<style scoped>
.party-fields {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  display: flex;
  flex-direction: column;
  gap: 14px;
  margin: 0;
  padding: 14px;
}

legend {
  color: var(--color-text-secondary);
  font-size: 12px;
  padding: 0 4px;
}

label {
  color: var(--color-text-secondary);
  display: flex;
  flex-direction: column;
  font-size: 13px;
  gap: 6px;
}

select,
input {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font: inherit;
  padding: 8px 12px;
}

.field-error {
  color: var(--color-danger);
  font-size: 12px;
}
</style>
