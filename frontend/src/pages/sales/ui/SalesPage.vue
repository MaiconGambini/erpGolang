<template>
  <AppShell>
    <AppPageHeader title="Vendas" description="Pedidos de venda com confirmação e baixa de estoque">
      <template #actions>
        <AppButton @click="showCreate = true">Nova venda</AppButton>
      </template>
    </AppPageHeader>
    <SaleTable />
    <CreateSaleDialog :visible="showCreate" @close="showCreate = false" />
  </AppShell>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import CreateSaleDialog from '@/features/sale/create/ui/CreateSaleDialog.vue'
import AppButton from '@/shared/ui/AppButton.vue'
import AppPageHeader from '@/shared/ui/AppPageHeader.vue'
import AppShell from '@/widgets/app-shell/ui/AppShell.vue'
import SaleTable from '@/widgets/sale-table/ui/SaleTable.vue'

const showCreate = ref(false)
</script>
