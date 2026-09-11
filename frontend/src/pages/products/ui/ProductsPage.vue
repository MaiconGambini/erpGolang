<template>
  <AppShell>
    <AppPageHeader
      :title="pageTitle"
      :description="pageDescription"
    >
      <template #actions>
        <div class="header-actions">
          <button
            v-if="!isLowStockView"
            type="button"
            class="outline"
            :disabled="exporting"
            @click="exportCsv()"
          >
            Exportar CSV
          </button>
          <AppButton v-if="canWriteUser" @click="showCreate = true">Novo produto</AppButton>
        </div>
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
import { useSessionStore } from '@/entities/session/model/session.store'
import CreateProductDialog from '@/features/product/create/ui/CreateProductDialog.vue'
import DeleteProductDialog from '@/features/product/delete/ui/DeleteProductDialog.vue'
import EditProductDialog from '@/features/product/edit/ui/EditProductDialog.vue'
import AppButton from '@/shared/ui/AppButton.vue'
import AppPageHeader from '@/shared/ui/AppPageHeader.vue'
import { canWrite } from '@/shared/lib/roles'
import { useCsvExport } from '@/shared/lib/use-csv-export'
import AppShell from '@/widgets/app-shell/ui/AppShell.vue'
import ProductTable from '@/widgets/product-table/ui/ProductTable.vue'

const route = useRoute()
const session = useSessionStore()
const canWriteUser = computed(() => canWrite(session.user?.role))
const { exporting, exportCsv } = useCsvExport('/products', 'produtos.csv')
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
