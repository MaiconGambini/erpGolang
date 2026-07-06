import { useQuery } from '@tanstack/vue-query'
import { computed, type Ref } from 'vue'
import { listLowStockProducts } from '@/entities/product/api/product.api'

export function useListLowStockProducts(enabled: Ref<boolean>) {
  return useQuery({
    queryKey: ['products', 'low-stock'],
    queryFn: () => listLowStockProducts(),
    enabled: computed(() => enabled.value),
  })
}
