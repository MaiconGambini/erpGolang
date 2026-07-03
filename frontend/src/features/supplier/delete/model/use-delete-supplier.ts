import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { removeSupplier } from '@/entities/supplier/api/supplier.api'

export function useDeleteSupplier() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => removeSupplier(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['suppliers'] })
    },
  })
}
