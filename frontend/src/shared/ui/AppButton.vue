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
  align-items: center;
  border: 1px solid transparent;
  border-radius: var(--radius-md);
  background: var(--color-brand);
  color: var(--color-text-on-brand);
  cursor: pointer;
  display: inline-flex;
  font: inherit;
  font-weight: 650;
  justify-content: center;
  min-height: 40px;
  padding: 0 16px;
  transition: background-color 0.15s ease, box-shadow 0.15s ease, transform 0.15s ease;
}

.app-button:hover:not(:disabled) {
  background: var(--color-brand-hover);
  box-shadow: var(--shadow-soft);
}

.app-button:active:not(:disabled) {
  transform: translateY(1px);
}

.app-button:disabled {
  cursor: not-allowed;
  opacity: 0.55;
}

.app-button:focus-visible {
  outline: none;
  box-shadow: var(--focus-ring);
}
</style>
