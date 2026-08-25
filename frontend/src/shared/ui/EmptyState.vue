<template>
  <div class="empty-state">
    <component :is="icon" v-if="icon" :size="30" :stroke-width="1.75" class="empty-icon" aria-hidden="true" />
    <p class="empty-title">{{ title }}</p>
    <p v-if="hint" class="empty-hint">{{ hint }}</p>
    <div v-if="$slots.default" class="empty-action">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Component } from 'vue'

defineProps<{
  icon?: Component
  title: string
  hint?: string
}>()
</script>

<style scoped>
.empty-state {
  align-items: center;
  color: var(--color-text-muted);
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 40px 16px;
  text-align: center;
}

.empty-icon {
  color: var(--color-text-muted);
  margin-bottom: 4px;
  opacity: 0.7;
}

.empty-title {
  color: var(--color-text-primary);
  font-size: 15px;
  font-weight: 600;
  margin: 0;
}

.empty-hint {
  font-size: 13px;
  margin: 0;
  max-width: 40ch;
}

.empty-action {
  margin-top: 10px;
}
</style>
