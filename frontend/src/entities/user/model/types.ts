import type { SessionUser } from '@/entities/session/model/types'

export type UserRole = SessionUser['role']

export interface User {
  id: string
  email: string
  name: string
  role: UserRole
  active: boolean
  createdAt: string
  updatedAt: string
}

export interface UserInput {
  name: string
  role: UserRole
  active: boolean
}
