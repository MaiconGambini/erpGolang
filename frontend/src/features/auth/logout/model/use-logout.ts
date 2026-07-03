import { useRouter } from 'vue-router'
import { logout } from '@/entities/session/api/session.api'
import { useSessionStore } from '@/entities/session/model/session.store'

export function useLogout() {
  const session = useSessionStore()
  const router = useRouter()

  return async () => {
    try {
      await logout()
    } finally {
      session.clear()
      router.push('/login')
    }
  }
}
