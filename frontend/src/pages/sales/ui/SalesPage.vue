<template>
  <AppShell>
    <AppPageHeader title="Vendas" description="Pedidos de venda com confirmação e baixa de estoque">
      <template #actions>
        <AppButton @click="showCreate = true">Nova venda</AppButton>
      </template>
    </AppPageHeader>
    <SaleTable @view="onView" @edit="onEdit" />
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
import { ref } from 'vue'
import CreateSaleDialog from '@/features/sale/create/ui/CreateSaleDialog.vue'
import SaleDetailDialog from '@/features/sale/detail/ui/SaleDetailDialog.vue'
import EditSaleDialog from '@/features/sale/edit/ui/EditSaleDialog.vue'
import AppButton from '@/shared/ui/AppButton.vue'
import AppPageHeader from '@/shared/ui/AppPageHeader.vue'
import AppShell from '@/widgets/app-shell/ui/AppShell.vue'
import SaleTable from '@/widgets/sale-table/ui/SaleTable.vue'

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
