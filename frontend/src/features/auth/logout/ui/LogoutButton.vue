<template>
  <button type="button" @click="onLogout">Sair</button>
</template>

<script setup lang="ts">
import { useLogout } from '../model/use-logout'

const onLogout = useLogout()
</script>

<style scoped>
button {
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font: inherit;
  padding: 6px 12px;
}
</style>
