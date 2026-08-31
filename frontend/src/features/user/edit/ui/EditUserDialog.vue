<template>
  <AppDialog v-if="visible && user" :visible="true" title-id="edit-user-title" @close="emit('close')">
      <h2 id="edit-user-title">Editar usuário</h2>
      <form class="form" @submit.prevent="onSubmit">
        <label>
          Nome *
          <input v-model="form.name" required />
        </label>
        <label>
          E-mail
          <input :value="user.email" type="email" disabled />
        </label>
        <label>
          Perfil *
          <select v-model="form.role" required>
            <option value="admin">Administrador</option>
            <option value="manager">Gerente</option>
            <option value="operator">Operador</option>
            <option value="viewer">Visualizador</option>
          </select>
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
  </AppDialog>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import type { User } from '@/entities/user/model/types'
import { useUpdateUser } from '../model/use-edit-user'
import AppDialog from '@/shared/ui/AppDialog.vue'
import AppButton from '@/shared/ui/AppButton.vue'
import { getApiErrorMessage } from '@/shared/api/errors'

const props = defineProps<{ visible: boolean; user: User | null }>()
const emit = defineEmits<{ close: [] }>()

const form = reactive({
  name: '',
  role: 'operator' as User['role'],
  active: true,
})
const submitError = ref('')

const { mutate, isPending } = useUpdateUser()

watch(
  () => [props.visible, props.user] as const,
  ([open, user]) => {
    if (open && user) {
      form.name = user.name
      form.role = user.role
      form.active = user.active
      submitError.value = ''
    }
  },
  { immediate: true },
)

function onSubmit() {
  if (!props.user) return
  submitError.value = ''
  mutate(
    {
      id: props.user.id,
      data: {
        name: form.name.trim(),
        role: form.role,
        active: form.active,
      },
    },
    {
      onSuccess: () => emit('close'),
      onError: (error) => {
        submitError.value = getApiErrorMessage(error, 'Não foi possível atualizar o usuário')
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

label {
  color: var(--color-text-secondary);
  display: flex;
  flex-direction: column;
  font-size: 13px;
  gap: 6px;
}

label.checkbox {
  align-items: center;
  flex-direction: row;
}

input,
select {
  background: var(--color-surface-raised);
  border: 1px solid var(--color-border-strong);
  border-radius: var(--radius-md);
  color: var(--color-text-primary);
  font: inherit;
  min-height: 40px;
  outline: none;
  padding: 8px 12px;
}

input:focus,
select:focus {
  border-color: var(--color-brand);
  box-shadow: var(--focus-ring);
}

input:disabled {
  background: var(--color-surface-muted);
  color: var(--color-text-muted);
}

.submit-error {
  color: var(--color-danger);
  font-size: 12px;
  margin: 0;
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
