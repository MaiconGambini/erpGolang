import { apiClient } from '@/shared/api/client'
import type { ApiResponse } from '@/shared/api/types'
import type { SessionUser } from '../model/types'

export interface LoginPayload {
  tenantSlug: string
  email: string
  password: string
}

interface LoginResponse {
  accessToken: string
  user: SessionUser
}

export async function login(payload: LoginPayload) {
  const response = await apiClient.post<ApiResponse<LoginResponse>>('/auth/login', payload)
  return response.data.data
}

export async function refresh() {
  const response = await apiClient.post<ApiResponse<LoginResponse>>('/auth/refresh')
  return response.data.data
}

export async function logout() {
  await apiClient.post('/auth/logout')
}

export async function me() {
  const response = await apiClient.get<ApiResponse<SessionUser>>('/auth/me')
  return response.data.data
}
