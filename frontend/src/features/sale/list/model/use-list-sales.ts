import { useQuery } from '@tanstack/vue-query'
import { computed, type Ref } from 'vue'
import { listSales } from '@/entities/sale/api/sale.api'
import type { SaleStatus } from '@/entities/sale/model/types'

export function useListSales(filters: {
  search: Ref<string>
  status: Ref<SaleStatus | undefined>
  from: Ref<string | undefined>
  to: Ref<string | undefined>
  limit: Ref<number>
  offset: Ref<number>
}) {
  return useQuery({
    queryKey: computed(() => ['sales', {
      search: filters.search.value,
      status: filters.status.value,
      from: filters.from.value,
      to: filters.to.value,
      limit: filters.limit.value,
      offset: filters.offset.value,
    }]),
    queryFn: () => listSales({
      search: filters.search.value,
      status: filters.status.value,
      from: filters.from.value,
      to: filters.to.value,
      limit: filters.limit.value,
      offset: filters.offset.value,
    }),
  })
}
