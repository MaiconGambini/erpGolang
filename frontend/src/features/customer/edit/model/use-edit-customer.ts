import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { updateCustomer, type CustomerInput } from '@/entities/customer/api/customer.api'

import { invalidateDashboardSummary } from '@/features/dashboard/summary/model/use-dashboard-summary'

export function useUpdateCustomer() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: CustomerInput }) => updateCustomer(id, data),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['customers'] })
      await invalidateDashboardSummary(queryClient)
    },
  })
}
