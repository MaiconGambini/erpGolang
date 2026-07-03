<template>
  <AppShell>
    <AppPageHeader title="Fornecedores" description="Gerencie fornecedores de produtos e serviços do tenant">
      <template #actions>
        <AppButton @click="showCreate = true">Novo fornecedor</AppButton>
      </template>
    </AppPageHeader>
    <SupplierTable @edit="onEdit" @delete="onDelete" />
    <CreateSupplierDialog :visible="showCreate" @close="showCreate = false" />
    <EditSupplierDialog :visible="!!editing" :supplier="editing" @close="editing = null" />
    <DeleteSupplierDialog :visible="!!deleting" :supplier="deleting" @close="deleting = null" />
  </AppShell>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { Supplier } from '@/entities/supplier/model/types'
import CreateSupplierDialog from '@/features/supplier/create/ui/CreateSupplierDialog.vue'
import DeleteSupplierDialog from '@/features/supplier/delete/ui/DeleteSupplierDialog.vue'
import EditSupplierDialog from '@/features/supplier/edit/ui/EditSupplierDialog.vue'
import AppButton from '@/shared/ui/AppButton.vue'
import AppPageHeader from '@/shared/ui/AppPageHeader.vue'
import AppShell from '@/widgets/app-shell/ui/AppShell.vue'
import SupplierTable from '@/widgets/supplier-table/ui/SupplierTable.vue'

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
