<template>
  <header class="page-header">
    <div>
      <h1>{{ title }}</h1>
      <p v-if="description">{{ description }}</p>
    </div>
    <slot name="actions" />
  </header>
</template>

<script setup lang="ts">
defineProps<{ title: string; description?: string }>()
</script>

<style scoped>
.page-header {
  align-items: flex-end;
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

h1 {
  font-size: 26px;
  font-weight: 600;
  margin: 0 0 4px;
}

p {
  color: var(--color-text-muted);
  font-size: 14px;
  margin: 0;
}
</style>
