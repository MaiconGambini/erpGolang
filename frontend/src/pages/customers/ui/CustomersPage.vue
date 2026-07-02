<template>
  <AppShell>
    <AppPageHeader title="Clientes" description="Gerencie pessoas físicas e jurídicas vinculadas ao tenant">
      <template #actions>
        <AppButton>Novo cliente</AppButton>
      </template>
    </AppPageHeader>
    <CustomerTable />
  </AppShell>
</template>

<script setup lang="ts">
import AppButton from '@/shared/ui/AppButton.vue'
import AppPageHeader from '@/shared/ui/AppPageHeader.vue'
import AppShell from '@/widgets/app-shell/ui/AppShell.vue'
import CustomerTable from '@/widgets/customer-table/ui/CustomerTable.vue'
</script>
