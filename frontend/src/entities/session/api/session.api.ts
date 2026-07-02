import { apiClient } from '@/shared/api/client'

export interface LoginPayload {
  tenantSlug: string
  email: string
  password: string
}

export async function login(payload: LoginPayload) {
  const response = await apiClient.post('/auth/login', payload)
  return response.data
}
