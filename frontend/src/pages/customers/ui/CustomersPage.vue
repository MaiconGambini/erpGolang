<template>
  <AppShell>
    <AppPageHeader title="Clientes" description="Gerencie pessoas físicas e jurídicas vinculadas ao tenant">
      <template #actions>
        <div class="header-actions">
          <button type="button" class="outline" :disabled="exporting" @click="exportCsv()">
            Exportar CSV
          </button>
          <AppButton v-if="canWriteUser" @click="showCreate = true">Novo cliente</AppButton>
        </div>
      </template>
    </AppPageHeader>
    <CustomerTable @edit="onEdit" @delete="onDelete" />
    <CreateCustomerDialog :visible="showCreate" @close="showCreate = false" />
    <EditCustomerDialog :visible="!!editing" :customer="editing" @close="editing = null" />
    <DeleteCustomerDialog :visible="!!deleting" :customer="deleting" @close="deleting = null" />
  </AppShell>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { Customer } from '@/entities/customer/model/types'
import { useSessionStore } from '@/entities/session/model/session.store'
import CreateCustomerDialog from '@/features/customer/create/ui/CreateCustomerDialog.vue'
import DeleteCustomerDialog from '@/features/customer/delete/ui/DeleteCustomerDialog.vue'
import EditCustomerDialog from '@/features/customer/edit/ui/EditCustomerDialog.vue'
import AppButton from '@/shared/ui/AppButton.vue'
import AppPageHeader from '@/shared/ui/AppPageHeader.vue'
import { canWrite } from '@/shared/lib/roles'
import { useCsvExport } from '@/shared/lib/use-csv-export'
import AppShell from '@/widgets/app-shell/ui/AppShell.vue'
import CustomerTable from '@/widgets/customer-table/ui/CustomerTable.vue'

const session = useSessionStore()
const canWriteUser = computed(() => canWrite(session.user?.role))
const { exporting, exportCsv } = useCsvExport('/customers', 'clientes.csv')

const showCreate = ref(false)
const editing = ref<Customer | null>(null)
const deleting = ref<Customer | null>(null)

function onEdit(customer: Customer) {
  editing.value = customer
}

function onDelete(customer: Customer) {
  deleting.value = customer
}
</script>

<style scoped>
.header-actions {
  display: flex;
  gap: 10px;
}

.outline {
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  cursor: pointer;
  font: inherit;
  font-weight: 600;
  min-height: 38px;
  padding: 0 16px;
}

.outline:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}
</style>
