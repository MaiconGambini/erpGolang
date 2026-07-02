<template>
  <button class="app-button" :type="type">
    <slot />
  </button>
</template>

<script setup lang="ts">
withDefaults(defineProps<{ type?: 'button' | 'submit' | 'reset' }>(), { type: 'button' })
</script>

<style scoped>
.app-button {
  border: 0;
  border-radius: var(--radius-md);
  background: var(--color-brand);
  color: #fff;
  cursor: pointer;
  font: inherit;
  font-weight: 600;
  min-height: 38px;
  padding: 0 16px;
}

.app-button:focus-visible {
  outline: none;
  box-shadow: var(--focus-ring);
}
</style>
