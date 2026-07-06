import type { SessionUser } from '@/entities/session/model/types'

export type UserRole = SessionUser['role']

export function canWrite(role: UserRole | undefined | null): boolean {
  return role === 'admin' || role === 'manager' || role === 'operator'
}

export function canDelete(role: UserRole | undefined | null): boolean {
  return role === 'admin'
}

export function isAdmin(role: UserRole | undefined | null): boolean {
  return role === 'admin'
}

export function canManageSales(role: UserRole | undefined | null): boolean {
  return role === 'admin' || role === 'manager'
}

export function canViewFinancial(role: UserRole | undefined | null): boolean {
  return role === 'admin' || role === 'manager'
}
