import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { removeProduct } from '@/entities/product/api/product.api'

export function useDeleteProduct() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => removeProduct(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['products'] })
    },
  })
}
