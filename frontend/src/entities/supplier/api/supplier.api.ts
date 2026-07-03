import { apiClient } from '@/shared/api/client'
import type { Paginated } from '@/shared/api/types'
import type { Supplier } from '../model/types'

export interface ListSuppliersParams {
  limit?: number
  offset?: number
  search?: string
  active?: boolean
}

export async function listSuppliers(params: ListSuppliersParams = {}): Promise<Paginated<Supplier>> {
  const response = await apiClient.get<Paginated<Supplier>>('/suppliers', { params })
  return response.data
}

export async function getSupplier(id: string): Promise<Supplier> {
  const response = await apiClient.get<{ data: Supplier }>(`/suppliers/${id}`)
  return response.data.data
}

export interface SupplierInput {
  name: string
  document?: string
  email?: string
  phone?: string
  active: boolean
}

export async function createSupplier(data: SupplierInput): Promise<Supplier> {
  const response = await apiClient.post<{ data: Supplier }>('/suppliers', data)
  return response.data.data
}

export async function updateSupplier(id: string, data: SupplierInput): Promise<Supplier> {
  const response = await apiClient.patch<{ data: Supplier }>(`/suppliers/${id}`, data)
  return response.data.data
}

export async function removeSupplier(id: string): Promise<void> {
  await apiClient.delete(`/suppliers/${id}`)
}
