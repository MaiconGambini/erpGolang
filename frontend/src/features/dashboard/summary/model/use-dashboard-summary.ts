import { useQuery } from '@tanstack/vue-query'
import { storeToRefs } from 'pinia'
import { getDashboardSummary } from '@/entities/dashboard/api/dashboard.api'
import { useSessionStore } from '@/entities/session/model/session.store'

export const dashboardSummaryQueryKey = ['dashboard', 'summary'] as const

export function invalidateDashboardSummary(queryClient: { invalidateQueries: (opts: { queryKey: readonly string[] }) => Promise<void> }) {
  return queryClient.invalidateQueries({ queryKey: dashboardSummaryQueryKey })
}

export function useDashboardSummary() {
  const { user } = storeToRefs(useSessionStore())

  return useQuery({
    queryKey: [...dashboardSummaryQueryKey, user.value?.id ?? 'anonymous'],
    queryFn: () => getDashboardSummary(),
    enabled: () => !!user.value?.id,
  })
}
