import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { createProduct, type ProductInput } from '@/entities/product/api/product.api'

import { invalidateDashboardSummary } from '@/features/dashboard/summary/model/use-dashboard-summary'

export function useCreateProduct() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: ProductInput) => createProduct(data),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['products'] })
      await invalidateDashboardSummary(queryClient)
    },
  })
}
