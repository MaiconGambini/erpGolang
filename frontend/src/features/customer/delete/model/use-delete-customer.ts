import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { removeCustomer } from '@/entities/customer/api/customer.api'

import { invalidateDashboardSummary } from '@/features/dashboard/summary/model/use-dashboard-summary'

export function useDeleteCustomer() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => removeCustomer(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['customers'] })
      await invalidateDashboardSummary(queryClient)
    },
  })
}
