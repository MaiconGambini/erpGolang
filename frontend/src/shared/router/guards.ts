import type { NavigationGuardNext, RouteLocationNormalized } from 'vue-router'
import { useSessionStore } from '@/entities/session/model/session.store'

export function requireAuth(
  to: RouteLocationNormalized,
  _from: RouteLocationNormalized,
  next: NavigationGuardNext,
) {
  const session = useSessionStore()
  if (!session.isAuthenticated) {
    next({ path: '/login', query: { redirect: to.fullPath } })
    return
  }
  next()
}
