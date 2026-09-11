<template>
  <div class="table-skeleton" role="status" :aria-label="label">
    <div v-for="row in rows" :key="row" class="sk-row" aria-hidden="true">
      <span
        v-for="col in cols"
        :key="col"
        class="skeleton sk-cell"
        :style="{ maxWidth: `${55 + ((row * 7 + col * 13) % 40)}%` }"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    rows?: number
    cols?: number
    label?: string
  }>(),
  { rows: 6, cols: 4, label: 'Carregando dados' },
)
</script>
