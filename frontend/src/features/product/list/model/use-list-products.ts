import { useQuery } from '@tanstack/vue-query'
import { computed, type Ref } from 'vue'
import { listProducts } from '@/entities/product/api/product.api'

export function useListProducts(filters: {
  search: Ref<string>
  active: Ref<boolean | undefined>
  limit: Ref<number>
  offset: Ref<number>
}, options?: { enabled?: Ref<boolean> }) {
  return useQuery({
    queryKey: computed(() => ['products', {
      search: filters.search.value,
      active: filters.active.value,
      limit: filters.limit.value,
      offset: filters.offset.value,
    }]),
    queryFn: () => listProducts({
      search: filters.search.value,
      active: filters.active.value,
      limit: filters.limit.value,
      offset: filters.offset.value,
    }),
    enabled: computed(() => options?.enabled?.value !== false),
  })
}
