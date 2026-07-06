<template>
  <AppShell>
    <AppPageHeader title="Auditoria" description="Registro de ações realizadas no tenant" />
    <section class="panel">
      <p v-if="isLoading" class="state">Carregando...</p>
      <p v-else-if="isError" class="state error">Erro ao carregar auditoria</p>
      <p v-else-if="!logs.length" class="state">Nenhum registro encontrado</p>
      <table v-else>
        <thead>
          <tr>
            <th>Data</th>
            <th>Ação</th>
            <th>Entidade</th>
            <th>ID da entidade</th>
            <th>Usuário</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="log in logs" :key="log.id">
            <td>{{ formatDate(log.createdAt) }}</td>
            <td>{{ log.action }}</td>
            <td>{{ log.entityType }}</td>
            <td class="mono">{{ log.entityId }}</td>
            <td class="mono">{{ log.actorUserId ?? '—' }}</td>
          </tr>
        </tbody>
      </table>
      <footer v-if="total > limit" class="pagination">
        <button type="button" :disabled="offset === 0" @click="offset = Math.max(0, offset - limit)">
          Anterior
        </button>
        <span>{{ offset + 1 }}–{{ Math.min(offset + limit, total) }} de {{ total }}</span>
        <button type="button" :disabled="offset + limit >= total" @click="offset += limit">
          Próximo
        </button>
      </footer>
    </section>
  </AppShell>
</template>

<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'
import { listAuditLogs } from '@/entities/audit/api/audit.api'
import AppPageHeader from '@/shared/ui/AppPageHeader.vue'
import AppShell from '@/widgets/app-shell/ui/AppShell.vue'

const limit = ref(50)
const offset = ref(0)

const { data, isLoading, isError } = useQuery({
  queryKey: computed(() => ['audit-logs', { limit: limit.value, offset: offset.value }]),
  queryFn: () => listAuditLogs({ limit: limit.value, offset: offset.value }),
})

const logs = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.pagination.total ?? 0)

function formatDate(value: string) {
  return new Date(value).toLocaleString('pt-BR')
}
</script>

<style scoped>
.panel {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 10px;
  overflow: hidden;
}

table {
  border-collapse: collapse;
  width: 100%;
}

th,
td {
  border-top: 1px solid var(--color-border);
  font-size: 13px;
  padding: 12px 16px;
  text-align: left;
}

.mono {
  font-family: ui-monospace, monospace;
  font-size: 12px;
}

.state {
  color: var(--color-text-muted);
  padding: 24px;
}

.state.error {
  color: #dc2626;
}

.pagination {
  align-items: center;
  border-top: 1px solid var(--color-border);
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  padding: 12px 16px;
}

.pagination button {
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font: inherit;
  padding: 6px 12px;
}

.pagination button:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}
</style>
