import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { cancelSale, confirmSale, removeSale } from '@/entities/sale/api/sale.api'
import { invalidateDashboardSummary } from '@/features/dashboard/summary/model/use-dashboard-summary'

export function useConfirmSale() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => confirmSale(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['sales'] })
      await queryClient.invalidateQueries({ queryKey: ['products'] })
      await invalidateDashboardSummary(queryClient)
    },
  })
}

export function useCancelSale() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => cancelSale(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['sales'] })
      await queryClient.invalidateQueries({ queryKey: ['products'] })
      await invalidateDashboardSummary(queryClient)
    },
  })
}

export function useDeleteSale() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => removeSale(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['sales'] })
      await invalidateDashboardSummary(queryClient)
    },
  })
}
