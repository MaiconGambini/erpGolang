<template>
  <button type="button" class="logout-button" @click="onLogout">Sair</button>
</template>

<script setup lang="ts">
import { useLogout } from '../model/use-logout'

const onLogout = useLogout()
</script>

<style scoped>
.logout-button {
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  cursor: pointer;
  font: inherit;
  font-size: 13px;
  font-weight: 600;
  min-height: 38px;
  padding: 0 13px;
  transition: background-color 0.15s ease, border-color 0.15s ease, color 0.15s ease;
}

.logout-button:hover {
  background: var(--color-danger-soft);
  border-color: var(--color-danger);
  color: var(--color-danger-strong);
}
</style>
