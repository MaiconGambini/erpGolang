import { refresh } from '@/entities/session/api/session.api'
import { useSessionStore } from '@/entities/session/model/session.store'

export async function bootAuth() {
  const session = useSessionStore()
  if (session.accessToken) {
    return
  }
  try {
    const data = await refresh()
    session.setAccess(data.accessToken)
    session.setUser(data.user)
  } catch {
    session.clear()
  }
}
