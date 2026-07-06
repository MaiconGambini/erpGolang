import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { updateUser } from '@/entities/user/api/user.api'
import type { UserInput } from '@/entities/user/model/types'

export function useUpdateUser() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UserInput }) => updateUser(id, data),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['users'] })
    },
  })
}
