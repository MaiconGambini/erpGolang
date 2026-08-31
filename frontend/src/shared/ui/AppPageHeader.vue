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
  align-items: flex-start;
  display: flex;
  gap: 20px;
  justify-content: space-between;
  margin-bottom: 24px;
}

h1 {
  font-size: clamp(24px, 2vw, 30px);
  font-weight: 680;
  letter-spacing: -0.03em;
  line-height: 1.15;
  margin: 0 0 6px;
}

p {
  color: var(--color-text-muted);
  font-size: 14px;
  line-height: 1.45;
  margin: 0;
  max-width: 68ch;
}

@media (max-width: 640px) {
  .page-header {
    flex-direction: column;
    gap: 14px;
    margin-bottom: 20px;
  }
}
</style>
