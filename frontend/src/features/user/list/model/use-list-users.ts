import { useQuery } from '@tanstack/vue-query'
import { computed, type Ref } from 'vue'
import { listUsers } from '@/entities/user/api/user.api'

export function useListUsers(filters: {
  limit: Ref<number>
  offset: Ref<number>
}) {
  return useQuery({
    queryKey: computed(() => ['users', {
      limit: filters.limit.value,
      offset: filters.offset.value,
    }]),
    queryFn: () => listUsers({
      limit: filters.limit.value,
      offset: filters.offset.value,
    }),
  })
}
