import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { updateSupplier, type SupplierInput } from '@/entities/supplier/api/supplier.api'

export function useUpdateSupplier() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: SupplierInput }) => updateSupplier(id, data),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['suppliers'] })
    },
  })
}
