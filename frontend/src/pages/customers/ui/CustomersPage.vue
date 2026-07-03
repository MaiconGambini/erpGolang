<template>
  <AppShell>
    <AppPageHeader title="Clientes" description="Gerencie pessoas físicas e jurídicas vinculadas ao tenant">
      <template #actions>
        <AppButton @click="showCreate = true">Novo cliente</AppButton>
      </template>
    </AppPageHeader>
    <CustomerTable @edit="onEdit" @delete="onDelete" />
    <CreateCustomerDialog :visible="showCreate" @close="showCreate = false" />
    <EditCustomerDialog :visible="!!editing" :customer="editing" @close="editing = null" />
    <DeleteCustomerDialog :visible="!!deleting" :customer="deleting" @close="deleting = null" />
  </AppShell>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { Customer } from '@/entities/customer/model/types'
import CreateCustomerDialog from '@/features/customer/create/ui/CreateCustomerDialog.vue'
import DeleteCustomerDialog from '@/features/customer/delete/ui/DeleteCustomerDialog.vue'
import EditCustomerDialog from '@/features/customer/edit/ui/EditCustomerDialog.vue'
import AppButton from '@/shared/ui/AppButton.vue'
import AppPageHeader from '@/shared/ui/AppPageHeader.vue'
import AppShell from '@/widgets/app-shell/ui/AppShell.vue'
import CustomerTable from '@/widgets/customer-table/ui/CustomerTable.vue'

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
