import { useMutation } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'
import { login } from '@/entities/session/api/session.api'
import { useSessionStore } from '@/entities/session/model/session.store'

export function useLogin() {
  const session = useSessionStore()
  const router = useRouter()

  return useMutation({
    mutationFn: login,
    onSuccess: (data) => {
      session.setAccess(data.accessToken)
      session.setUser(data.user)
      router.push('/')
    },
  })
}
