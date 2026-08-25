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

.actions-col {
  white-space: nowrap;
  width: 100px;
}

.link {
  background: none;
  border: 0;
  color: var(--color-brand);
  cursor: pointer;
  font: inherit;
  font-size: 13px;
  padding: 0;
}

.badge {
  background: #dcfce7;
  border-radius: 999px;
  color: #15803d;
  padding: 3px 10px;
}

.badge.inactive {
  background: #f3f4f6;
  color: #6b7280;
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
