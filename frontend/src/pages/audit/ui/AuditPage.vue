<template>
  <AppShell>
    <AppPageHeader title="Auditoria" description="Registro de ações realizadas no tenant" />
    <section class="panel">
      <TableSkeleton v-if="isLoading" :cols="5" />
      <p v-else-if="isError" class="state error">Erro ao carregar auditoria</p>
      <EmptyState
        v-else-if="!logs.length"
        :icon="ScrollText"
        title="Nenhum registro de auditoria"
        hint="Ações de criação, edição e exclusão aparecem aqui automaticamente."
      />
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
import EmptyState from '@/shared/ui/EmptyState.vue'
import TableSkeleton from '@/shared/ui/TableSkeleton.vue'
import { ScrollText } from 'lucide-vue-next'
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
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-card);
  overflow-x: auto;
}

table {
  border-collapse: collapse;
  min-width: 700px;
  width: 100%;
}

th,
td {
  border-top: 1px solid var(--color-border);
  font-size: 13px;
  padding: 13px 16px;
  text-align: left;
}

th {
  background: var(--color-surface-muted);
  color: var(--color-text-muted);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  white-space: nowrap;
}

tbody tr {
  transition: background-color 0.15s ease;
}

tbody tr:hover {
  background: var(--color-surface-muted);
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
  color: var(--color-danger);
}

.pagination {
  align-items: center;
  background: var(--color-surface-muted);
  border-top: 1px solid var(--color-border);
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  justify-content: flex-end;
  padding: 12px 16px;
}

.pagination button {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-secondary);
  cursor: pointer;
  font: inherit;
  min-height: 34px;
  padding: 0 12px;
  transition: background-color 0.15s ease, border-color 0.15s ease;
}

.pagination button:hover:not(:disabled) {
  background: var(--color-surface-raised);
  border-color: var(--color-border-strong);
}

.pagination button:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}
</style>
