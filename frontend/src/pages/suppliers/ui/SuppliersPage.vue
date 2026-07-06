<template>
  <AppShell>
    <AppPageHeader title="Fornecedores" description="Gerencie fornecedores de produtos e serviços do tenant">
      <template #actions>
        <div class="header-actions">
          <button type="button" class="outline" :disabled="exporting" @click="exportCsv()">
            Exportar CSV
          </button>
          <AppButton v-if="canWriteUser" @click="showCreate = true">Novo fornecedor</AppButton>
        </div>
      </template>
    </AppPageHeader>
    <SupplierTable @edit="onEdit" @delete="onDelete" />
    <CreateSupplierDialog :visible="showCreate" @close="showCreate = false" />
    <EditSupplierDialog :visible="!!editing" :supplier="editing" @close="editing = null" />
    <DeleteSupplierDialog :visible="!!deleting" :supplier="deleting" @close="deleting = null" />
  </AppShell>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { Supplier } from '@/entities/supplier/model/types'
import { useSessionStore } from '@/entities/session/model/session.store'
import CreateSupplierDialog from '@/features/supplier/create/ui/CreateSupplierDialog.vue'
import DeleteSupplierDialog from '@/features/supplier/delete/ui/DeleteSupplierDialog.vue'
import EditSupplierDialog from '@/features/supplier/edit/ui/EditSupplierDialog.vue'
import AppButton from '@/shared/ui/AppButton.vue'
import AppPageHeader from '@/shared/ui/AppPageHeader.vue'
import { canWrite } from '@/shared/lib/roles'
import { useCsvExport } from '@/shared/lib/use-csv-export'
import AppShell from '@/widgets/app-shell/ui/AppShell.vue'
import SupplierTable from '@/widgets/supplier-table/ui/SupplierTable.vue'

const session = useSessionStore()
const canWriteUser = computed(() => canWrite(session.user?.role))
const { exporting, exportCsv } = useCsvExport('/suppliers', 'fornecedores.csv')

const showCreate = ref(false)
const editing = ref<Supplier | null>(null)
const deleting = ref<Supplier | null>(null)

function onEdit(supplier: Supplier) {
  editing.value = supplier
}

function onDelete(supplier: Supplier) {
  deleting.value = supplier
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
