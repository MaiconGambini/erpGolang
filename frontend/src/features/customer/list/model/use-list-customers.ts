import { useQuery } from '@tanstack/vue-query'
import { computed, type Ref } from 'vue'
import { listCustomers } from '@/entities/customer/api/customer.api'

export function useListCustomers(filters: {
  search: Ref<string>
  active: Ref<boolean | undefined>
  limit: Ref<number>
  offset: Ref<number>
}) {
  return useQuery({
    queryKey: computed(() => ['customers', {
      search: filters.search.value,
      active: filters.active.value,
      limit: filters.limit.value,
      offset: filters.offset.value,
    }]),
    queryFn: () => listCustomers({
      search: filters.search.value,
      active: filters.active.value,
      limit: filters.limit.value,
      offset: filters.offset.value,
    }),
  })
}
