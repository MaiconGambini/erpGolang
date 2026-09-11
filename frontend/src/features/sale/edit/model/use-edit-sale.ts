import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { updateSale, type SaleInput } from '@/entities/sale/api/sale.api'
import { invalidateDashboardSummary } from '@/features/dashboard/summary/model/use-dashboard-summary'

export function useEditSale() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: SaleInput }) => updateSale(id, data),
    onSuccess: async (_, { id }) => {
      await queryClient.invalidateQueries({ queryKey: ['sales'] })
      await queryClient.invalidateQueries({ queryKey: ['sale', id] })
      await queryClient.invalidateQueries({ queryKey: ['products'] })
      await invalidateDashboardSummary(queryClient)
    },
  })
}
