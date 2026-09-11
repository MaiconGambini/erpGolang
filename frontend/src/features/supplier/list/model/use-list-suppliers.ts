import { useQuery } from '@tanstack/vue-query'
import { computed, type Ref } from 'vue'
import { listSuppliers } from '@/entities/supplier/api/supplier.api'

export function useListSuppliers(filters: {
  search: Ref<string>
  active: Ref<boolean | undefined>
  limit: Ref<number>
  offset: Ref<number>
}) {
  return useQuery({
    queryKey: computed(() => ['suppliers', {
      search: filters.search.value,
      active: filters.active.value,
      limit: filters.limit.value,
      offset: filters.offset.value,
    }]),
    queryFn: () => listSuppliers({
      search: filters.search.value,
      active: filters.active.value,
      limit: filters.limit.value,
      offset: filters.offset.value,
    }),
  })
}
