import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { createSale, type SaleInput } from '@/entities/sale/api/sale.api'

import { invalidateDashboardSummary } from '@/features/dashboard/summary/model/use-dashboard-summary'

export function useCreateSale() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: SaleInput) => createSale(data),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['sales'] })
      await queryClient.invalidateQueries({ queryKey: ['products'] })
      await invalidateDashboardSummary(queryClient)
    },
  })
}
