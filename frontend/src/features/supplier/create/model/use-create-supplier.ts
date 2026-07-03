import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { createSupplier, type SupplierInput } from '@/entities/supplier/api/supplier.api'

export function useCreateSupplier() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: SupplierInput) => createSupplier(data),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['suppliers'] })
    },
  })
}
