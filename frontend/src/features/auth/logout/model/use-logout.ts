import { useRouter } from 'vue-router'
import { useQueryClient } from '@tanstack/vue-query'
import { logout } from '@/entities/session/api/session.api'
import { useSessionStore } from '@/entities/session/model/session.store'

export function useLogout() {
  const session = useSessionStore()
  const router = useRouter()
  const queryClient = useQueryClient()

  return async () => {
    try {
      await logout()
    } finally {
      session.clear()
      queryClient.clear()
      router.push('/login')
    }
  }
}
