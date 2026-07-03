import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { updateProduct, type ProductInput } from '@/entities/product/api/product.api'

export function useUpdateProduct() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: ProductInput }) => updateProduct(id, data),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['products'] })
    },
  })
}
