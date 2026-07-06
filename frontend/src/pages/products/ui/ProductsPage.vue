<template>
  <AppShell>
    <AppPageHeader
      :title="pageTitle"
      :description="pageDescription"
    >
      <template #actions>
        <AppButton @click="showCreate = true">Novo produto</AppButton>
      </template>
    </AppPageHeader>
    <ProductTable :low-stock="isLowStockView" @edit="onEdit" @delete="onDelete" />
    <CreateProductDialog :visible="showCreate" @close="showCreate = false" />
    <EditProductDialog :visible="!!editing" :product="editing" @close="editing = null" />
    <DeleteProductDialog :visible="!!deleting" :product="deleting" @close="deleting = null" />
  </AppShell>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import type { Product } from '@/entities/product/model/types'
import CreateProductDialog from '@/features/product/create/ui/CreateProductDialog.vue'
import DeleteProductDialog from '@/features/product/delete/ui/DeleteProductDialog.vue'
import EditProductDialog from '@/features/product/edit/ui/EditProductDialog.vue'
import AppButton from '@/shared/ui/AppButton.vue'
import AppPageHeader from '@/shared/ui/AppPageHeader.vue'
import AppShell from '@/widgets/app-shell/ui/AppShell.vue'
import ProductTable from '@/widgets/product-table/ui/ProductTable.vue'

const route = useRoute()
const isLowStockView = computed(() => route.query.lowStock === '1')

const pageTitle = computed(() => (isLowStockView.value ? 'Estoque baixo' : 'Produtos'))
const pageDescription = computed(() =>
  isLowStockView.value
    ? 'Produtos com estoque abaixo do limite configurado'
    : 'Gerencie catálogo, preços e estoque do tenant',
)

const showCreate = ref(false)
const editing = ref<Product | null>(null)
const deleting = ref<Product | null>(null)

function onEdit(product: Product) {
  editing.value = product
}

function onDelete(product: Product) {
  deleting.value = product
}
</script>
