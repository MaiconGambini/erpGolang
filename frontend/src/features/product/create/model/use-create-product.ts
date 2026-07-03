import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { createProduct, type ProductInput } from '@/entities/product/api/product.api'

export function useCreateProduct() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: ProductInput) => createProduct(data),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['products'] })
    },
  })
}
