import { apiClient } from '@/shared/api/client'
import type { Paginated } from '@/shared/api/types'
import type { User, UserInput } from '../model/types'

export interface ListUsersParams {
  limit?: number
  offset?: number
}

export async function listUsers(params: ListUsersParams = {}): Promise<Paginated<User>> {
  const response = await apiClient.get<Paginated<User>>('/users', { params })
  return response.data
}

export async function getUser(id: string): Promise<User> {
  const response = await apiClient.get<{ data: User }>(`/users/${id}`)
  return response.data.data
}

export async function updateUser(id: string, data: UserInput): Promise<User> {
  const response = await apiClient.patch<{ data: User }>(`/users/${id}`, data)
  return response.data.data
}
