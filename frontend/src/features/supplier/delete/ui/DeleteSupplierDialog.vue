<template>
  <div v-if="visible && supplier" class="overlay" @click.self="emit('close')">
    <div class="dialog" role="dialog" aria-labelledby="delete-title">
      <h2 id="delete-title">Excluir fornecedor</h2>
      <p>
        Tem certeza que deseja excluir <strong>{{ supplier.name }}</strong>?
        Esta ação não pode ser desfeita.
      </p>
      <p v-if="submitError" class="submit-error">{{ submitError }}</p>
      <div class="actions">
        <button type="button" class="secondary" @click="emit('close')">Cancelar</button>
        <button type="button" class="danger" :disabled="isPending" @click="onConfirm">
          Excluir
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { Supplier } from '@/entities/supplier/model/types'
import { getApiErrorMessage } from '@/shared/api/errors'
import { useDeleteSupplier } from '../model/use-delete-supplier'

const props = defineProps<{ visible: boolean; supplier: Supplier | null }>()
const emit = defineEmits<{ close: [] }>()

const submitError = ref('')
const { mutate, isPending } = useDeleteSupplier()

watch(() => props.visible, (open) => {
  if (open) submitError.value = ''
})

function onConfirm() {
  if (!props.supplier) return
  mutate(props.supplier.id, {
    onSuccess: () => emit('close'),
    onError: (error) => {
      submitError.value = getApiErrorMessage(error, 'Não foi possível excluir o fornecedor')
    },
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
  max-width: 420px;
  padding: 24px;
  width: 100%;
}

h2 {
  font-size: 18px;
  margin: 0 0 12px;
}

p {
  color: var(--color-text-secondary);
  margin: 0 0 20px;
}

.submit-error {
  color: #dc2626;
  font-size: 13px;
}

.actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}

.secondary,
.danger {
  border-radius: var(--radius-md);
  cursor: pointer;
  font: inherit;
  padding: 8px 14px;
}

.secondary {
  background: transparent;
  border: 1px solid var(--color-border);
}

.danger {
  background: #dc2626;
  border: 0;
  color: #fff;
  font-weight: 600;
}

.danger:disabled {
  opacity: 0.6;
}
</style>
