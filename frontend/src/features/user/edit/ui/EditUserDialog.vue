<template>
  <div v-if="visible && user" class="overlay" @click.self="emit('close')">
    <div class="dialog" role="dialog" aria-labelledby="edit-user-title">
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
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import type { User } from '@/entities/user/model/types'
import { useUpdateUser } from '../model/use-edit-user'
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
  align-items: center;
  flex-direction: row;
}

input,
select {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font: inherit;
  padding: 8px 12px;
}

input:disabled {
  background: #f9fafb;
  color: var(--color-text-muted);
}

.submit-error {
  color: #dc2626;
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
