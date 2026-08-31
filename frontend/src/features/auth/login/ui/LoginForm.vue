<template>
  <form class="form" @submit.prevent="onSubmit">
    <label>
      Empresa
      <input v-model="tenantSlug" autocomplete="organization" placeholder="acme" required />
    </label>
    <label>
      E-mail
      <input v-model="email" autocomplete="username" type="email" placeholder="admin@acme.com" required />
    </label>
    <label>
      Senha
      <input v-model="password" autocomplete="current-password" type="password" placeholder="Sua senha" required />
    </label>
    <p v-if="error" class="error" role="alert" aria-live="assertive">{{ error }}</p>
    <AppButton type="submit" :disabled="isPending">Entrar</AppButton>
  </form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import AppButton from '@/shared/ui/AppButton.vue'
import { useLogin } from '../model/use-login'

const tenantSlug = ref('')
const email = ref('')
const password = ref('')
const error = ref('')

const { mutate, isPending } = useLogin()

function onSubmit() {
  error.value = ''
  mutate(
    { tenantSlug: tenantSlug.value, email: email.value, password: password.value },
    { onError: () => { error.value = 'Credenciais inválidas' } },
  )
}
</script>

<style scoped>
.form {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

label {
  color: var(--color-text-secondary);
  display: flex;
  flex-direction: column;
  font-size: 13px;
  font-weight: 600;
  gap: 7px;
}

input {
  background: var(--color-surface-raised);
  border: 1px solid var(--color-border-strong);
  border-radius: var(--radius-md);
  color: var(--color-text-primary);
  min-height: 44px;
  outline: none;
  padding: 0 12px;
  transition: border-color 0.15s ease, box-shadow 0.15s ease, background-color 0.15s ease;
}

input:focus {
  background: var(--color-surface);
  border-color: var(--color-brand);
  box-shadow: var(--focus-ring);
}

.error {
  align-items: flex-start;
  background: var(--color-danger-soft);
  border: 1px solid color-mix(in srgb, var(--color-danger) 24%, transparent);
  border-radius: var(--radius-md);
  color: var(--color-danger-strong);
  display: flex;
  font-size: 13px;
  margin: 0;
  padding: 10px 12px;
}
</style>
