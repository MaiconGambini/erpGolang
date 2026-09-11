<template>
  <AppShell>
    <AppPageHeader title="Usuários" description="Gerencie perfis e acesso dos usuários do tenant" />
    <section class="panel">
      <TableSkeleton v-if="isLoading" :cols="5" />
      <p v-else-if="isError" class="state error">Erro ao carregar usuários</p>
      <EmptyState
        v-else-if="!users.length"
        :icon="UserCog"
        title="Nenhum usuário encontrado"
        hint="Usuários do tenant aparecem aqui após o convite ou criação."
      />
      <table v-else>
        <thead>
          <tr>
            <th>Nome</th>
            <th>E-mail</th>
            <th>Perfil</th>
            <th>Status</th>
            <th class="actions-col">Ações</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id">
            <td>{{ user.name }}</td>
            <td>{{ user.email }}</td>
            <td>{{ roleLabel(user.role) }}</td>
            <td>
              <span class="badge" :class="{ inactive: !user.active }">
                {{ user.active ? 'Ativo' : 'Inativo' }}
              </span>
            </td>
            <td class="actions-col">
              <button type="button" class="link" @click="editing = user">Editar</button>
            </td>
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
    <EditUserDialog :visible="!!editing" :user="editing" @close="editing = null" />
  </AppShell>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { User } from '@/entities/user/model/types'
import EditUserDialog from '@/features/user/edit/ui/EditUserDialog.vue'
import { useListUsers } from '@/features/user/list/model/use-list-users'
import AppPageHeader from '@/shared/ui/AppPageHeader.vue'
import AppShell from '@/widgets/app-shell/ui/AppShell.vue'
import EmptyState from '@/shared/ui/EmptyState.vue'
import TableSkeleton from '@/shared/ui/TableSkeleton.vue'
import { UserCog } from 'lucide-vue-next'
import type { UserRole } from '@/shared/lib/roles'

const limit = ref(20)
const offset = ref(0)
const editing = ref<User | null>(null)

const { data, isLoading, isError } = useListUsers({ limit, offset })

const users = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.pagination.total ?? 0)

function roleLabel(role: UserRole) {
  const map: Record<UserRole, string> = {
    admin: 'Administrador',
    manager: 'Gerente',
    operator: 'Operador',
    viewer: 'Visualizador',
  }
  return map[role]
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
  min-width: 540px;
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

.actions-col {
  white-space: nowrap;
  width: 100px;
}

.link {
  background: none;
  border: 0;
  border-radius: var(--radius-sm);
  color: var(--color-brand);
  cursor: pointer;
  font: inherit;
  font-size: 13px;
  min-height: 32px;
  padding: 0 5px;
  transition: background-color 0.15s ease, color 0.15s ease;
}

.link:hover {
  background: var(--color-brand-soft);
  color: var(--color-brand-hover);
  text-decoration: underline;
  text-underline-offset: 3px;
}

.badge {
  background: var(--color-success-soft);
  border-radius: 999px;
  color: var(--color-success-strong);
  font-size: 12px;
  font-weight: 650;
  line-height: 1;
  padding: 6px 9px;
  white-space: nowrap;
}

.badge.inactive {
  background: var(--color-surface-muted);
  color: var(--color-text-muted);
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
