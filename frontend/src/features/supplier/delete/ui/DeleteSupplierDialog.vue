<template>
  <AppDialog v-if="visible && supplier" :visible="true" title-id="delete-title" size="sm" @close="emit('close')">
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
  </AppDialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { Supplier } from '@/entities/supplier/model/types'
import { getApiErrorMessage } from '@/shared/api/errors'
import AppDialog from '@/shared/ui/AppDialog.vue'
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

h2 {
  font-size: 18px;
  font-weight: 650;
  margin: 0 0 12px;
}

p {
  color: var(--color-text-secondary);
  margin: 0 0 20px;
}

.submit-error {
  color: var(--color-danger);
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
  background: var(--color-danger);
  border: 0;
  color: var(--color-text-on-danger);
  font-weight: 600;
}

.danger:disabled {
  opacity: 0.6;
}
</style>
