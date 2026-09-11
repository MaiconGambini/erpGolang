<template>
  <AppShell>
    <AppPageHeader title="Vendas" description="Pedidos de venda com confirmação e baixa de estoque">
      <template #actions>
        <div class="header-actions">
          <button type="button" class="outline" :disabled="exporting" @click="onExportCsv">
            Exportar CSV
          </button>
          <AppButton v-if="canWriteUser" @click="showCreate = true">Nova venda</AppButton>
        </div>
      </template>
    </AppPageHeader>
    <SaleTable ref="saleTableRef" @view="onView" @edit="onEdit" />
    <CreateSaleDialog :visible="showCreate" @close="showCreate = false" />
    <SaleDetailDialog
      :visible="!!viewingId"
      :sale-id="viewingId"
      @close="viewingId = null"
      @edit="onEditFromDetail"
    />
    <EditSaleDialog :visible="!!editingId" :sale-id="editingId" @close="editingId = null" />
  </AppShell>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useSessionStore } from '@/entities/session/model/session.store'
import CreateSaleDialog from '@/features/sale/create/ui/CreateSaleDialog.vue'
import SaleDetailDialog from '@/features/sale/detail/ui/SaleDetailDialog.vue'
import EditSaleDialog from '@/features/sale/edit/ui/EditSaleDialog.vue'
import AppButton from '@/shared/ui/AppButton.vue'
import AppPageHeader from '@/shared/ui/AppPageHeader.vue'
import { canWrite } from '@/shared/lib/roles'
import { useCsvExport } from '@/shared/lib/use-csv-export'
import AppShell from '@/widgets/app-shell/ui/AppShell.vue'
import SaleTable from '@/widgets/sale-table/ui/SaleTable.vue'

const session = useSessionStore()
const canWriteUser = computed(() => canWrite(session.user?.role))
const { exporting, exportCsv } = useCsvExport('/sales', 'vendas.csv')
const saleTableRef = ref<InstanceType<typeof SaleTable> | null>(null)

function onExportCsv() {
  exportCsv(saleTableRef.value?.exportParams ?? {})
}

const showCreate = ref(false)
const viewingId = ref<string | null>(null)
const editingId = ref<string | null>(null)

function onView(id: string) {
  viewingId.value = id
}

function onEdit(id: string) {
  editingId.value = id
}

function onEditFromDetail(id: string) {
  viewingId.value = null
  editingId.value = id
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
