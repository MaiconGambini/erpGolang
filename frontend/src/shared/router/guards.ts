import type { NavigationGuardNext, RouteLocationNormalized } from 'vue-router'
import { useSessionStore } from '@/entities/session/model/session.store'
import { isAdmin } from '@/shared/lib/roles'

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

export function requireAdmin(
  to: RouteLocationNormalized,
  _from: RouteLocationNormalized,
  next: NavigationGuardNext,
) {
  const session = useSessionStore()
  if (!session.isAuthenticated) {
    next({ path: '/login', query: { redirect: to.fullPath } })
    return
  }
  if (!isAdmin(session.user?.role)) {
    next({ path: '/' })
    return
  }
  next()
}
