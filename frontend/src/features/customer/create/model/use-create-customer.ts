import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { createCustomer, type CustomerInput } from '@/entities/customer/api/customer.api'
import { invalidateDashboardSummary } from '@/features/dashboard/summary/model/use-dashboard-summary'

export function useCreateCustomer() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: CustomerInput) => createCustomer(data),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['customers'] })
      await invalidateDashboardSummary(queryClient)
    },
  })
}
