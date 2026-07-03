<template>
  <form class="form" @submit.prevent="onSubmit">
    <label>
      Empresa
      <input v-model="tenantSlug" placeholder="acme" required />
    </label>
    <label>
      E-mail
      <input v-model="email" type="email" placeholder="admin@acme.com" required />
    </label>
    <label>
      Senha
      <input v-model="password" type="password" placeholder="Sua senha" required />
    </label>
    <p v-if="error" class="error">{{ error }}</p>
    <AppButton type="submit" :disabled="isPending">Entrar</AppButton>
  </form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import AppButton from '@/shared/ui/AppButton.vue'
import { useLogin } from '../model/use-login'

const tenantSlug = ref('acme')
const email = ref('admin@acme.com')
const password = ref('admin123')
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
  gap: 16px;
}

label {
  color: var(--color-text-secondary);
  display: flex;
  flex-direction: column;
  font-size: 13px;
  gap: 6px;
}

input {
  border: 1px solid #cbd5e1;
  border-radius: var(--radius-md);
  font: inherit;
  padding: 10px 12px;
}

.error {
  color: #dc2626;
  font-size: 13px;
  margin: 0;
}
</style>
